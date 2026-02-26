# OAuth Local Server: Unhandled Serve() Error

## The Problem

`go server.Serve(listener)` (`ims/authz_user.go:130`) launches a goroutine. The `error`
return from `Serve` is discarded.

`http.Server.Serve` (`ims-go/login/server.go:136-138`) delegates to the standard
library's `http.Server.Serve`. If the listener fails, `Serve` returns an error
immediately. No HTTP handler ever runs.

Since no handler runs, nothing writes to `resCh` or `errCh`
(`ims-go/login/server.go:91-92`).

The `select` (`ims/authz_user.go:137-145`) is waiting on `server.Error()`,
`server.Response()`, and a 5-minute timeout. Neither channel will ever receive a value.

The user waits 5 minutes for the timeout to fire, then gets a generic "user timed out"
error instead of the actual listener failure.

## Practical Likelihood

In practice, this is hard to trigger. The most obvious cause — port already in use — is
caught earlier at `net.Listen` (`ims/authz_user.go:112-114`), before `Serve` is called.
Once `net.Listen` succeeds, `Serve` with a valid listener will block on `Accept()` until
`Shutdown()` is called. A failure would require something unusual like the file descriptor
being invalidated between `Listen` and `Serve`.

The fix is applied as a defensive pattern: capturing goroutine errors and reacting to
them is good practice regardless of how likely the failure is, and it has no cost.

## The Fix

Capture the `Serve` error via a buffered channel and add it as a case in the `select`:

```go
serveCh := make(chan error, 1)
go func() {
    serveCh <- server.Serve(listener)
}()
```

The new `select` case:

```go
case serr = <-serveCh:
    log.Println("The local server stopped unexpectedly.")
```

This gives the user an immediate, accurate error instead of a 5-minute wait.

## Why the Channel Must Be Buffered

In the normal flow, `Serve` runs in the background while the `select`
(`ims/authz_user.go:137-145`) waits for one of three things: response, error, or
timeout.

The browser callback arrives. The handler writes to `resCh`
(`ims-go/login/result.go:35-36`). The `select` reads from `server.Response()` and
proceeds.

We reach `Shutdown()` (`ims/authz_user.go:159-163`). `Shutdown` closes the listener
and waits for handlers to finish, then returns.

Once the listener is closed, `server.Serve(listener)` (`ims-go/login/server.go:136-138`)
returns `http.ErrServerClosed`.

Our goroutine tries to write this error to `serveCh`:

```go
go func() {
    serveCh <- server.Serve(listener)
}()
```

Nobody is reading from `serveCh` — we already left the `select` back when the browser
callback arrived.

If `serveCh` were **unbuffered**, the goroutine blocks on the write forever — a goroutine
leak (the same class of problem described in `docs/oauth-shutdown-deadlock.md`).

With `serveCh` buffered at capacity 1, the write succeeds immediately into the buffer.
The goroutine exits. The channel and its buffered value are cleaned up by GC when nothing
references them anymore.
