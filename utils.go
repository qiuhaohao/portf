package main

import "os"

func MustOpenFile(name string) *os.File {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return f
}
