package utils

import (
	"bytes"
	"github.com/alexmullins/zip"
	"github.com/pygrum/siphon/internal/db"
	"io"
	"os"
)

func ZipFile(sample *db.Sample, pass string) (string, error) {
	contents, err := os.ReadFile(sample.Path)
	if err != nil {
		return "", err
	}
	zipName := sample.Hash + ".zip"
	zipFile, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	w, err := zipWriter.Encrypt(sample.Name, pass)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		return "", err
	}
	return zipName, nil
}
