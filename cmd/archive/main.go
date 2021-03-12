package main

import "github.com/mholt/archiver/v3"

func main() {
	archiveDetails := GetArchiveDetails()
	fileList := archiveDetails.FileList()
	err := archiver.Archive(fileList, archiveDetails.TarGz())
	if err != nil {
		panic(err)
	}
}