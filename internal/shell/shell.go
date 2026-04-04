// Package shell generates the shell-specific init scripts that hook fk into the terminal.
// Usage: eval "$(fk --shell-init bash)"  or  eval "$(fk --shell-init zsh)"
package shell

import "fmt"

// Init returns the shell snippet for the given shell name.
// The snippet:
//  1. Hooks into the shell's pre-command mechanism to record exit code + last command
//  2. Defines a `fk` function that invokes the real binary
func Init(shell, binaryPath string) (string, error) {
	switch shell {
	case "bash":
		return bashInit(binaryPath), nil
	case "zsh":
		return zshInit(binaryPath), nil
	case "fish":
		return fishInit(binaryPath), nil
	default:
		return "", fmt.Errorf("unsupported shell %q — supported: bash, zsh, fish", shell)
	}
}

// bashInit emits PROMPT_COMMAND-based hooks for bash
func bashInit(bin string) string {
	return fmt.Sprintf(`
# fk shell integration — added by: fk --shell-init bash
# Paste this into ~/.bashrc or run: eval "$(fk --shell-init bash)"

_fk_hook() {
  # Capture exit code immediately before anything else runs
  local _fk_ec=$?
  # Grab last history entry, strip the history number prefix
  local _fk_cmd=$(HISTTIMEFORMAT="" history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
  # Skip fk commands so fk_LAST_CMD always holds the last non-fk command
  case "$_fk_cmd" in
    fk|fk\ *) ;;
    *) export fk_EXIT_CODE=$_fk_ec; export fk_LAST_CMD=$_fk_cmd ;;
  esac
}

# Prepend to PROMPT_COMMAND so it fires before any other hooks
if [[ "$PROMPT_COMMAND" != *"_fk_hook"* ]]; then
  PROMPT_COMMAND="_fk_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi

# fk shell function — invokes the binary directly
fk() {
  %s "$@"
}
`, bin)
}

// zshInit emits precmd-based hooks for zsh
func zshInit(bin string) string {
	return fmt.Sprintf(`
# fk shell integration — added by: fk --shell-init zsh
# Paste this into ~/.zshrc or run: eval "$(fk --shell-init zsh)"

_fk_hook() {
  local _fk_ec=$?
  # fc -ln -1 prints last command without line number
  local _fk_cmd=$(fc -ln -1 2>/dev/null | sed 's/^[[:space:]]*//')
  # Skip fk commands so fk_LAST_CMD always holds the last non-fk command
  case "$_fk_cmd" in
    fk|fk\ *) ;;
    *) export fk_EXIT_CODE=$_fk_ec; export fk_LAST_CMD=$_fk_cmd ;;
  esac
}

# Register with zsh's precmd array (runs after every command, before prompt)
autoload -Uz add-zsh-hook
add-zsh-hook precmd _fk_hook

# fk shell function
fk() {
  %s "$@"
}
`, bin)
}

// fishInit emits event-based hooks for fish shell
func fishInit(bin string) string {
	return fmt.Sprintf(`
# fk shell integration — fish
# Add to ~/.config/fish/config.fish or run: fk --shell-init fish | source

function _fk_hook --on-event fish_postexec
  # Skip fk commands so fk_LAST_CMD always holds the last non-fk command
  set -l _fk_cmd $argv[1]
  if not string match -qr '^fk(\s|$)' -- "$_fk_cmd"
    set -gx fk_EXIT_CODE $status
    set -gx fk_LAST_CMD $_fk_cmd
  end
end

function fk
  %s $argv
end
`, bin)
}
