# TODO

## Slack Webhook Notifications

Add Slack notifications on CI failures. This would post to a team channel when
the CI pipeline, govulncheck, or release workflows fail.

## Comprehensive Test Coverage

Current coverage is ~46%. The tested code covers pure functions (validation,
decoding, client construction, profile parsing) and fuzz tests. The untested
code is mostly thin wrappers that call `ims-go` through `newIMSClient()`.

### What we tried

**httptest approach**: Create a local HTTP server, point `Config.URL` at it, and
let the real code path execute end-to-end. This works mechanically — the tests
pass — but the assertions are shallow. The fake server returns canned responses
regardless of input, so we're testing that our code passes arguments through and
wraps results, not that it does anything meaningful. The tests are also tightly
coupled to `ims-go` implementation details (endpoint paths like
`/ims/validate_token/v1`), meaning they break if the library changes internals
even though our code is fine.

### What the untested code actually looks like

Each public function follows the same pattern:

```
validate config → create client → call ims-go → wrap result
```

The validation is already tested. The client creation is already tested. What
remains is ~5 lines of glue per function where httptest gives line coverage but
not meaningful assertions.

### Open question

Is there an approach that tests our behavior rather than `ims-go`'s wire format?
Options to explore:

- Accept the shallow httptest coverage as "good enough" for glue code
- Focus testing effort on new features where logic is non-trivial
- Explore contract testing if `ims-go` publishes response schemas
- **Mock IMS server**: Build a standalone fake IMS service that understands the
  real API contract (paths, params, auth headers) and returns realistic responses.
  Run the CLI binary against it as an integration test. Tests real behavior
  end-to-end including flag parsing, config loading, and output formatting —
  not just the `ims/` package. Tradeoff: the mock needs maintaining as the API
  evolves, but catches issues that unit tests never will.
- **Real IMS integration tests**: Run a subset of tests against the actual IMS
  API using a dedicated test client. Gate behind a build tag
  (`//go:build integration`) or env var so they don't run in CI by default.
  Tradeoff: requires credentials, is slow, can flake on network issues — but
  is the only way to catch real API drift or behavioral changes.
