package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
)

type Cli struct {
	Shift struct {
		Count int
	} `cmd:""`

	Keyword struct {
		Keyword string
	} `cmd:""`

	Atbash struct {} `cmd:""`

	Display bool
	Word string
}

type Cipher map[int]int

func makeShift(cli *Cli) Cipher {
	c := Cipher{}
	for i := 0; i < 26; i++ {
		j := (i + cli.Shift.Count) % 26
		c[j] = i
	}
	return c
}

func (cipher Cipher) contains(c int) bool {
	for v := range cipher {
		if v == c {
			return true
		}
	}
	return false
}

func makeKeyword(cli *Cli) Cipher {
	cipher := Cipher{}
	w := []rune(cli.Keyword.Keyword)
	j := 0
	for _, c := range w {
		ci := int(c - 'a')
		if cipher.contains(ci) {
			continue
		}
		fmt.Printf("Adding %d to %d\n", j, ci)
		cipher[ci] = j
		j++
	}
	for i := 0; i < 26; i++ {
		if cipher.contains(i) {
			continue
		}
		fmt.Printf("adding %d to %d\n", j, i)
		cipher[i] = j
		j++
	}
	return cipher
}

func makeAtbash(cli *Cli) Cipher {
	cipher := Cipher{}
	for i := 0; i < 26; i++ {
		cipher[i] = 25-i
	}
	return cipher
}

func (cipher Cipher) display() {
	for i := 0; i < 26; i++ {
		fmt.Printf("%c ", rune(i+int('a')))
	}
	fmt.Printf("\n")
	for i := 0; i < 26; i++ {
		fmt.Printf("%c ", rune(cipher[i]+int('a')))
	}
	fmt.Printf("\n")
}

func (cipher Cipher) encode(s string) string {
	r := ""
	for _, a := range []rune(s) {
		i := int(a-'a')
		j, found := cipher[i]
		if found {
			r += string(rune(j + 'a'))
		} else {
			r += string(a)
		}
	}
	return r
}

func main() {
	var cli Cli
	var cipher Cipher

	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "shift":
		cipher = makeShift(&cli)

	case "keyword":
		cipher = makeKeyword(&cli)

	case "atbash":
		cipher = makeAtbash(&cli)

	default:
		panic("unhandled")
	}

	if cli.Display {
		cipher.display()
	}

	if cli.Word != "" {
		fmt.Printf("%s\n", strings.ToLower(cipher.encode(cli.Word)))
	}
}