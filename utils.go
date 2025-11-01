package main

import (
	"strings"
)

func parseWords(s string) []string {
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, w := range raw {
		word := strings.TrimSpace(w)
		if word != "" {
			out = append(out, word)
		}
	}
	return out
}
