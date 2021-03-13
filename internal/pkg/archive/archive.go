package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ArchiveConfig struct {
	SourceDir string
	ArchiveFileName string
	DestinationDir string
	ExclusionList []string
}

func (a *ArchiveConfig) ParentDir() string {
	return filepath.Dir(a.SourceDir)
}

func (a *ArchiveConfig) TarGz() string {
	return path.Join(a.DestinationDir, a.ArchiveFileName + ".tar.gz")
}

func (a *ArchiveConfig) CreateArchive(buf io.Writer) error {
	gzipWriter := gzip.NewWriter(buf)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return filepath.Walk(a.SourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

     	pathWithoutParent := strings.Trim(path, a.ParentDir())

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = pathWithoutParent

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}

		return nil
	})
}