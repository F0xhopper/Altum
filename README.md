# Altum

[![CI](https://github.com/F0xhopper/Altum/actions/workflows/ci.yml/badge.svg)](https://github.com/F0xhopper/Altum/actions/workflows/ci.yml)

A minimalist deep work companion for the terminal. Altum times your focused
work sessions, walks you through a short review when you stop, and keeps
everything in a local SQLite database so you can see how your deep work adds
up over time.

The name comes from the Latin *altum*, "deep".

![Altum demo](assets/demo.gif)

## How it works

1. `altum start` opens a full-screen timer. Work until you're done, then press
   Enter to stop.
2. A four-step review captures the session: the milestone you reached, focus
   quality (1–5), any interruptions, and a short reflection. Only the
   milestone is required.
3. The session is saved to a local SQLite database.
4. `altum report` summarises your recent sessions: total deep work time,
   average session length, average focus rating, best day, and more.

## Install

With Go 1.25+:

```sh
go install github.com/F0xhopper/Altum/cmd/altum@latest
```

Or download a prebuilt binary from the
[releases page](https://github.com/F0xhopper/Altum/releases/latest):

```sh
# macOS (Apple Silicon shown; use Darwin_x86_64 for Intel)
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_Darwin_arm64.tar.gz
tar -xzf Altum_Darwin_arm64.tar.gz
sudo mv altum /usr/local/bin/
```

## Usage

```sh
altum            # interactive menu
altum start      # start a deep work session
altum report     # report on the last 7 days
altum report -d 30
altum config set db_path ~/notes/altum.db
altum config get
```

Example report:

```
═══════════════════════════════════════════════════════════
  Deep Work Report: Jul 19, 2026 - Jul 25, 2026
═══════════════════════════════════════════════════════════

Total sessions: 9
Total deep work: 11.0 hours (660 minutes)
Average session: 73 minutes
Average rating: 4.0 / 5
Best day: Jul 22 – 2.9 hours (2 sessions)
Longest session: 110 minutes (Jul 22)
Days with deep work: 7 / 7 (100%)

Top performing days:
Jul 22: 2.9h (★★★★☆)
Jul 20: 2.3h (★★★☆☆)
Jul 21: 1.8h (★★☆☆☆)
```

## Keybindings

During a session:

| Key | Action |
| --- | --- |
| `enter` / `space` | Stop the session and begin the review |
| `ctrl+c` / `q` | Quit without saving |

During the review:

| Key | Action |
| --- | --- |
| `enter` | Confirm and continue to the next step |
| `tab` | Skip an optional step |
| `shift+tab` | Go back a step |
| `ctrl+s` | Save the session (on the final step) |

In the menu: `↑`/`k` and `↓`/`j` to navigate, `enter` to select, `q` to quit.

## Data and configuration

Sessions are stored in a single SQLite database. By default:

```
~/.local/share/altum/altum.db   # $XDG_DATA_HOME is respected if set
```

The location can be overridden with the `--db_path` flag, or persistently:

```sh
altum config set db_path /path/to/altum.db
```

Configuration lives at `~/.config/altum/config.yaml`. A different config file
can be passed with `--config`.

## Development

```sh
go build -o altum ./cmd/altum
go vet ./...
```

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea),
[Cobra](https://github.com/spf13/cobra), and a pure-Go SQLite driver
(no CGO required).

## License

Apache License 2.0. See [LICENSE](LICENSE).
