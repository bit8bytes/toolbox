package main

import "github.com/bit8bytes/toolbox/hello"

func main() {
	geeting := hello.Greet("Tobi!")
	println(geeting)

	purpose := hello.Purpose
	println(purpose)
}
