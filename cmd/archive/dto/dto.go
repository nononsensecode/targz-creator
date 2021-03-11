package dto

import (
	"path/filepath"
	"strings"
	"time"
)

// ArchiveDetails describes name of the archive, dir to be archived and destination directory
type ArchiveDetails struct {
	ArchiveDir string
	ArchiveName string
	DestinationDir string
	ExclusionList []string
}

// DestinationFilename returns the targz file name including destination directory
func (a *ArchiveDetails) DestinationFilename() string {
	return filepath.Join(a.DestinationDir, a.targz())
}

func (a *ArchiveDetails) targz() string {
	currentTime := time.Now()
	return a.ArchiveName + "-" + currentTime.Format("02-01-2006-15-04-05") + ".tar.gz"
}

// ArchiveDirTrim removes all trailing and leading spaces from the source dir name
func (a *ArchiveDetails) ArchiveDirTrim() string {
	return strings.Trim(a.ArchiveDir, " ")
}