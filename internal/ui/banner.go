package ui

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// catArt is Gary. Please be nice to him.
var catArt = []string{
	`  /\_/\  `,
	` ( o.o ) `,
	`  > ^ <  `,
	` /  |  \ `,
	`(___|___)`,
}

// catMoods are the cat-like noises and behaviors Gary greets you with.
// One is picked at random each startup.
var catMoods = []string{
	"*purrs*",
	"*knocks your mug off the desk*",
	"mrrrp?",
	"*sits on the keyboard*",
	"meow.",
	"*headbutts your hand*",
	"*stares at nothing*",
	"*kneads the stack trace*",
	"prrrt!",
	"*blinks slowly at you*",
	"*demands to be fed*",
	"*chirps at a bird*",
	"mrow?",
	"*naps in the sunbeam*",
	"*sharpens claws on the sofa*",
	"hsssss",
	"*brings you a dead pointer*",
	"*zoomies at 3am*",
}

// tagline returns the lines printed to the right of the cat.
func tagline(name string) []string {
	return []string{
		"",
		colorBold + name + colorReset,
		colorDim + catMoods[rand.IntN(len(catMoods))] + colorReset,
		"",
		colorDim + "ctrl-c to quit" + colorReset,
	}
}

// Banner returns the startup splash: an ASCII cat alongside the agent's name.
func Banner(name string) string {
	out := "\n"
	lines := tagline(name)
	for i, line := range catArt {
		right := ""
		if i < len(lines) {
			right = lines[i]
		}
		row := fmt.Sprintf("%s%s%s   %s", colorCat, line, colorReset, right)
		out += strings.TrimRight(row, " ") + "\n"
	}
	return out + "\n"
}
