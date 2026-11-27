package main

import (
	"strings"
	"testing"
	"time"
)

func TestSleepPrintsTheWakeTime(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "sleep", "0")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	for _, want := range []string{"Seconds: 0", "Now: 2026-11-27T10:00:00+10:00"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout does not contain %q:\n%s", want, out)
		}
	}
}

func TestSleepWaitsTheRequestedDuration(t *testing.T) {
	start := time.Now()

	if _, _, err := runCommand(t, pinnedDeps(), "sleep", "0.05"); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
		t.Errorf("elapsed = %v, want at least 50ms", elapsed)
	}
}

func TestSleepRefusesMoreThanThirtySeconds(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "sleep", "31")
	if err == nil || !strings.Contains(err.Error(), "at most 30") {
		t.Errorf("Execute() error = %v, want a cap rejection", err)
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}

func TestSleepRefusesNegativeSeconds(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "sleep", "--", "-1")
	if err == nil || !strings.Contains(err.Error(), "zero or more") {
		t.Errorf("Execute() error = %v, want a negative rejection", err)
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}

func TestSleepRefusesANonNumber(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "sleep", "soon")
	if err == nil || !strings.Contains(err.Error(), "not a number") {
		t.Errorf("Execute() error = %v, want a number rejection", err)
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}
