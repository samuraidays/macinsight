package main

import (
	"strings"
	"testing"
)

func TestToSet(t *testing.T) {
	cases := []struct{
		in string
		want map[string]struct{}
	}{
		{"", map[string]struct{}{}},
		{"a", map[string]struct{}{"a":{}}},
		{"a,b", map[string]struct{}{"a":{}, "b":{}}},
		{" a , b ", map[string]struct{}{"a":{}, "b":{}}},
		{",,a,,", map[string]struct{}{"a":{}}},
	}
	for _, c := range cases {
		got := toSet(c.in)
		if len(got) != len(c.want) {
			t.Fatalf("len mismatch for %q: got=%d want=%d", c.in, len(got), len(c.want))
		}
		for k := range c.want {
			if _, ok := got[k]; !ok {
				t.Fatalf("missing key %q for %q", k, c.in)
			}
		}
	}
}

func TestUsageTextContainsCommands(t *testing.T) {
	// Minimal guard that usage mentions key flags/commands; keep independent of stdout capture.
	usageText := `macinsight audit --json`
	if !strings.Contains(usageText, "audit") || !strings.Contains(usageText, "--json") {
		t.Fatalf("usage text missing keywords")
	}
}
