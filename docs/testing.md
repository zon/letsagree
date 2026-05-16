# Testing

## Unit and integration tests

Write tests for code logic. Use unit tests for pure functions and integration tests where real I/O (HTTP, database, file system) is involved. Prefer integration tests over mocking when the real dependency is lightweight enough to use in a test.

**Backend (Go):** use [Testify](https://github.com/stretchr/testify) — `go test ./...` from `backend/`.  
**Frontend (TypeScript):** use [Bun Test](https://bun.sh/docs/cli/test) — `bun test` from `frontend/`.

## Commands

Commands (CLI flags, just recipes, shell invocations) are not covered by automated tests. Test them manually by running them.

## Agents

Run commands yourself — do not describe what to run and wait for a human to verify. If a command is part of the work, execute it, observe the output, and confirm it behaves correctly before marking the requirement done.
