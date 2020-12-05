package main

import (
	"os"
	"io"
	// "fmt"
	"github.com/dexter3k/vault13/file/dat"
)

func main() {
	arc, err := dat.Load(os.Args[1])
	check(err)
	ff, err := arc.LoadFile("COLOR.PAL")
	check(err)
	defer ff.Close()

	f, err := os.Create("colorz.pal")
	check(err)
	defer f.Close()

	_, err = io.Copy(f, ff)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
