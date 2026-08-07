// Command useragent parses User-Agent strings and Client Hints.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "useragent: %v\n", err)
		}
		os.Exit(exitCode(err))
	}
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	var ue *usageError
	if errors.As(err, &ue) {
		return 2
	}
	return 1
}

type usageError struct {
	msg string
}

func (e *usageError) Error() string { return e.msg }

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printRootUsage(stderr)
		return &usageError{msg: "missing subcommand"}
	}
	switch args[0] {
	case "parse":
		return runParse(args[1:], stdin, stdout, stderr)
	case "-h", "-help", "--help", "help":
		printRootUsage(stderr)
		return flag.ErrHelp
	default:
		printRootUsage(stderr)
		return &usageError{
			msg: fmt.Sprintf("unknown subcommand %q", args[0]),
		}
	}
}

func printRootUsage(w io.Writer) {
	fmt.Fprint(w, `useragent analyzes User-Agent strings and Client Hints.

Usage:
  useragent parse [flags] <user-agent>
  useragent parse [flags] -          # read HTTP headers from stdin

Examples:
  useragent parse 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...'
  useragent parse -sec-ch-ua-platform '"Windows"' \
    -sec-ch-ua-platform-version '"15.0.0"' \
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...'
  printf 'User-Agent: Mozilla/5.0\nSec-CH-UA-Platform: "Android"\n\n' \
    | useragent parse -

`)
}
