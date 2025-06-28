package main

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// Downloads the requested file from file server in current executeable folder.
func downloadFile(fileName, filepath string) error {
	// Since we want to resume the download - create file with append mode.
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		slog.Error("Error opening file:", "file", err)
		return err
	}
	defer file.Close()

	// Get the current file size for resume
	fileInfo, err := file.Stat()
	if err != nil {
		slog.Error("Error getting file info:", "file", err)
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", fileServer+fileName, nil)
	if err != nil {
		slog.Error("Error creating request:", "request", err)
		return err
	}

	// Add Range header for resume capability
	if fileInfo.Size() > 0 {
		req.Header.Add("Range", "bytes="+fmt.Sprint(fileInfo.Size())+"-")
	}

	// Perform our request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error performing request/download:", "request", err)
		return err
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength, "downloading game")

	// Check server supports resume - should be
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		slog.Error("Server does not support resume:", "status", resp.StatusCode)
		return err
	}

	// Create progress tracking
	written, err := io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		slog.Error("Error copying file:", "file", err)
		return err
	}
	slog.Info("File downloaded:", "written", written)
	return nil
}

// Checks if fileName (located in executeable directory) matches with given checksum.
func verifyChecksum(fileName string, checksum string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("Error opening file for checksum verification:", "file", err)
		return false
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.Error("Error calculating checksum:", "error", err)
		return false
	}

	calculatedChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	if calculatedChecksum != checksum {
		slog.Error("Checksum mismatch:", "expected", checksum, "got", calculatedChecksum)
		return false
	}

	slog.Info("Checksum verified successfully")
	return true
}

// Tries to extract the given file in the current executeable folder.
func extractZip(fileName string) error {
	// Open the zip file
	reader, err := zip.OpenReader(fileName)
	if err != nil {
		slog.Error("Error opening zip file:", "file", err)
		return err
	}
	defer reader.Close()

	// Calculate total size for progress bar
	var totalSize int64
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	// Create progress bar
	bar := progressbar.DefaultBytes(
		totalSize,
		"Extracing",
	)

	// Iterate over each file in the archive.
	for _, file := range reader.File {
		// Open the file inside zip
		rc, err := file.Open()
		if err != nil {
			slog.Error("Error opening file:", "file", err)
			return err
		}
		defer rc.Close()

		// Create the directory structure, if needed.
		path := filepath.Join(".", file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		// Create the file
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		outFile, err := os.Create(path)
		if err != nil {
			slog.Error("Error creating file:", "file", err)
			return err
		}
		defer outFile.Close()

		writer := io.MultiWriter(outFile, bar)

		// Copy the contents
		_, err = io.Copy(writer, rc)
		if err != nil {
			slog.Error("Error copying file:", "file", err)
			return err
		}
	}
	return nil
}
