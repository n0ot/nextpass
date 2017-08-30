// Copyright (c) 2017 Niko Carpenter
// Use of this source code is governed by the MIT License,
// which can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/n0ot/nextpass"
	flag "github.com/spf13/pflag"
)

var (
	length     uint
	lower      bool
	upper      bool
	digits     bool
	special    bool
	additional bool
	noNewline  bool
	verbose    bool
	randSrc    string
)

func init() {
	flag.UintVarP(&length, "length", "l", 64, "length of resulting password")
	flag.BoolVarP(&upper, "upper", "U", false, "include uppercase letters A-Z")
	flag.BoolVarP(&lower, "lower", "L", false, "include lowercase letters a-z")
	flag.BoolVarP(&digits, "digits", "D", false, "include digits 0-9")
	flag.BoolVarP(&special, "special", "S", false, "include special characters, which are the printable ascii characters excluding letters, digits, and the space")
	flag.BoolVarP(&additional, "additional", "A", false, "read additional characters from standard input, encoded in UTF-8; newline characters will NOT be ignored")
	flag.BoolVarP(&noNewline, "no-newline", "n", false, "Don't print a newline after the password")
	flag.BoolVarP(&verbose, "verbose", "v", false, "print more information, in addition to the generated password")
	flag.StringVarP(&randSrc, "random-source", "r", "", "specify a file to be used as an alternate source of randomness. Don't use this unless you know what you're doing.")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [options]

nextpass generates a cryptographically random password. Whenever you need to
create your next password, use nextpass.

`, os.Args[0])

	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
If the included characters are not enough,
use -A, and pass your favorite foreign characters or emojis into standard input.

Examples:
    %s -l 32 -LUDS
        generates a password with length 32, including
        lowercase and uppercase letters, digits, and special characters.
    echo -n ABCDEF | %s -DA
        generates a 64 digit hexadecimal string (256 bits).

`, os.Args[0], os.Args[0])
}

func main() {
	var alphabet []rune
	if lower {
		alphabet = append(alphabet, []rune(nextpass.LowerChars)...)
	}
	if upper {
		alphabet = append(alphabet, []rune(nextpass.UpperChars)...)
	}
	if digits {
		alphabet = append(alphabet, []rune(nextpass.DigitChars)...)
	}
	if special {
		alphabet = append(alphabet, []rune(nextpass.SpecialChars)...)
	}
	if additional {
		additionalChars, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read from standard input. If additional characters were passed in,\nthey will be skipped.\n\n%v\n", err)
		}
		alphabet = append(alphabet, []rune(string(additionalChars))...)
	}

	if len(alphabet) == 0 {
		fmt.Fprintln(os.Stderr, "No characters included in password; cannot generate.\nDid you forget to enable one of the character types?")
		os.Exit(1)
	}

	g := nextpass.NewGenerator(alphabet, int(length))
	if randSrc != "" {
		r, err := os.Open(randSrc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open %s: %v\n", randSrc, err)
			os.Exit(1)
		}
		defer r.Close()
		g.SetRandomSource(r)
	}
	if verbose {
		fmt.Printf(`password length: %d
alphabet size: %d
complexity in bits: about %d

`, length, len(alphabet), g.Bits())
	}

	password, err := g.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot generate password: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(password)
	if !noNewline {
		fmt.Println("")
	}
}
