package main

import "github.com/adelrosarioh/tarantulas/collectors/textbooks"

func main() {

	if err := textbooks.Run(); err != nil {
		panic(err)
	}
}
