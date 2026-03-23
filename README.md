# kidmode

A macOS terminal app that turns your keyboard into a toy for kids. Every keypress triggers colorful visual effects and fun sounds while the system is locked down so they can't escape.

## Install

```bash
brew tap tdiderich/kidmode
brew install kidmode
```

Or build from source:

```bash
git clone https://github.com/tdiderich/kidmode.git
cd kidmode
make build
```

## Usage

```bash
kidmode
```

Your terminal goes full-screen with visual effects. Kids can bang on the keyboard freely.

To exit, type: `adulttakeover`

You can set a custom exit password:

```bash
KIDMODE_PASSWORD=mypassword kidmode
```

## Requirements

- macOS
- **Accessibility permissions** — on first run, macOS will prompt you to grant access. Go to System Settings > Privacy & Security > Accessibility and enable your terminal app.

## What gets blocked

- System shortcuts (Cmd+Tab, Cmd+Q, Cmd+Space, etc.)
- Mouse input (clicks, movement, scrolling)
- Media keys (volume, play/pause, skip)
- Screenshot shortcuts (Cmd+Shift+3/4/5)
- Mission Control, Launchpad, Exposé
- Terminal escape sequences (Ctrl+C, Ctrl+Z, etc.)

## License

MIT
