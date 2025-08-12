# Stay Active

This is a CLI application designed to help avoid micromanagement. Use it alongside your Slack app (or any other platform you find useful), especially if you need to take a quick break or deal with an emergency and don’t want to lose any of your daily hours from your paycheck.

Micromanagement is something we may all experience from time to time, and this tool provides a solution by helping keep your Slack status set to "online". Slack automatically marks users as "active" when they interact with the app at least once every 30 minutes. With this application, you can open a chat with yourself and let it run to keep your status active, allowing you to step away as needed without appearing offline.

## Requirements

- [Go (Golang)](https://go.dev/doc/install)

## Installation & usage

To install, you can run:

```bash
go install github.com/danielmesquitta/stay-active@latest
```

Now, to start, run:

```bash
stay-active --interval 5m --duration 1h --verbose
```

If you want to know the command with more details, run:

```bash
stay-active --help
```
