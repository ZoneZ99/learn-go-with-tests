package main

import (
	"os"
	"time"
	"./clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}