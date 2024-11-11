# Stay Active

This is a CLI application designed to help avoid micromanagement. Use it alongside your Slack app (or any other platform you find useful), especially if you need to take a quick break or deal with an emergency and don’t want to lose any of your daily hours from your paycheck.

![Demo](doc/preview.gif?raw=true "Demo")

Micromanagement can be a reality we all face at times, and this tool helps address that by keeping your Slack status active. Slack considers users active when they interact with the app at least every 30 minutes. By using this application, you can set it up according to your preferences, open a chat with yourself, and step away as needed.

## Requirements

- [Go (Golang) >= 1.23.2](https://go.dev/doc/install)
- [Robotgo requirements](https://github.com/go-vgo/robotgo#user-content-requirements): Make sure to install every requirement for your OS

## How to run

Simply run:

```bash
make install
make run

```

That's it! Now, your computer will continue pressing keys until you either terminate this application or the timeout period ends.
