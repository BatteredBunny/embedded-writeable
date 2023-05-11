package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

var LastDate string
var LastTime string

//go:embed *
var fsdir embed.FS

func CopyEmbedToTemp(fs embed.FS, tempdir string) {
	files, err := fs.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var data []byte
		data, err = fsdir.ReadFile(file.Name())
		if err != nil {
			log.Panic(err)
		}

		if err = os.WriteFile(path.Join(tempdir, file.Name()), data, 0600); err != nil {
			log.Panic(err)
		}
	}
}

func main() {
	tempdir, err := os.MkdirTemp("", "*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	CopyEmbedToTemp(fsdir, tempdir)

	if LastDate == "" || LastTime == "" {
		fmt.Println("This is the first time you ran me!")
	} else {
		fmt.Println("Last ran at:", LastDate, LastTime)
	}

	now := time.Now()
	command := fmt.Sprintf(
		"go build -ldflags '-X main.LastDate=%s -X main.LastTime=%s' -o %s %s/*.go",
		now.Format("2006/01/01"),
		now.Format("15:04:05"),
		os.Args[0], // Replaces currently ran binary with the new one
		tempdir,
	)

	if err = exec.Command("sh", "-c", command).Run(); err != nil {
		log.Println(err)
	}
}
