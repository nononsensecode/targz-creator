package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"nononsensecode.com/backup-creator/internal/pkg/archive"
)

func main() {
	start := time.Now()
	config := archive.GetArchiveConfig()

	// Create the archive file
	archive, err := os.Create(config.TarGz())
	if err != nil {
		log.Fatalf("Error creating archive %s\n %v\n", config.TarGz(), err)
	}
	defer archive.Close()

	// Add files to archive
	err = config.CreateArchive(archive)
	if err != nil {
		fmt.Println("Error writing to archive", err)
		debug.PrintStack()
		os.Exit(2)
	}

	elapsed := time.Since(start)
	fmt.Printf("Archive created successfully in %v\n", elapsed)
}

