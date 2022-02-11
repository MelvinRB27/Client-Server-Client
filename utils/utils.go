package utils

import (
	"math/rand"
	"time"
)

//esto es para genera un Id de número aleatorio
func GenRandomID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}