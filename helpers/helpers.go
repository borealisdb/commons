package helpers

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/big"
)

var passwordChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		maxN := big.NewInt(int64(len(passwordChars)))
		if n, err := cryptoRand.Int(cryptoRand.Reader, maxN); err != nil {
			panic(fmt.Errorf("Unable to generate secure, random password: %v", err))
		} else {
			b[i] = passwordChars[n.Int64()]
		}
	}
	return string(b)
}
