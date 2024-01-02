package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"errors"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

const macKeyLen = 128

var ErrNoMatch = errors.New("password does not match")

type SaltedSHA512PBKDF2Dictionary struct {
	Iterations int    `plist:"iterations"`
	Salt       []byte `plist:"salt"`
	Entropy    []byte `plist:"entropy"`
}

func SaltedSHA512PBKDF2(plaintext string) (SaltedSHA512PBKDF2Dictionary, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return SaltedSHA512PBKDF2Dictionary{}, err
	}
	iterations, err := secureRandInt(20000, 40000)
	if err != nil {
		return SaltedSHA512PBKDF2Dictionary{}, err
	}
	return SaltedSHA512PBKDF2Dictionary{
		Iterations: iterations,
		Salt:       salt,
		Entropy: pbkdf2.Key([]byte(plaintext),
			salt, iterations, macKeyLen, sha512.New),
	}, nil
}

func Verify(plaintext string, h SaltedSHA512PBKDF2Dictionary) error {
	hashed := pbkdf2.Key([]byte(plaintext), h.Salt, h.Iterations, macKeyLen, sha512.New)
	if 1 != subtle.ConstantTimeCompare(h.Entropy, hashed) {
		return ErrNoMatch
	}
	return nil
}

func secureRandInt(min, max int64) (int, error) {
	var random int
	for {
		iter, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return 0, err
		}
		if iter.Int64() >= min {
			random = int(iter.Int64())
			break
		}
	}
	return random, nil
}
