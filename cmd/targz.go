package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Job struct {
	FullPath string
	RelPath  string
	FInfo    fs.FileInfo
}
type Result struct {
	File   *os.File
	Header *tar.Header
	Err    error
}

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func targzComp(fromPath *string, toPath *string) error {
	slog.SetDefault(logger)

	return compressF(fromPath, toPath)

}

func compressF(fromPath *string, toPath *string) error {

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

	jobs := make(chan Job)
	results := make(chan Result)
	done := make(chan bool)

	numOfWorker := 4
	var wg sync.WaitGroup
	for range numOfWorker {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		for res := range results {
			if res.Err != nil {
				slog.Error("logts - error processing file", "errmsg", res.Err)
				continue
			}
			writeTar(res, tarw)
		}
		done <- true
	}()

	err := filepath.WalkDir(*fromPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(*fromPath, path)
		if err != nil {
			return err
		}
		jobs <- Job{
			FullPath: path,
			RelPath:  filepath.ToSlash(relPath),
			FInfo:    info,
		}
		return nil
	})

	close(jobs)
	wg.Wait()
	close(results)
	<-done

	return err
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		file, err := os.Open(job.FullPath)
		if err != nil {
			results <- Result{Err: err}
			continue
		}
		header, err := tar.FileInfoHeader(job.FInfo, job.FInfo.Name())
		if err != nil {
			results <- Result{Err: err}
			file.Close()
			continue
		}
		header.Name = job.RelPath
		results <- Result{
			File:   file,
			Header: header,
		}
	}
}

func writeTar(res Result, tarw *tar.Writer) {
	defer res.File.Close()

	if err := tarw.WriteHeader(res.Header); err != nil {
		slog.Error("logts - header write error", "file", res.Header.Name, "errmsg", err)
		return
	}

	if _, err := io.Copy(tarw, res.File); err != nil {
		slog.Error("logts - copy error", "file", res.Header.Name, "errmsg", err)
		return
	}
	slog.Info("logts - file archived, success", "file", res.Header.Name)
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
