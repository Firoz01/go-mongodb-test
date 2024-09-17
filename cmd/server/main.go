package main

import (
	"github.com/Firoz01/go-mongodb-test/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application := app.NewApplication()

	application.Init()
	application.Run()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal (like Ctrl+C)
	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	application.Cleanup()

	// Exit the program explicitly
	log.Println("Program terminated successfully")
	os.Exit(0)
}
