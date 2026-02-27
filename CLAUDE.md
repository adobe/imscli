# Code Review & Improvement Plan for imscli

## Summary

Full review of all Go source files in the `imscli` project — a CLI tool for Adobe's
IMS (Identity Management Service). The codebase has ~44 Go source files across `cmd/`,
`ims/`, and `output/` packages.

---

## 1. Bugs

### 1.5 Suspicious `Expires` calculations (likely overflow)

**Files:** `ims/jwt_exchange.go:65`, `ims/exchange.go:73`

```go
Expires: int(r.ExpiresIn * time.Millisecond),
```

`r.ExpiresIn` is a `time.Duration` and `time.Millisecond` is also a `time.Duration`.
Multiplying two durations produces `int64` nanoseconds-squared, which is semantically
wrong and will overflow for typical expiry values. If `ExpiresIn` is already in
milliseconds, the multiplication is incorrect.

### 1.6 Inconsistent `Expires` calculation in refresh

**File:** `ims/refresh.go:67`

Uses `int(r.ExpiresIn * time.Second)` while jwt_exchange and exchange use
`time.Millisecond`. At least one of these is wrong.

---

## 2. Idiomatic Go Issues

### 2.6 `validateURL` switch can be simplified (skipped)

**File:** `ims/config.go:61-75`

```go
return parsedURL.Scheme != "" && parsedURL.Host != ""
```

---

## 3. Structural Improvements

### 3.1 Extract IMS client creation into a helper

Every method in the `ims` package repeats:

```go
httpClient, err := i.httpClient()
// ...
c, err := ims.NewClient(&ims.ClientConfig{URL: i.URL, Client: httpClient})
```

This 8-line boilerplate appears in ~10 files. Extract into:

```go
func (i Config) newIMSClient() (*ims.Client, error) { ... }
```

### 3.2 Deduplicate validate subcommands

All four files in `cmd/validate/` have identical `RunE` logic — only the flag binding
and error message differ. These could be consolidated with a factory function.

### 3.3 Deduplicate invalidate subcommands

Same situation as validate — four nearly identical files.

### 3.5 Remove empty `internal/output/` directory

This directory exists but contains no files. Not tracked by git — local cleanup only.

### 3.8 Missing `cobra.MarkFlagRequired` on mandatory flags (skipped)

Skipped: cobra's `MarkFlagRequired` checks flag `.Changed` (CLI only), so it would
reject values provided via config file or env vars through viper.

---

## 4. Second-Wave Findings (Remaining)

### Critical / High Severity

#### 6.1.7 Private key material not zeroed from memory

**File:** `ims/jwt_exchange.go:36-51`

The private key is read into a `[]byte` and passed to `ExchangeJWT`, but the byte
slice is never zeroed after use. Key material persists in memory until GC. Should
add `defer func() { for i := range key { key[i] = 0 } }()`.

#### 6.1.8 No mutual exclusion for token fields in validate/invalidate

**Files:** `ims/validate.go:26-48`, `ims/invalidate.go:23-48`

If multiple token fields are populated simultaneously, only the first match in the
switch wins silently. No validation checks that exactly one token is provided.

#### 6.1.12 PKCE flag mutation persists on shared Config pointer

**File:** `cmd/authz/pkce.go:30`

`imsConfig.PKCE = true` mutates the shared Config pointer.

### Medium Severity

#### 6.2.2 `TokenInfo.Expires` is computed but never consumed

Set in `jwt_exchange.go:65`, `exchange.go:73`, `refresh.go:67` with wrong calculations,
but no caller ever reads the field. Dead code with active bugs.

#### 6.2.3 Missing input validation in `AuthorizeService` and `AuthorizeClientCredentials`

`ims/authz_service.go` and `ims/authz_client.go` are the only two operation methods
without a `validateXxxConfig()` function.

#### 6.2.5 `browser.Stdout = nil` has global side effects

**File:** `ims/authz_user.go:123`

Modifies a package-level variable in `pkg/browser`. Should save/restore the original.

#### 6.2.6 Inconsistent flag shorthands across commands

`-s` and `-a` mean different things in different commands.

#### 6.2.8 `cmd/profile.go:41` — default API version outdated

Profile API version defaults to `"v1"` but latest is `"v3"`.

#### 6.2.9 Hardcoded serviceCode whitelist in profile decoding

**File:** `ims/profile.go:108-111`

Four service codes are hardcoded. Should be configurable.

#### 6.2.10 URL trailing slash creates malformed JWT claim keys

**File:** `ims/jwt_exchange.go:47`

If `i.URL` ends with `/`, the claim key gets a double slash.

#### 6.2.12 API version whitelist is not future-proof

**Files:** `ims/organizations.go:23`, `ims/admin_organizations.go:23`

Hardcoded `v1`-`v6` whitelist requires code changes for new versions.

### Low Severity / Nits

