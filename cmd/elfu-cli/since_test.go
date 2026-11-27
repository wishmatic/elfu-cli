package main

import (
	"strings"
	"testing"
)

func TestSinceReportsAPastTimestamp(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "since", "2026-11-27T07:00:00+10:00")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	for _, want := range []string{
		"Human: 3 hours ago",
		"At: 2026-11-27T07:00:00+10:00",
		"Seconds: 10800",
		"Now: 2026-11-27T10:00:00+10:00",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout does not contain %q:\n%s", want, out)
		}
	}
}

func TestSinceReportsAFutureTimestamp(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "since", "2026-11-27T13:00:00+10:00")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out, "Human: in 3 hours") {
		t.Errorf("stdout = %q, want a future duration", out)
	}

	if !strings.Contains(out, "Seconds: -10800") {
		t.Errorf("stdout = %q, want negative seconds", out)
	}
}

func TestSinceReadsAZonelessTimestampInTheLocalZone(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "since", "2026-11-27T07:00:00")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out, "Human: 3 hours ago") {
		t.Errorf("stdout = %q, want the local zone reading", out)
	}

	if !strings.Contains(out, "At: 2026-11-27T07:00:00+10:00") {
		t.Errorf("stdout = %q, want the timestamp in the local zone", out)
	}
}

func TestSinceRejectsAnUnparseableTimestamp(t *testing.T) {
	out, _, err := runCommand(t, pinnedDeps(), "since", "yesterday")
	if err == nil {
		t.Fatal("Execute() error = nil, want a parse failure")
	}

	if !strings.Contains(err.Error(), "not an ISO 8601 timestamp") {
		t.Errorf("error = %v, want a parse failure", err)
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}
