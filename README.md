# fk 🐷

> **You typed it wrong. `fk` fixes it.**

AI-powered shell command corrector. Run a broken command, type `fk` — get the fix instantly. Built in Go. Zero runtime dependencies.

```
$ git psuh origin main
fatal: 'psuh' is not a git command. See 'git --help'.

$ fk
git push origin main
```

```
$ docker-compose up -d
ERROR: Version in "./docker-compose.yml" is unsupported.

$ fk
docker compose up -d
```

Inspired by [thefuck](https://github.com/nvbn/thefuck). Faster. No Python. No rules to maintain.

---

## Install

### With [install-release](https://pypi.org/project/install-release/) (`ir`)

```bash
pip install -U install-release
ir get https://github.com/Rishang/fk
```

Binaries go to `~/bin` — add it to your PATH:

```bash
export PATH="$HOME/bin:$PATH"
```

### Pre-built binary

Grab the binary for your platform from the [latest release](https://github.com/Rishang/fk/releases/latest):

| Platform    | Asset               |
|-------------|---------------------|
| Linux x64   | `fk-linux-amd64`   |
| Linux arm64 | `fk-linux-arm64`   |
| macOS x64   | `fk-darwin-amd64`  |
| macOS arm64 | `fk-darwin-arm64`  |

```bash
# Example: Linux x86_64
curl -fL -o fk https://github.com/Rishang/fk/releases/latest/download/fk-linux-amd64
chmod +x fk
sudo mv fk /usr/local/bin/
```

### From source

```bash
git clone https://github.com/Rishang/fk
cd fk
task install        # installs to $(go env GOPATH)/bin/fk
```

Requires Go 1.21+.

---

## Setup

### 1. Pick your AI provider

`fk` works with any major AI provider — or your own local model:

```bash
# Claude (Anthropic)
fk config set --provider claude --api-token sk-ant-xxxx --model claude-sonnet-4-20250514

# OpenAI
fk config set --provider openai --api-token sk-xxxx --model gpt-4o

# OpenRouter (100+ models, one key)
fk config set --provider openrouter --api-token sk-or-xxxx --model openai/gpt-4o-mini

# Google Gemini
fk config set --provider gemini --api-token AIzaSy-xxxx --model gemini-1.5-flash

# Local model (Ollama, LiteLLM, etc.)
fk config set --provider openai --api-token x --base-url http://localhost:11434/v1 --model llama3.2
```

Config lives at `~/.config/fk/config.yaml`.

### 2. Hook into your shell

```bash
# bash — add to ~/.bashrc
echo 'eval "$(fk --shell-init bash)"' >> ~/.bashrc && source ~/.bashrc

# zsh — add to ~/.zshrc
echo 'eval "$(fk --shell-init zsh)"' >> ~/.zshrc && source ~/.zshrc

# fish — add to ~/.config/fish/config.fish
echo 'fk --shell-init fish | source' >> ~/.config/fish/config.fish
```

### 3. Break things. Fix them.

```bash
$ kubectll get pods
command not found: kubectll

$ fk
kubectl get pods
```

---

## How it works

```
failing command + exit code
         │
         ▼
   shell hook captures context
         │
         ▼
      fk sends prompt to AI
         │
         ▼
   AI returns fix (commands only, no prose)
         │
         ▼
      fk prints the fix
```

1. The shell hook (`PROMPT_COMMAND` / `precmd`) captures the last **non-fk** command and its exit code into env vars — so running `fk` multiple times always points at the original failing command.
2. `fk` builds a terse prompt and calls your configured AI provider.
3. AI responds with `FIX: <cmd>` or `STEPS:\n$ cmd1\n$ cmd2…`
4. `fk` prints the corrected command(s) as plain text.

---

## Options

| Config key   | Default                    | Description                                            |
|--------------|----------------------------|--------------------------------------------------------|
| `provider`   | `claude`                   | AI backend: `claude`, `openai`, `openrouter`, `gemini` |
| `api_key`    | —                          | Your provider API token                                |
| `model`      | `claude-sonnet-4-20250514` | Model name                                             |
| `base_url`   | *(provider default)*       | Override endpoint — for proxies or local models        |
| `max_tokens` | `512`                      | Max tokens in AI response                              |
| `auto_run`   | `false`                    | Execute suggestion without confirmation                |

```bash
fk config show          # view current settings
fk config set --provider openai --api-token sk-xxxx --model gpt-4o
```

---

## Flags

| Flag              | Description                                                        |
|-------------------|--------------------------------------------------------------------|
| `--rerun` / `-r`  | Re-run the failed command to capture live output before asking AI  |
| `--auto-run`      | Run the fix immediately without prompting                          |
| `--debug`         | Print raw AI response before parsing                               |
| `--cmd`           | Provide the failed command explicitly (no shell hook needed)       |
| `--exit-code`     | Provide the exit code explicitly                                   |
| `--output`        | Provide captured output explicitly                                 |

```bash
# Direct usage — no shell integration needed
fk --cmd "kubectl get pods" --exit-code 1
fk --cmd "cargo build" --exit-code 101 --rerun --debug
fk --cmd "pip install numpy" --exit-code 1 --auto-run
```

---

## `fk cat` — files to prompt

Dump files or a directory into a clean `<file>…</file>` format, ready to paste into any AI prompt.

```bash
fk cat go.mod go.sum          # specific files
fk cat ./internal             # walk a directory (respects .gitignore)
fk cat                        # walk current directory
```

Output:

```
<file go.mod>
module github.com/Rishang/fk
...
</file go.mod>
<file internal/config/config.go>
...
</file internal/config/config.go>
```

---

## Building

```bash
task build        # build for current platform → ./dist/fk
task dist         # cross-compile: linux/darwin/windows amd64+arm64
task test         # run tests
task install      # install to $GOPATH/bin
```
