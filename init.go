package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		os.Chdir(path)
	}
	rand.Seed(time.Now().UnixNano())
}
