package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("offset is negative")
	ErrNegativeLimit         = errors.New("limit is negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkParameters(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	stat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	fileSize := stat.Size()
	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	err = copyChunk(dstFile, srcFile, offset, limit)
	if err != nil {
		return err
	}

	err = dstFile.Close()
	if err != nil {
		return err
	}

	err = srcFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func copyChunk(dstFile, srcFile *os.File, offset, limit int64) error {
	var copied int64 = 0

	_, err := srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	progressor := getProgressor(50)
	progressor(0)

	for copied < limit {
		copiedChunkSize, err := io.CopyN(dstFile, srcFile, 10)

		copied += copiedChunkSize
		progressor(int(copied * 100 / limit))

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}
	}

	fmt.Println()

	return nil
}

func getProgressor(loadingSymbolsCount int) func(int) {
	const additionalSymbolsCount = 7
	const colorGreen = "\033[32m"
	const colorGrey = "\033[37m"
	const colorReset = "\033[0m"

	isFirst := true

	return func(percent int) {
		// При условии, что во время загрузки не будет ничего печататься в консоль
		if !isFirst {
			fmt.Print(strings.Repeat("\b", loadingSymbolsCount+additionalSymbolsCount))
		}

		isFirst = false
		progress := percent * loadingSymbolsCount / 100

		start := strings.Repeat("|", progress)
		end := strings.Repeat("-", loadingSymbolsCount-progress)

		fmt.Printf("[%s%s%s%s%s] %3d%%", colorGreen, start, colorGrey, end, colorReset, percent)
	}
}

func checkParameters(fromPath, toPath string, offset, limit int64) error {
	srcExt := filepath.Ext(fromPath)
	dstExt := filepath.Ext(toPath)

	if srcExt != ".txt" || dstExt != ".txt" {
		return ErrUnsupportedFile
	}

	if offset < 0 {
		return ErrNegativeOffset
	}

	if limit < 0 {
		return ErrNegativeLimit
	}

	return nil
}
