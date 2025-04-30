package generator

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	LowerCaseChars    = "abcdefghijklmnopqrstuvwxyz"
	UpperCaseChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberChars       = "0123456789"
	AlphabetChars     = LowerCaseChars + UpperCaseChars
	AlphaNumericChars = AlphabetChars + NumberChars
)

type CharSet struct {
	letterBytes         string
	bitMask             byte
	availableCharLength int
}

func New(charSet string) (*CharSet, error) {
	availableCharLength := len(charSet)
	if availableCharLength < 2 || availableCharLength > 256 {
		return nil, errors.New("availableCharBytes length must be greater" +
			" than 0 and less than or equal to 256")
	}

	var bitLength byte
	var bitMask byte
	for bits := availableCharLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1
	return &CharSet{
		letterBytes:         charSet,
		bitMask:             bitMask,
		availableCharLength: availableCharLength,
	}, nil
}

func (g *CharSet) RandomBytes(length int) ([]byte, error) {
	bufferSize := length + length/3
	var err error
	result := make([]byte, length)
	for i, j, rb := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			rb, err = randomBytes(bufferSize)
			if err != nil {
				return nil, fmt.Errorf("unable to generate secure random bytes: %s", err)
			}
		}

		if idx := int(rb[j%length] & g.bitMask); idx < g.availableCharLength {
			result[i] = g.letterBytes[idx]
			i++
		}
	}
	return result, nil
}

func randomBytes(length int) ([]byte, error) {
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	if err != nil {
		return nil, err
	}
	return rb, nil
}
