package archive

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func GetArchiveConfig() *ArchiveConfig {
	var config = ArchiveConfig{
		ArchiveFileName: "content",
		DestinationDir: ".",
		ExclusionList: []string{},		
	}

	flag.Func("s", "`source directory`, is the directory to be archived", func(s string) error {
		path, err := isFieldEmpty(s, "source")
		if err != nil {
			return err
		}
	
		if !doesPathExist(path) {
			return errors.New("source " + path + " does not exist")
		}
		
		config.SourceDir = path
		return nil
	})

	flag.Func("d", "`destination directory` is the directory where archive will be stored", func(s string) error {
		path, err := isFieldEmpty(s, "destination")
		if err != nil {
			return err
		}

		if !doesPathExist(path) {
			return errors.New("destination " + path + " does not exist")
		}

		config.DestinationDir = path
		return nil
	})

	flag.Func("f", "`archive filename` is the name of the archive", func(s string) error {
		filename, err := isFieldEmpty(s, "archive filename")
		if err != nil {
			return err
		}

		config.ArchiveFileName = filename
		return nil
	})

	flag.Func("e", "`exclusion list`, is a comma separated list of directories/files to be excluded from the archive", func(s string) error {
		config.ExclusionList = strings.Split(s, ",")
		return nil
	})

	flag.Parse()

	if config.SourceDir == "" {
		fmt.Println("source cannot be set empty")
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	return &config
}

func doesPathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isFieldEmpty(path string, field string) (string, error) {
	path = strings.Trim(path, " ")
	if path == "" {
		return "", errors.New(field + " cannot be set empty")
	}

	return path, nil
}