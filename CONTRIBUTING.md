# Contributing

Thanks for taking the time to contribute.

This repository is maintained as a **portfolio-quality reference implementation**. Contributions are welcome—bug fixes, documentation improvements, tests, and small enhancements are all appreciated.

---

## Ground Rules

- Keep changes **small and focused** (one improvement per PR when possible).
- Prefer **clarity over cleverness**.
- Keep **public APIs stable** unless the change is clearly justified.
- Add or update tests when behavior changes.
- Follow existing style and patterns in the codebase.

---

## Getting Started

### Prerequisites

- **Go** (use the version defined in `go.mod` if present)
- **Docker** (optional, if the project supports containerized dev)
- `make` (optional, if a `Makefile` is provided)

### Common Commands

If the repo has a `Makefile`, these are typical:

```bash
make install
make dev
make test
make check
```

If not, these are good defaults:

```bash
go test ./...
go vet ./...
gofmt -w .
```

---

## Git Workflow

### Branches

- `main`  
  The stable branch. Always expected to be green in CI.

- Feature branches  
  Create a branch for any change.

Branch naming conventions:

- `feat/<short-description>` – new features
- `fix/<short-description>` – bug fixes
- `chore/<short-description>` – tooling/refactors/cleanup
- `docs/<short-description>` – documentation-only changes

Examples:

```text
feat/add-status-endpoint
fix/handle-empty-upload
docs/update-readme
```

### Commits

Write commit messages that explain intent:

```text
feat: add submission status endpoint
fix: return 409 on duplicate email
docs: clarify local setup steps
chore: simplify test helpers
```

### Pull Requests

1. Create a feature branch:
   ```bash
   git checkout -b feat/your-change
   ```
2. Make changes + run checks:
   ```bash
   gofmt -w .
   go test ./...
   go vet ./...
   ```
3. Commit and push:
   ```bash
   git add -A
   git commit -m "feat: describe the change"
   git push origin feat/your-change
   ```
4. Open a PR targeting `main`

---

## Code Quality Checklist

Before opening a PR, please ensure:

- [ ] Code is formatted (`gofmt`)
- [ ] Tests pass (`go test ./...`)
- [ ] Static checks pass (`go vet ./...`)
- [ ] Docs updated (if behavior/usage changed)
- [ ] No secrets committed (API keys, tokens, `.env`, etc.)

---

## Tests

- Add tests for new behavior.
- Prefer **table-driven** tests where it improves readability.
- Keep tests deterministic.

---

## Security

If you find a security issue, please **open a GitHub issue** and clearly mark it as **Security**.

Include:
- description of the issue
- steps to reproduce
- impact
- suggested fix (if you have one)

---

## Documentation Style

- Keep README instructions copy/paste friendly.
- Use placeholder sample data (e.g., `John Smith`, `john.smith@example.com`) rather than real personal data.
- Prefer short sections and bullet lists.

---

## License

By contributing, you agree that your contributions will be licensed under the project’s license.
