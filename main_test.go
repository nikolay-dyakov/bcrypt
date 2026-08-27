package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestRunHashesPasswordFromStdin(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := run([]string{"-cost", "10"}, strings.NewReader("secret\n"), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}

	hash := strings.TrimSpace(stdout.String())
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("secret")); err != nil {
		t.Fatalf("generated hash does not match the password: %v", err)
	}
}

func TestRunAcceptsExplicitEmptyPassword(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{"-p", "", "-c", "10"}, strings.NewReader("ignored"), &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}

	hash := strings.TrimSpace(stdout.String())
	if err := bcrypt.CompareHashAndPassword([]byte(hash), nil); err != nil {
		t.Fatalf("generated hash does not match the empty password: %v", err)
	}
}

func TestRunRejectsUnsafeCost(t *testing.T) {
	err := run([]string{"-p", "secret", "-c", "9"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "cost must be between 10 and 31") {
		t.Fatalf("expected an invalid-cost error, got %v", err)
	}
}

func TestRunRejectsLongPassword(t *testing.T) {
	password := strings.Repeat("a", maxPasswordBytes+1)
	err := run([]string{"-p", password, "-c", "10"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "at most 72 bytes") {
		t.Fatalf("expected a password-length error, got %v", err)
	}
}

func TestRunPrintsVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "command", args: []string{"version"}},
		{name: "single-dash flag", args: []string{"-version"}},
		{name: "double-dash flag", args: []string{"--version"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := run(test.args, strings.NewReader("must not be read"), &stdout, &bytes.Buffer{})
			if err != nil {
				t.Fatalf("run returned an error: %v", err)
			}
			if got, want := stdout.String(), "bcrypt v"+appVersion+"\n"; got != want {
				t.Fatalf("version output = %q, want %q", got, want)
			}
		})
	}
}
