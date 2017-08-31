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
	charClass  string
	noNewline  bool
	verbose    bool
	version    bool
	randSrc    string
)

func init() {
	flag.UintVarP(&length, "length", "l", 64, "length of resulting password")
	flag.BoolVarP(&upper, "upper", "U", false, "include uppercase letters A-Z")
	flag.BoolVarP(&lower, "lower", "L", false, "include lowercase letters a-z")
	flag.BoolVarP(&digits, "digits", "D", false, "include digits 0-9")
	flag.BoolVarP(&special, "special", "S", false, "include special characters, which are the printable ascii characters excluding letters, digits, and the space")
	flag.BoolVarP(&additional, "additional", "A", false, "read additional characters from standard input, encoded in UTF-8; newline characters will NOT be ignored")
	flag.StringVarP(&charClass, "type", "t", "", "Use a predefined character set.")
	flag.BoolVarP(&noNewline, "no-newline", "n", false, "Don't print a newline after the password")
	flag.BoolVarP(&verbose, "verbose", "v", false, "print more information, in addition to the generated password")
	flag.BoolVarP(&version, "version", "V", false, "print the version and exit")
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
Character sets available with --type: base64|base58|url|hex|octal|binary

If the included characters are not enough,
use -A, and pass your favorite foreign characters or emojis into standard input.

Duplicate characters are not allowed in the final alphabet

Examples:
    %s -l 32 -LUDS
        generates a password with length 32, including
        lowercase and uppercase letters, digits, and special characters.
    echo -n ABCDEF | %s -DA
        generates a 64 digit hexadecimal string (256 bits).
    %s -t hex
        does the same thing as above.

`, os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	if version {
		fmt.Printf("nextpass version %s\n", nextpass.Version)
		os.Exit(0)
	}

	var alphabet []rune
	switch charClass {
	case "base64":
		alphabet = []rune(nextpass.Base64Chars)
	case "base58":
		alphabet = []rune(nextpass.Base58Chars)
	case "url":
		alphabet = []rune(nextpass.URLChars)
	case "hex":
		alphabet = []rune(nextpass.HexChars)
	case "octal":
		alphabet = []rune(nextpass.OctalChars)
	case "binary":
		alphabet = []rune(nextpass.BinaryChars)
	case "":
	default:
		fmt.Fprintf(os.Stderr, "Unknown character set name: %s\n", charClass)
		os.Exit(1)
	}

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
		usage()
		os.Exit(1)
	}

	g, err := nextpass.NewGenerator(alphabet, int(length))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new password generator: %v\n", err)
		os.Exit(1)
	}

	if randSrc != "" {
		r, err := os.Open(randSrc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open %s: %v\n", randSrc, err)
			os.Exit(1)
		}
		defer r.Close()
		g.SetRandomSource(r)
	}

	password, n, err := g.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot generate password: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf(`password length: %d
alphabet size: %d
complexity in bits: about %d
bytes read: %d

`, length, len(alphabet), g.Bits(), n)
	}

	fmt.Print(password)
	if !noNewline {
		fmt.Println("")
	}
}
