# gocowsay

A cowsay clone in Go. Reads a message from arguments or a pipe, wraps it in an
ASCII speech balloon, and prints a figure you can swap with a flag.

## Build

    go build
    go install    # optional: puts gocowsay in $(go env GOPATH)/bin

## Usage

    fortune | gocowsay
    gocowsay hello world
    echo "Go is fun" | gocowsay -f cat
    gocowsay -t "hmm..."
    gocowsay -W 20 "a long message that will be wrapped at twenty columns"
    gocowsay -l

## Flags

    -f name   figure to use (default "cow")
    -l        list available figures and exit
    -t        think instead of speak
    -W n      wrap the message at n columns (default 40)

## Adding a figure

Drop a text file into `figures/`, for example `figures/dog.txt`, and rebuild.
The filename (without `.txt`) becomes the figure name. Put `$thoughts` where
the tail connecting the figure to the balloon goes; it is replaced with `\`
when speaking and `o` when thinking. Figures are embedded in the binary with
`go:embed`.

## Layout

    main.go        the program
    figures/*.txt  the figures
