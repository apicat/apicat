package project

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"/Users/alice/dev/shop", "-Users-alice-dev-shop"},
		{"/home/dev/my.project", "-home-dev-my-project"},
		{"/srv/api_v2", "-srv-api_v2"},
		{"/tmp/空格 dir", "-tmp----dir"},
	}
	for _, c := range cases {
		if got := Encode(c.in); got != c.want {
			t.Errorf("Encode(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDirFor(t *testing.T) {
	dir, err := DirFor("/Users/alice/dev/shop")
	if err != nil {
		t.Fatalf("DirFor: %v", err)
	}
	if filepath.Base(dir) != "-Users-alice-dev-shop" {
		t.Errorf("unexpected base: %s", filepath.Base(dir))
	}
	if !strings.Contains(dir, filepath.Join(".apicat", "projects")) {
		t.Errorf("expected path under .apicat/projects, got %s", dir)
	}
}
