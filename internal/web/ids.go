package web

import (
	"errors"
	"math/rand"
)

var TooManyIdCollisions = errors.New("too many ids collided when trying to generate a new one")

const characters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Id(length int) string {
	var id = make([]byte, length)
	for i := range length {
		id[i] = characters[rand.Intn(len(characters))]
	}
	return string(id)
}

func UniqueId[T any](length int, used map[string]T) (string, error) {
	for range 50 {
		var id = Id(length)
		if _, ok := used[id]; ok == false {
			return id, nil
		}
	}
	return "", TooManyIdCollisions
}
