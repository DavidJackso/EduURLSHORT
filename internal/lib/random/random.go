package random

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, length)
	for i := range res {
		res[i] = letters[rnd.Intn(len(letters))]
	}
	return string(res)
}
