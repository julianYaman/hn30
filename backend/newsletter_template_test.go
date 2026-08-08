package main

import "testing"

func TestSourceLabel(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"https://github.com/xoreaxeaxeax/rosenbridge", "github.com/xoreaxeaxeax/rosenbridge"},
		{"https://github.com/xoreaxeaxeax/rosenbridge/blob/main/README.md", "github.com/xoreaxeaxeax/rosenbridge"},
		{"https://www.runnersworld.com/news/a73355106/hamster/", "runnersworld.com"},
		{"https://0xkrt26.github.io/math/", "0xkrt26.github.io"},
		{"https://news.ycombinator.com/item?id=1", "news.ycombinator.com"},
	}
	for _, tt := range tests {
		if got := sourceLabel(tt.in); got != tt.want {
			t.Errorf("sourceLabel(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
