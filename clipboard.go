package main

import "github.com/atotto/clipboard"

func ReadAll() (string, error) {
	return clipboard.ReadAll()
}
