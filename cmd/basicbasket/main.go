package main

import (
	basicbasket "github.com/adelrosarioh/tarantulas/collectors/basickbaskets"
)

func main() {

	if err := basicbasket.Run(); err != nil {
		panic(err)
	}
}
