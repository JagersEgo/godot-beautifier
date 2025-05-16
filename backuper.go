package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mholt/archives"
)

func NewBackup(base string, input_files []string) (string, error) {
	ctx := context.TODO()

	save_location := fmt.Sprintf(
		"%s/godot-linter-backup_%d.tar.zstd",
		os.TempDir(),
		time.Now().Unix(),
	)

	filemap := make(map[string]string, len(input_files))
	for _, f := range input_files {
		filemap[f] = f
	}

	files, err := archives.FilesFromDisk(ctx, nil, filemap)
	if err != nil {
		return "ERROR", err
	}

	// create the output file we'll write to
	out, err := os.Create(save_location)
	if err != nil {
		return "ERROR", err
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (since we're writing, we only set Archival, but if you're
	// going to read, set Extraction)
	format := archives.CompressedArchive{
		Compression: archives.Zstd{},
		Archival:    archives.Tar{},
	}

	// create the archive
	err = format.Archive(ctx, out, files)
	if err != nil {
		return "ERROR", err
	}

	return save_location, nil
}
