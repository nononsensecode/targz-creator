package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"nononsensecode.com/backup-creator/cmd/archive/dto"
)

func main() {
	archiveDetails := archiveDetails()

	// Create archive file
	archiveFile, err := os.Create(archiveDetails.DestinationFilename())
	if err != nil {
		log.Fatalln("Error creating archive file", err)
	}
	defer archiveFile.Close()

	gzipFile := gzip.NewWriter(archiveFile)
	defer gzipFile.Close()
	tarFile := tar.NewWriter(gzipFile)
	defer tarFile.Close()

	_, err = os.Stat(archiveDetails.ArchiveDir)
	if err != nil {
		log.Fatalf("source directory %s does not exist\n", archiveDetails.ArchiveDir)
	}

	err = filepath.Walk(archiveDetails.ArchiveDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		exclude, isDir := exclude(archiveDetails, info)
		if exclude {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}
		header.Name = path

		err = tarFile.WriteHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarFile, file)
		return err
	})
}

func exclude(archiveDetails *dto.ArchiveDetails, info fs.FileInfo) (bool, bool) {
	for _, filename := range archiveDetails.ExclusionList {
		if strings.HasSuffix(filename, "/") {
			if matchDir(filename, info) {
				return true, true
			}
		} else {
			if matchFile(filename, info) {
				return true, false
			}
		}
	}

	return false, false
}

func matchDir(dirname string, info fs.FileInfo) bool {
	dirname = strings.Trim(dirname, "/")
	if info.IsDir() && info.Name() == dirname {
		return true
	}
	return false
}

func matchFile(filename string, info fs.FileInfo) bool {
	if info.Name() == filename {
		return true
	}
	return false
}

func archiveDetails() *dto.ArchiveDetails {
	archiveDetails := dto.ArchiveDetails{
		ArchiveName: "content",
		DestinationDir: ".",
		ExclusionList: []string{},
	}

	flag.Func("f", "`filename` for the archive", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("file name cannot be set empty")
		}

		archiveDetails.ArchiveName = s
		return nil
	})

	flag.Func("s", "archive `directory name`", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("archive directory is not set")
		}

		archiveDetails.ArchiveDir = s
		return doesExist(archiveDetails.ArchiveDir)
	})

	flag.Func("d", "`directory name` where arhive to be stored", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("destination directory cannot be set empty")
		}

		archiveDetails.DestinationDir = s
		return doesExist(archiveDetails.DestinationDir)
	})

	flag.Func("e", "comma separated `list of files` or dirs that are to be excluded", func(s string) error {
		trimmed := strings.Trim(s, " ")
		archiveDetails.ExclusionList = strings.Split(trimmed, ",")
		return nil
	})

	flag.Parse()	

	if archiveDetails.ArchiveDir == "" {
		println("archive directory is not set")
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	return &archiveDetails
}


func doesExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}