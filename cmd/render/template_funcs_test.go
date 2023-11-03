package main

import (
	"fmt"
	"testing"
)

func TestToSentence(t *testing.T) {
	tests := []struct {
		items []string
		want  string
	}{
		{
			items: []string{},
			want:  "",
		},
		{
			items: []string{"Earth"},
			want:  "Earth",
		},
		{
			items: []string{"Earth", "Wind"},
			want:  "Earth and Wind",
		},
		{
			items: []string{"Earth", "Wind", "Fire"},
			want:  "Earth, Wind and Fire",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("slice with %d items", i), func(t *testing.T) {
			got := toSentence(tc.items)
			if got != tc.want {
				t.Errorf("Want %q, got %q", tc.want, got)
			}
		})
	}
}
