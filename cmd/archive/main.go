package main

import "github.com/mholt/archiver/v3"

func main() {
	archiveDetails := GetArchiveDetails()
	err := archiver.Archive(archiveDetails.FileList(), archiveDetails.TarGz())
	if err != nil {
		panic(err)
	}
}