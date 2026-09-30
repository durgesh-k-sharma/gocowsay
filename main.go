package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

//go:embed figures/*.txt
var figuresFS embed.FS

func figureNames() []string {
	entries, _ := figuresFS.ReadDir("figures") // returned sorted by name
	var names []string
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".txt"))
	}
	return names
}

// loadFigure returns the art for a figure. The $thoughts placeholder is
// replaced by the tail character: "\" for speech, "o" for thought.
func loadFigure(name, tail string) (string, error) {
	data, err := figuresFS.ReadFile("figures/" + name + ".txt")
	if err != nil {
		return "", fmt.Errorf("unknown figure %q (available: %s)",
			name, strings.Join(figureNames(), ", "))
	}
	art := strings.TrimRight(string(data), "\n")
	return strings.ReplaceAll(art, "$thoughts", tail), nil
}

// wrap splits a line into pieces of at most width runes, breaking on
// spaces. Words longer than width are split mid-word. Lines that already
// fit are returned untouched, so indentation is preserved.
func wrap(line string, width int) []string {
	if utf8.RuneCountInString(line) <= width {
		return []string{line}
	}

	var out []string
	cur := ""
	for _, word := range strings.Fields(line) {
		r := []rune(word)
		for len(r) > width {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			out = append(out, string(r[:width]))
			r = r[width:]
		}
		word = string(r)

		switch {
		case cur == "":
			cur = word
		case utf8.RuneCountInString(cur)+1+len(r) <= width:
			cur += " " + word
		default:
			out = append(out, cur)
			cur = word
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

// buildBalloon wraps already-padded lines (all `width` runes wide) in a balloon.
func buildBalloon(lines []string, width int, think bool) string {
	out := []string{" " + strings.Repeat("_", width+2)}
	last := len(lines) - 1

	for i, l := range lines {
		var left, right string
		switch {
		case think:
			left, right = "(", ")"
		case last == 0:
			left, right = "<", ">"
		case i == 0:
			left, right = "/", "\\"
		case i == last:
			left, right = "\\", "/"
		default:
			left, right = "|", "|"
		}
		out = append(out, fmt.Sprintf("%s %s %s", left, l, right))
	}

	out = append(out, " "+strings.Repeat("-", width+2))
	return strings.Join(out, "\n")
}

// readMessage takes the message from arguments if given, otherwise from a pipe.
func readMessage(args []string) ([]string, error) {
	if len(args) > 0 {
		return strings.Split(strings.Join(args, " "), "\n"), nil
	}

	info, err := os.Stdin.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice != 0 {
		return nil, fmt.Errorf("no input: pass a message or pipe one in")
	}

	var lines []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		lines = append(lines, strings.TrimSuffix(sc.Text(), "\r"))
	}
	return lines, sc.Err()
}

func main() {
	figure := flag.String("f", "cow", "figure to use")
	list := flag.Bool("l", false, "list available figures and exit")
	think := flag.Bool("t", false, "think instead of speak")
	maxWidth := flag.Int("W", 40, "wrap the message at this many columns")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: gocowsay [flags] [message...]")
		fmt.Fprintln(os.Stderr, "       fortune | gocowsay [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *list {
		fmt.Println(strings.Join(figureNames(), "\n"))
		return
	}
	if *maxWidth < 1 {
		fmt.Fprintln(os.Stderr, "-W must be at least 1")
		os.Exit(1)
	}

	tail := "\\"
	if *think {
		tail = "o"
	}
	art, err := loadFigure(*figure, tail)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	raw, err := readMessage(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	// Expand tabs (one rune, many columns) and wrap long lines.
	var lines []string
	for _, l := range raw {
		l = strings.ReplaceAll(l, "\t", "    ")
		lines = append(lines, wrap(l, *maxWidth)...)
	}
	if len(lines) == 0 {
		return
	}

	// Pad every line to the widest one so the right border lines up.
	width := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > width {
			width = n
		}
	}
	for i, l := range lines {
		lines[i] = l + strings.Repeat(" ", width-utf8.RuneCountInString(l))
	}

	fmt.Println(buildBalloon(lines, width, *think))
	fmt.Println(art)
}
