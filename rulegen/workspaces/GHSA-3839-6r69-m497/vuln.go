package main

	"fmt"
	"math"
	"math/big"
	"regexp"
	"unicode"
)

	if count == 0 {
		return "", nil
	}
	RandomString, err := CryptoRandom(count, 0, 0, true, true)
	if err != nil {
		return "", fmt.Errorf("Error: %s", err)
	}
	match, err := regexp.MatchString("([0-9]+)", RandomString)
	if err != nil {
		panic(err)
	}

	if !match {
		//Get the position between 0 and the length of the string-1  to insert a random number
		position := getCryptoRandomInt(count)
		//Insert a random number between [0-9] in the position
		RandomString = RandomString[:position] + string('0' + getCryptoRandomInt(10)) + RandomString[position + 1:]
		return RandomString, err
	}
	return RandomString, err

}

/*
		if chars == nil {
			ch = rune(getCryptoRandomInt(gap) + int64(start))
		} else {
			ch = chars[getCryptoRandomInt(gap) + int64(start)]
		}

		if letters && unicode.IsLetter(ch) || numbers && unicode.IsDigit(ch) || !letters && !numbers {
package goutils

import (
	"testing"
	"unicode/utf8"
)
		}
	}
}
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"time"
	"unicode"
)

/*
RandomAlphabetic creates a random string whose length is the number of characters specified.
Characters will be chosen from the set of alpha-numeric characters as indicated by the arguments.

Parameters:
	count - the length of random string to create
	letters - if true, generated string may include alphabetic characters
	numbers - if true, generated string may include numeric characters

Returns:
	string - the random string
	if err != nil {
		return "", fmt.Errorf("Error: %s", err)
	}
	match, err := regexp.MatchString("([0-9]+)", RandomString)
	if err != nil {
		panic(err)
	}

	if !match {
		//Get the position between 0 and the length of the string-1  to insert a random number
		position := rand.Intn(count)
		//Insert a random number between [0-9] in the position
		RandomString = RandomString[:position] + string('0'+rand.Intn(10)) + RandomString[position+1:]
		return RandomString, err
	}
	return RandomString, err

}

import (
	"fmt"
	"math/rand"
	"testing"
)

	// H_I;E
	// 2b2ca
}
