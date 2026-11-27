package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestTimePrintsThePinnedInstant(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "time")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	for _, want := range []string{
		"2026-11-27T10:00:00+10:00",
		"AEST +10:00",
		"Unix " + strconv.FormatInt(testNow.Unix(), 10),
		"Friday, 27 November 2026 at 10:00 AEST",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout does not contain %q:\n%s", want, out)
		}
	}

	if lines := strings.Count(strings.TrimSpace(out), "\n"); lines != 0 {
		t.Errorf("stdout is not one line:\n%s", out)
	}
}

func TestTimeRejectsArguments(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "time", "extra")
	if err == nil {
		t.Fatal("Execute() error = nil, want an argument rejection")
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}
