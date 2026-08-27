package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var appVersion string

const (
	sourceVersion     = "1.0.0"
	defaultCost       = 13
	minimumSecureCost = 10
	maxPasswordBytes  = 72
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "bcrypt: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("bcrypt", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var password string
	var cost int
	var copyToClipboard bool
	var showVersion bool

	flags.StringVar(&password, "p", "", "password to hash (prefer the hidden prompt or stdin)")
	flags.StringVar(&password, "password", "", "password to hash (prefer the hidden prompt or stdin)")
	flags.IntVar(&cost, "c", defaultCost, "bcrypt cost (10-31)")
	flags.IntVar(&cost, "cost", defaultCost, "bcrypt cost (10-31)")
	flags.BoolVar(&copyToClipboard, "cc", false, "copy the bcrypt hash to the clipboard")
	flags.BoolVar(&copyToClipboard, "clipboard", false, "copy the bcrypt hash to the clipboard")
	flags.BoolVar(&showVersion, "version", false, "print version information")

	if err := flags.Parse(args); err != nil {
		return err
	}
	versionCommand := flags.NArg() == 1 && flags.Arg(0) == "version"
	if showVersion || versionCommand {
		if showVersion && flags.NArg() != 0 {
			return fmt.Errorf("the version flag does not accept arguments")
		}
		_, err := fmt.Fprintf(stdout, "bcrypt v%s\n", currentVersion())
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %q", flags.Args())
	}
	if cost < minimumSecureCost || cost > bcrypt.MaxCost {
		return fmt.Errorf("cost must be between %d and %d", minimumSecureCost, bcrypt.MaxCost)
	}

	passwordWasSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "p" || f.Name == "password" {
			passwordWasSet = true
		}
	})

	passwordBytes := []byte(password)
	if !passwordWasSet {
		var err error
		passwordBytes, err = readPassword(stdin, stderr)
		if err != nil {
			return err
		}
	}
	defer clear(passwordBytes)

	if len(passwordBytes) > maxPasswordBytes {
		return fmt.Errorf("password is %d bytes; bcrypt accepts at most %d bytes", len(passwordBytes), maxPasswordBytes)
	}

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, cost)
	if err != nil {
		return fmt.Errorf("generate hash: %w", err)
	}

	if copyToClipboard {
		if err := clipboard.WriteAll(string(hash)); err != nil {
			return fmt.Errorf("copy hash to clipboard: %w", err)
		}
		fmt.Fprintln(stderr, "bcrypt hash copied to clipboard")
		return nil
	}

	_, err = fmt.Fprintln(stdout, string(hash))
	return err
}

func currentVersion() string {
	if appVersion != "" {
		return strings.TrimPrefix(appVersion, "v")
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		version := info.Main.Version
		if version != "" && version != "(devel)" && !strings.HasPrefix(version, "v0.0.0-") && !strings.Contains(version, "+dirty") {
			return strings.TrimPrefix(version, "v")
		}
	}
	return sourceVersion
}

func readPassword(stdin io.Reader, stderr io.Writer) ([]byte, error) {
	if input, ok := stdin.(*os.File); ok && term.IsTerminal(int(input.Fd())) {
		fmt.Fprint(stderr, "Password: ")
		password, err := term.ReadPassword(int(input.Fd()))
		fmt.Fprintln(stderr)
		if err != nil {
			return nil, fmt.Errorf("read password: %w", err)
		}
		return password, nil
	}

	password, err := io.ReadAll(io.LimitReader(stdin, maxPasswordBytes+3))
	if err != nil {
		return nil, fmt.Errorf("read password from stdin: %w", err)
	}
	password = bytes.TrimSuffix(password, []byte("\n"))
	password = bytes.TrimSuffix(password, []byte("\r"))
	return password, nil
}
