package main

import (
	"crypto/rand"
	"fmt"
	"github.com/russross/blackfriday"
)

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func convertMarkDown(s string) string {
	return string(blackfriday.MarkdownBasic([]byte(s)))
}
