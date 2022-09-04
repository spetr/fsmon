package main

import "math/rand"

func getRandString(n int) string {
	var r = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = r[rand.Intn(len(r))]
	}
	return string(s)
}
