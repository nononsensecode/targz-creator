package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ArchiveDetails struct {
	ArchiveDir string
	DestinationDir string
	ArchiveFilename string
	ExclusionList []string
}

func (a *ArchiveDetails) Destination() string {
	return path.Join(a.DestinationDir, a.ArchiveFilename)
}

func (a *ArchiveDetails) TarGz() string {
	return a.Destination() + ".tar.gz"
}

func (a *ArchiveDetails) FileList() []string {
	var fileList []string
	filepath.Walk(a.ArchiveDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Do not include the root directory in the list
		if a.ArchiveDir == path {
			return nil
		}

		for _, filename := range a.ExclusionList {
			if isDir(filename) {
				if info.IsDir() && info.Name() == strings.Trim(filename, "/") {
					return filepath.SkipDir
				}
			} else {
				if info.Name() == filename {
					return nil
				}
			}
		}

		// Add first level directory only, don't walk too deep
		if info.IsDir() {
			fileList = append(fileList, path)
			return filepath.SkipDir
		}

		fileList = append(fileList, path)

		return nil		
	})

	return fileList
}

func isDir(path string) bool {
	return strings.HasSuffix(path, "/")
}

func GetArchiveDetails() *ArchiveDetails {
	var archiveDetails = ArchiveDetails{
		DestinationDir: ".",
		ArchiveFilename: "content",
		ExclusionList: []string{},
	}

	flag.Func("s", "`source directory`", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("source directory cannot be set empty")
		}

		err := doesExist(s)
		if err != nil {
			return err
		}

		archiveDetails.ArchiveDir = s
		return nil
	})

	flag.Func("d", "`destination directory`", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("destination directory cannot be set empty")
		}

		err := doesExist(s)
		if err != nil {
			return err
		}

		archiveDetails.DestinationDir = s
		return nil
	})

	flag.Func("f", "`archiv file name`", func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("archive filename cannot be set empty")
		}

		archiveDetails.ArchiveFilename = s
		return nil
	})

	flag.Func("e", "`list of filenames or dirs` separated by comma", func(s string) error {
		archiveDetails.ExclusionList = strings.Split(s, ",")
		return nil
	})

	flag.Parse()

	if archiveDetails.ArchiveDir == "" {
		fmt.Println("source directory is not set")
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	return &archiveDetails
}

func doesExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New("file '" + path + "' does not exist")
	}

	return nil
}