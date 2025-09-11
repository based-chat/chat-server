// main - entrypoint for chat-server application
package main

import (
	"log"

	"github.com/fatih/color"
)

// main is the entrypoint for chat-server application.
// It prints a green "Hello, World!" to the console.
func main() {
	log.Println(color.GreenString("Hello, World!"))
}
