package main

import (
	"embed"
	"log"
	"os"
	"path"
)

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

func PrintError(input ...any) {
	log.Println("\033[91m", input, "\033[0m")
}
