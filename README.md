# Zencode

**Zencode** is a powerful, terminal-first AI coding agent built for **Termux** (Android) and other Linux environments. It brings a Grok Build + OpenCode style experience to your phone or small device: explore codebases, plan changes, precisely edit files, run commands safely, track todos, delegate sub-tasks, all from a beautiful keyboard-centric TUI.

- Single static binary (small footprint, fast startup — perfect for Termux).
- First-class support for **Grok / xAI** (just set `XAI_API_KEY`).
- Works with any OpenAI-compatible provider (OpenAI, Groq, Anthropic via gateways, Gemini, Ollama, llama.cpp servers, etc.).
- Precise edits (`search_replace`), visible task tracking, Plan Mode, subagents, project rules, permissions with YOLO mode.
- Built for real work on the go in Termux.

## Quick Start (Termux)

```bash
# 1. One-time Termux prep (recommended)
pkg update && pkg upgrade
pkg install git ripgrep

# Optional but useful for building from source later
# pkg install golang

# 2. Get the latest arm64 binary (example — replace with real release URL when available)
curl -L https://github.com/aungmyatmoe/zencode/releases/latest/download/zencode-linux-arm64 -o $PREFIX/bin/zencode
chmod +x $PREFIX/bin/zencode

# Or build from source (after cloning)
# go build -o $PREFIX/bin/zencode ./cmd/zencode

# 3. Set your key (xAI/Grok recommended to start)
export XAI_API_KEY="xai-your-key-here"

# 4. Run in any project
cd ~/storage/shared/my-project   # or wherever your code lives (run termux-setup-storage if needed)
zencode
```

On first run it will load your project context (respects `.gitignore`), show a short tree + any `ZENCODE.md` / `AGENTS.md` rules, and you're chatting with a capable coding agent.

Press `?` for keyboard help, `/` for slash commands, `Shift+Tab` to enter Plan Mode.

## Features (Grok + OpenCode inspired)

- Rich TUI scrollback with live streaming, collapsible thinking/tool blocks, and beautiful diffs for edits.
- Core tools: `read_file`, `list_dir`, `grep` (uses `rg` when available), `search_replace` (exact unique string match — safe & precise), `run_terminal_command` (with smart safe-command fast-paths), `todo_write`.
- **Plan Mode**: Explore and design before any edits. Writes a reviewable `plan.md` in the session; approve to execute.
- Subagents: Delegate research (`explore`), planning, or full work in parallel.
- Project rules: Drop a `ZENCODE.md` (or `AGENTS.md`, `CLAUDE.md`) at repo root or in subdirs — automatically loaded with proper precedence.
- Permissions & safety: Asks before risky edits or shell commands. `Ctrl+O` or `/yolo` for auto-approve. Configurable allow/deny rules.
- Model switching, headless mode for scripts/CI, sessions you can resume.
- Vim-friendly navigation (toggle with `/vim-mode`).

## Configuration

Environment variables (highest priority):
- `XAI_API_KEY` — for Grok / xAI (recommended starting point)
- `ZENCODE_API_KEY`, `ZENCODE_BASE_URL`, `ZENCODE_MODEL`

Config file: `~/.zencode/config.toml` (or `$XDG_CONFIG_HOME/zencode/config.toml`).

See docs or run `zencode` and use the settings UI / slash commands.

## Headless / Scripting

```bash
zencode -p "Review the auth changes and suggest improvements" --cwd ./my-api --yolo --output json
```

## Philosophy & Limitations (v1)

Zencode targets **practical coding on Termux**. It emphasizes reliable tool use, safety, visible planning (todos + Plan Mode), and a delightful terminal experience over bleeding-edge local inference in v1.

- v1 is API-first (bring your keys). Local models supported later via any OpenAI-compatible server you run yourself (Ollama, llama.cpp, etc.).
- No LSP in v1 (heavy on device). Excellent `rg` + LLM understanding gets you very far. LSP is on the roadmap.
- Keep an eye on RAM — use compact models or aggressive `/compact` when working on big codebases on phones.

## Development

```bash
git clone https://github.com/aungmyatmoe/zencode
cd zencode
go build -o zencode ./cmd/zencode
./zencode --help
```

Contributions, issues, and Termux-specific feedback very welcome.

## License

MIT

---

Built for the terminal. Made for Termux. Powered by great models (starting with Grok).
