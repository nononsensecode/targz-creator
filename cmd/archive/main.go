package main

import (
	"fmt"
	"time"

	"github.com/mholt/archiver/v3"
)

func main() {
	start := time.Now()
	archiveDetails := GetArchiveDetails()
	fileList := archiveDetails.FileList()
	err := archiver.Archive(fileList, archiveDetails.TarGz())
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Archived in %s\n", elapsed)
}