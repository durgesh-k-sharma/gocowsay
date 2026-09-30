# gocowsay

A cowsay clone in Go. Reads a message from arguments or a pipe, wraps it in an ASCII speech balloon, and prints a figure you can swap with a flag.

```sh
fortune | gocowsay
```

## Features

- Reads the message from arguments or stdin
- Word-wraps the message into a speech balloon
- Multiple figures, selectable with a flag
- Thought balloons as well as speech balloons
- Figures are embedded in the binary, so there are no data files to install
- Standard library only

## Installation

With Go 1.16 or newer:

```sh
go install github.com/durgesh-k-sharma/gocowsay@latest
```

Or build from source:

```sh
git clone https://github.com/durgesh-k-sharma/gocowsay.git
cd gocowsay
go build
```

## Usage

```sh
# Pipe anything into it
fortune | gocowsay
echo "Go is fun" | gocowsay -f cat

# Or pass the message as arguments
gocowsay hello world
gocowsay -t "hmm..."
gocowsay -W 20 "a long message that will be wrapped at twenty columns"

# List available figures
gocowsay -l
```

### Flags

| Flag        | Default | Description                              |
|-------------|---------|------------------------------------------|
| `-f name`   | `cow`   | Figure to use                            |
| `-l`        |         | List available figures and exit          |
| `-t`        | off     | Think instead of speak                   |
| `-W n`      | `40`    | Wrap the message at `n` columns          |

## Adding a figure

Drop a text file into `figures/`, for example `figures/dog.txt`, and rebuild. The filename (without `.txt`) becomes the figure name, so this one is used with `gocowsay -f dog`.

Put `$thoughts` where the tail connecting the figure to the balloon goes. It is replaced with `\` when speaking and `o` when thinking.

Figures are embedded in the binary with `go:embed`, so rebuild after adding one.

## Project layout

```
main.go         the program
figures/*.txt   the figures
```
