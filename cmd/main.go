package main

import (
	"MSRM/internal/api"
	"fmt"
)

func main() {
	fmt.Println("Hi!")
	api.StartServer()
	fmt.Println("Goddbye!")
}
