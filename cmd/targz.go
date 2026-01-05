package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func targzComp(fromPath *string, toPath *string) error {
	slog.SetDefault(logger)

	return compressF(fromPath, toPath)

}

func compressF(fromPath *string, toPath *string) error {
	*fromPath = strings.TrimSuffix(*fromPath, "/")
	*toPath = strings.TrimSuffix(*toPath, "/")

	slog.Info("logts - started process", "src", *fromPath, "des", *toPath)

	output, err := os.Create(*toPath + "/" + archiveName(fromPath))
	if err != nil {
		return err
	}
	defer output.Close()

	if err := createArchive(fromPath, output); err != nil {
		return err
	}
	slog.Info("logts - ended process - success", "src", *fromPath, "des", *toPath)

	return nil
}

func createArchive(fromPath *string, buf io.Writer) error {
	gzipw := gzip.NewWriter(buf)
	defer gzipw.Close()
	tarw := tar.NewWriter(gzipw)
	defer tarw.Close()

	filesDir, err := os.ReadDir(*fromPath)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, file := range filesDir {
		wg.Add(1)
		go addFile(file, fromPath, tarw, &wg, &mu)
	}

	wg.Wait()

	return nil
}

func addFile(entry os.DirEntry, filePath *string, tarw *tar.Writer, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	fullFilePath := filepath.Join(*filePath, entry.Name())

	slog.Info("logts - started file", "file", fullFilePath)

	file, err := os.Open(fullFilePath)
	if err != nil {
		slog.Error("logts - error processing file 1", "file", fullFilePath, "errmsg", err.Error())
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		slog.Error("logts - error processing file 2", "file", fullFilePath, "errmsg", err.Error())
		return
	}
	relPath, err := filepath.Rel(fromPath, fullFilePath)
	if err != nil {
		slog.Error("logts - error processing file 1", "file", fullFilePath, "errmsg", err.Error())
		return
	}
	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())

	if err != nil {
		slog.Error("logts - error processing file 3", "file", fullFilePath, "errmsg", err.Error())
		return
	}

	header.Name = filepath.ToSlash(relPath)

	mu.Lock()
	defer mu.Unlock()

	if err := tarw.WriteHeader(header); err != nil {
		slog.Error("logts - error processing file 4", "file", fullFilePath, "errmsg", err.Error())
		return
	}

	if _, err := io.Copy(tarw, file); err != nil {
		slog.Error("logts - error processing file 5", "file", fullFilePath, "errmsg", err.Error())
		return
	}
	slog.Info("logts - finished processing file, success", "file", fullFilePath)
}

func archiveName(path *string) string {
	newPath := strings.Split(*path, "/")
	date := time.Now().Format(time.DateOnly)
	times := time.Now().Format(time.TimeOnly)

	return fmt.Sprintf("%s_%s-%s-%s.tar.gz",
		"logts",
		newPath[len(newPath)-1],
		strings.ReplaceAll(date, "-", ""),
		strings.ReplaceAll(times, ":", ""),
	)
}
