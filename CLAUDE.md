# CLAUDE.md

Go Discord bot (`eod/`) + web backend. Prefer the smallest change that works.

## Style
- Read neighboring code first and copy it. New code should be indistinguishable from what's around it.
- Positional SQL params (`$1`), not named/`NamedExec`.
- Comments are terse fragments inside functions (`// Get poll & vote cnt`, `// Update`, `// Check`). No godoc sentences.
- `var x []T`, not `x := []T{}`. Consts at the top of the file.
- Name things to pair with what exists (`pollSuccess` → `pollReject`).

## Design
- No new abstraction unless the existing one can't do it. Per-guild state goes in `types.ServerMem`, not a new map.
- No callback fields, worker queues, or extra mutexes to solve a problem that hasn't happened.
- Don't touch unrelated code or change signatures package-wide for one caller.
- `polls` imports `base`, so `base` can never import `polls`. Cross-module wiring goes in package `eod`.

## Before claiming
- Verify with the compiler/grep instead of asserting. `go build ./... && go vet ./...` before saying it's done.
- Say plainly what wasn't tested (nothing here runs against live Discord or the DB).
