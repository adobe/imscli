# OAuth Local Server Shutdown Deadlock

## Context

The `authz user` command starts a local HTTP server to handle the OAuth2 authorization
code grant flow. A browser is opened, the user authenticates, and the browser redirects
back to the local server with the authorization code. The server exchanges it for a token
and sends the result through an unbuffered channel.

## The Problem

There are two issues in the original shutdown logic:

1. `Shutdown()` is called with `context.Background()`, which has no timeout
2. A deadlock can occur between the handler trying to write to a channel and `Shutdown()`
   waiting for the handler to return

## Setup: Unbuffered Channels

The login server creates unbuffered channels (`ims-go/login/server.go:91-92`):

```go
resCh = make(chan *ims.TokenResponse)
errCh = make(chan error)
```

An unbuffered channel blocks the writer until a reader is ready to receive.

## Normal Flow (No Deadlock)

1. The `select` in `ims/authz_user.go:137-145` waits on three cases:

   ```go
   select {
   case serr = <-server.Error():
       log.Println("The IMS HTTP handler returned an error message.")
   case resp = <-server.Response():
       log.Println("The IMS HTTP handler returned a message.")
   case <-time.After(time.Minute * 5):
       fmt.Fprintf(os.Stderr, "Timeout reached waiting for the user to finish the authentication ...\n")
       serr = fmt.Errorf("user timed out")
   }
   ```

2. Browser callback arrives. The handler writes to `resCh` (`ims-go/login/result.go:35-40`):

   ```go
   select {
   case h.resCh <- result:
       // Result sent.
   case <-r.Context().Done():
       // Request cancelled.
   }
   ```

3. Our `select` reads from `server.Response()` (which returns `resCh`). Both sides proceed.

4. `Shutdown()` is called. No in-flight handlers, so it completes immediately.

## Deadlock Flow

1. The `select` in `ims/authz_user.go:137-145` is waiting on the three cases.

2. The 5-minute timeout fires first. We leave the `select` and reach `Shutdown()`
   (`ims/authz_user.go:147-150`).

3. Now **nobody is reading** from `resCh` or `errCh`.

4. The browser callback arrives late (the user finished auth just after the timeout).
   The HTTP handler chain reaches `resultHandler.ServeHTTP` (`ims-go/login/result.go:27`),
   processes the token exchange successfully, and hits (`ims-go/login/result.go:35-40`):

   ```go
   select {
   case h.resCh <- result:
       // Result sent.
   case <-r.Context().Done():
       // Request cancelled.
   }
   ```

5. `h.resCh <- result` **blocks** — unbuffered channel, no reader.

6. `r.Context().Done()` — the request context is controlled by `http.Server`. `Shutdown`
   does **not** cancel request contexts while waiting for handlers to return. So this
   case doesn't fire either.

7. The handler is **stuck** in this select forever (`ims-go/login/result.go:35-40`).

8. `Shutdown()` (`ims-go/login/server.go:143-148`) calls `s.server.Shutdown(ctx)` which
   waits for in-flight handlers to return. It waits for step 7 to complete. It never does.

   ```go
   func (s *Server) Shutdown(ctx context.Context) error {
       defer close(s.errCh)
       defer close(s.resCh)
       return s.server.Shutdown(ctx)
   }
   ```

9. `defer close(s.resCh)` and `defer close(s.errCh)` (`ims-go/login/server.go:144-145`)
   **never execute** because `s.server.Shutdown` never returns.

10. **Circular dependency**: `Shutdown` waits for the handler, the handler waits for
    someone to read the channel, and nobody is reading.

## Practical Impact for a CLI

The deadlock doesn't prevent the CLI from exiting. With a timeout on `Shutdown`, the
context expires, `Shutdown` returns an error, and execution continues. The stuck handler
goroutine is cleaned up when the process exits. For a long-running server this would be
a real resource leak; for a CLI it's cosmetic (the user sees an error message about the
shutdown timeout instead of a clean exit).

## The Fix

Two changes in `ims/authz_user.go`:

### 1. Add a timeout to Shutdown

Replace `context.Background()` (no deadline) with a 10-second timeout:

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err = server.Shutdown(shutdownCtx); err != nil {
    return "", fmt.Errorf("error shutting down the local server: %w", err)
}
```

The `defer cancel()` releases the internal timer resources if `Shutdown` completes
before the deadline. This is idiomatic Go — `go vet` warns if the cancel function
from `context.WithTimeout` is not called.

### 2. Drain channels before Shutdown

Start goroutines that read from both channels **before** calling `Shutdown`:

```go
go func() {
    for range server.Response() {
    }
}()
go func() {
    for range server.Error() {
    }
}()
```

This breaks the deadlock:

1. Drainers start reading from `resCh` and `errCh`
2. `Shutdown` is called and waits for in-flight handlers
3. If a handler writes to `resCh` (`ims-go/login/result.go:36`), the drainer reads it.
   The handler unblocks and returns.
4. `Shutdown` sees all handlers are done. It calls `close(s.resCh)` and `close(s.errCh)`
   (`ims-go/login/server.go:144-145`).
5. `for range` exits on a closed channel. The drainer goroutines return.

Everyone cleans up.