| Finding | File | Description |
|---------|------|-------------|
| No test for unicode/special chars | `cmd/pretty/json_test.go` | Missing edge case coverage |
| `RawURLEncoding` rejects padded base64 | `ims/decode.go:53` | Some JWT impls include `=` padding |
| No JWE support (5-part tokens) | `ims/decode.go:43` | Only 3-part JWTs supported |
| Gzip `io.Copy(io.Discard)` unnecessary | `ims/profile.go:147` | JSON decoder already consumed stream |
| `PersistentPreRunE` can be overridden | `cmd/root.go:31` | If a subcommand defines its own, parent's is lost |
| Duplicate port defaults (cobra + const) | `cmd/authz/user.go:44` + `ims/authz_user.go:27` | Maintenance risk |
| `cmd/refresh.go:36-39` uses `map[string]interface{}` | `cmd/refresh.go` | No guaranteed field order in JSON; use struct |
| `viper.Unmarshal` ignores unknown keys | `cmd/params.go:69` | Config typos silently ignored |
| Listener not closed on panic path | `ims/authz_user.go:112-130` | Missing `defer listener.Close()` |
| Cascading/ClientSecret unconditionally sent | `ims/invalidate.go:91-97` | Should be conditional on token type |
| `version` variable has no default | `main.go:22` | Shows empty if ldflags not set |
| Double error printing from cobra | `main.go:27-29` | Root cmd doesn't SilenceErrors; leaf cmds do |

---

## 5. Remaining Change Plan

### Structural Refactoring

| # | Change | Files |
|---|--------|-------|
| 33 | Deduplicate validate subcommands with factory | `cmd/validate/*.go` |
| 34 | Deduplicate invalidate subcommands with factory | `cmd/invalidate/*.go` |
| 38 | Remove empty `internal/output/` directory | directory |
| 40 | Merge `pkce.go`/`user.go` into factory | `cmd/authz/` |
| 42 | Update default profile API version to v3 | `cmd/profile.go` |

### Testing

| # | Change | Files |
|---|--------|-------|
| 43 | Add unit tests for all validators | `ims/*_test.go` (new) |
| 44 | Add unicode/special char tests for pretty | `cmd/pretty/json_test.go` |
| 46 | Add integration tests for flag/config/env precedence | `cmd/*_test.go` (new) |

---

## Completed

The following items have been implemented and merged:

- **1.1** Fix `InvalidateToken` calling wrong validator (#60)
- **1.2** Fix log message "authorization code" → "service token" (#60)
- **1.3** Fix error message "validating" → "invalidating" (#60)
- **1.4** Replace deprecated `ioutil.ReadFile` with `os.ReadFile` (#60)
- **1.7** Fix duplicate `"exch"` alias → `"ref"` for refresh (#60)
- **1.8** Safe type assertion in profile decoding (#60)
- **1.9** Fix empty-string default scopes across 6 files (#60)
- **1.10** Fix nil/empty Scopes panic in validation (#60)
- **2.2** Remove redundant `== true` boolean comparison (#60)
- **2.3** Remove redundant `= false` zero-value init (#60)
- **2.4** Replace `%v` with `%w` in error wrapping — 78 instances (#60)
- **3.9** Fix fragile `strings.Replace` → `strings.Trim` (#60)
- **6.1.1** Fix viper BindPFlags for persistent flags (#61)
- **6.1.3** Use `http.DefaultTransport.Clone()` for proxy Transport (#62)
- **6.2.1** Make Config.Timeout functional and expose as --timeout flag (#63)
- **6.2.13** Wrap bare errors in organizations/admin methods (#61)
- **3.4** Extract shared `resolveToken()` helper (#65)
- **3.7** Use `output.PrintPrettyJSON` in admin commands (#65)
- **6.1.5** Fix OAuth shutdown deadlock and add timeout (#66)
- **6.1.6** Capture server.Serve() error in OAuth flow (#67)
- **3.6** Refactor pretty-print into `cmd/pretty` package (#68)
- **N4** Remove dead `DecodedToken.Signature` field (#68)
- **N5** Consolidate duplicate `prettyJSON` with output package (#68)
- **41** Make decode output valid JSON (single object) (#69)
- **6.1.9** Fix `fmt.Printf` → `log.Printf` in `findFulfillableData` (#69)
- **6.1.11** Fix wrong OAuth flow description in `cmd/authz/service.go` (#69)
- **6.2.4** Add trailing newline in refresh fullOutput mode (#69)
- **6.2.7** Fix grammar "an user" → "a user" (#69)
- **6.2.11** Admin commands use `pretty.JSON` for pretty-print (#68)
- **6.2.14** Declare `--organization`/`--userID` mutual exclusion via cobra (#69)
- **6.2.15** Decode output is now parseable JSON (#69)
- **6.2.16** Fix wrong doc comment "access token" → "service token" (#69)
- **N3** Fix "decodification" → "decoding" (#69)
- **2.1** Rename `ProfileApiVersion` → `ProfileAPIVersion`, `OrgsApiVersion` → `OrgsAPIVersion` (#70)
- **2.5** Replace C-style `/* */` block comments with `//` line comments (#70)

### Skipped (not applicable)

- **6.1.2** mapstructure tags — not needed, case-insensitive matching works
- **6.1.4** Client-per-request — not an issue for a CLI tool that exits after one command
- **6.1.10** Default metascopes — covered by 1.9
- **3.8** `MarkFlagRequired` — incompatible with viper config/env workflow
- **6.1.13** Missing Config `String()` — Config is never printed in the CLI
