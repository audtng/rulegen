package main

	"fmt"
	"math"
	"math/big"
	"unicode"
)

	if count == 0 {
		return "", nil
	}

	return CryptoRandom(count, 0, 0, true, true)
}

/*
		if chars == nil {
			ch = rune(getCryptoRandomInt(gap) + int64(start))
		} else {
			ch = chars[getCryptoRandomInt(gap)+int64(start)]
		}

		if letters && unicode.IsLetter(ch) || numbers && unicode.IsDigit(ch) || !letters && !numbers {
package goutils

import (
	"regexp"
	"strconv"
	"testing"
	"unicode/utf8"
)
		}
	}
}

func TestCryptoRandAlphaNumeric_FuzzOnlyNumeric(t *testing.T) {

	// Testing for a reported regression in which some versions produced
	// a predictably small set of chars.
	iters := 1000
	charlen := 0
	for i := 0; i < 16; i++ {
		numOnly := 0
		charlen++
		for i := 0; i < iters; i++ {
			out, err := CryptoRandomAlphaNumeric(charlen)
			println(out)
			if err != nil {
				t.Fatal("func failed to produce a random thinger")
			}
			if _, err := strconv.Atoi(out); err == nil {
				numOnly++
			}

			m, err := regexp.MatchString("^[0-9a-zA-Z]+$", out)
			if err != nil {
				t.Fatal(err)
			}
			if !m {
				t.Fatal("Character is not alphanum")
			}
		}

		if numOnly == iters {
			t.Fatalf("Got %d numeric-only random sequences", numOnly)
		}
	}

}
	"fmt"
	"math"
	"math/rand"
	"time"
	"unicode"
)

/*
RandomAlphabetic creates a random string whose length is the number of characters specified.
Characters will be chosen from the set of alphabetic characters.

Parameters:
	count - the length of random string to create

Returns:
	string - the random string
	if err != nil {
		return "", fmt.Errorf("Error: %s", err)
	}

	return RandomString[:count], nil

}

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
)

	// H_I;E
	// 2b2ca
}

func TestRandAlphaNumeric_FuzzOnlyNumeric(t *testing.T) {

	// Testing for a reported regression in which some versions produced
	// a predictably small set of chars.
	iters := 1000
	charlen := 0
	for i := 0; i < 16; i++ {
		numOnly := 0
		charlen++
		for i := 0; i < iters; i++ {
			out, err := RandomAlphaNumeric(charlen)
			println(out)
			if err != nil {
				t.Fatal("func failed to produce a random thinger")
			}
			if _, err := strconv.Atoi(out); err == nil {
				numOnly++
			}

			m, err := regexp.MatchString("^[0-9a-zA-Z]+$", out)
			if err != nil {
				t.Fatal(err)
			}
			if !m {
				t.Fatal("Character is not alphanum")
			}
		}

		if numOnly == iters {
			t.Fatalf("Got %d numeric-only random sequences", numOnly)
		}
	}

}
