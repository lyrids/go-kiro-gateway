# CLAUDE.md

Go gateway that exposes OpenAI- and Anthropic-compatible APIs on top of the Kiro service. Two backends: direct HTTP (default) and ACP (local `kiro-cli acp`). See `README.md` and `docs/ARCHITECTURE.md` for details.

## Layout

- Go module lives in `gateway/` . Entry point: `gateway/cmd/gateway`.
- Packages under `gateway/internal/`: `server` (routes), `converter` (OpenAI/Anthropic/Responses to Kiro), `streaming`, `parser`, `backend` (HTTP + ACP), `auth`, `config`, `middleware`, `models`, `resolver`, `cache`, `thinking`, `truncation`, `tokenizer`, `errors`, `debug`, `logging`, `client`.
- `plugins/go-kiro-gateway-models` is a separate TypeScript package. `.agents/personas/` is the source for `AGENTS.md` and `.claude/agents/`.

## Commands

Run from `gateway/`:

- `make build` builds the `go-kiro-gateway` binary (CGO disabled)
- `make test` runs `go test -race -count=1 ./...`
- `make lint` runs `go vet ./...`
- `go test ./internal/<pkg> -run <Name>` runs a single test

## Conventions

- Configuration comes from environment variables (`.env`, template in `.env.example`). Add new settings to `internal/config` and to `.env.example`.
- Never read, print, or commit `.env` files or credentials.
- Do not edit `AGENTS.md` or `.claude/agents/` directly. Edit `.agents/personas/` and run `.agents/sync-agents.sh`.
- Add or update tests next to the code you change, and run `make test` and `make lint` before finishing.
