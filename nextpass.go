// Copyright (c) 2017 Niko Carpenter
// Use of this source code is governed by the MIT License,
// which can be found in the LICENSE file.

// Package nextpass provides a cryptographically secure password generator.
package nextpass

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/pkg/errors"
)

const (
	LowerChars   = "abcdefghijklmnopqrstuvwxyz"
	UpperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DigitChars   = "0123456789"
	SpecialChars = "`" + `~!@#$%^&*()-=_+[]{}\|;:'"/?<>,.`
)

var Version = "unset"

// A Generator contains options to generate passwords.
type Generator struct {
	alphabet []rune
	length   int
	source   io.Reader
}

// NewGenerator creates a Generator from the given alphabet and length.
// The default source of random bytes, crypto/rand.Reader will be used.
func NewGenerator(alphabet []rune, length int) Generator {
	return Generator{alphabet, length, rand.Reader}
}

// SetRandomSource changes the source of entropy when generating a password.
// By default, this is crypto/rand.Reader.
// Changing this to a non random source is not secure. Only do this if you know what you're doing.
func (g *Generator) SetRandomSource(source io.Reader) {
	g.source = source
}

// Max returns the total number of password combinations possible for this Generator.
func (g Generator) Max() *big.Int {
	return big.NewInt(0).Exp(big.NewInt(int64(len(g.alphabet))), big.NewInt(int64(g.length)), nil)
}

// bits returns the complexity in bits
// of a password generated with this Generator.
// If the number of bits is not a whole number, the returned value
// will be rounded up.
func (g Generator) Bits() int {
	return big.NewInt(0).Sub(g.Max(), big.NewInt(1)).BitLen()
}

// Length returns the length of generated passwords.
func (g Generator) Length() int {
	return g.length
}

// Alphabet returns the set of characters that may be included in generated passwords.
func (g Generator) Alphabet() []rune {
	return g.alphabet
}

// Generate generates a password.
// The password is returned, along with the number of bytes read from the source of entropy.
func (g Generator) Generate() (password string, n int, err error) {
	if g.length == 0 {
		return "", 0, nil
	}
	if len(g.alphabet) == 0 {
		return "", 0, errors.New("Alphabet has length 0")
	}

	// By reading from g.source once for the entire password,
	// it is possible to consume less entropy.
	base := big.NewInt(int64(len(g.alphabet)))
	r := newReadCounter(g.source)
	num, err := rand.Int(r, g.Max())
	if err != nil {
		return "", r.count, errors.Wrap(err, "Cannot get random data")
	}

	m := big.NewInt(0)
	passChars := make([]rune, g.length)
	// Add characters in reverse, so that the encoded characters
	// follow the same order as the bytes read from g.source.
	for i := g.length - 1; i >= 0; i-- {
		num.DivMod(num, base, m)
		passChars[i] = g.alphabet[int(m.Int64())]
	}

	return string(passChars), r.count, nil
}
