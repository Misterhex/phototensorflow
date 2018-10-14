package main

import (
	"io"
	"log"
	"os"
	"os/exec"

	"path"
	"path/filepath"

	"github.com/google/uuid"
)

func ClassifyImage(r io.Reader) (string, error) {
	fileName := uuid.New().String()
	f, err := os.Create(fileName)

	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, r)

	if err != nil {
		return "", err
	}

	pyscripts := path.Join(getExecutingDirectory(), "pyscripts/classify_image.py")
	imageAbs := path.Join(getExecutingDirectory(), fileName)

	c := "/usr/bin/python"
	args := []string{pyscripts, "--image_file", imageAbs}

	log.Printf("running %s %s\n", c, args)

	cmd := exec.Command(c, args...)

	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}

	log.Println("output: " + string(out))

	return string(out), nil
}

func getExecutingDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
