package utils

import (
	"errors"
	"math/rand"
)

type Container[K any] interface {
	Contains(key K) bool
}

var TooManyIdCollisions = errors.New("too many ids collided when trying to generate a new one")

const characters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Id(length int) string {
	var id = make([]byte, length)
	for i := range length {
		id[i] = characters[rand.Intn(len(characters))]
	}
	return string(id)
}

func UniqueId(length int, used Container[string]) (string, error) {
	for range 50 {
		var id = Id(length)
		if !used.Contains(id) {
			return id, nil
		}
	}
	return "", TooManyIdCollisions
}
