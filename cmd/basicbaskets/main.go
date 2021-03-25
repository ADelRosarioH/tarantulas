package main

import (
	basicbasket "github.com/adelrosarioh/tarantulas/collectors/basicbaskets"
)

func main() {

	if err := basicbasket.Run(); err != nil {
		panic(err)
	}
}
