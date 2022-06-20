package main

import (
	"fmt"
	"github.com/AgeroFlynn/crud/foundation/config"
)

func main() {
	err := config.NewConfigFromFile()
	if err != nil {
		fmt.Println(err)
	}
}
