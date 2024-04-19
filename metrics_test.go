package main

import (
	"strings"
	"testing"

	"github.com/mihailtudos/chirpy/pkg/utils"
)

func TestIsWordInSlice(t *testing.T) {
	cases := []struct {
		Word  string
		Words []string
		Found bool
	}{
		{
			Word:  "Doku",
			Words: strings.Fields("I don't know this player."),
			Found: false,
		},
		{
			Word:  "Doku",
			Words: strings.Fields("I heared something about doku is it a player?."),
			Found: true,
		},
		{
			Word:  "Doku",
			Words: strings.Fields("Doku! Doku! Doku!."),
			Found: false,
		},
	}

	for i := range cases {
		if utils.InSlice(cases[i].Word, cases[i].Words) != cases[i].Found {
			t.Errorf("'%s' not found in %v", cases[i].Word, cases[i].Words)
		}
	}
}

func TestProfineWordsAreRemoved(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "This is a kerfuffle opinion I need to share with the world",
			Expected: "This is a **** opinion I need to share with the world",
		},
		{
			Input:    "I had something interesting for breakfast",
			Expected: "I had something interesting for breakfast",
		},
		{
			Input:    "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			Expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			Input:    "I hear Mastodon is better than Chirpy. Sharbert! I need to migrate",
			Expected: "I hear Mastodon is better than Chirpy. Sharbert! I need to migrate",
		},
	}

	for _, c := range cases {
		res := utils.ReplaceProfaneWords(c.Input)
		if res != c.Expected {
			t.Errorf("given: '%s'\nexpected: '%s'\ngot: '%s'\n", c.Input, c.Expected, res)
			continue
		}
	}
}
