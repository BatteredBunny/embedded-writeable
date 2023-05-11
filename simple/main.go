package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
)

const dataSignifier = "uwunya"

type FileWrapper struct {
	bs []byte
}

func newFileWrapper(fileName string) (f FileWrapper, err error) {
	var bs []byte
	bs, err = os.ReadFile(fileName)
	if err != nil {
		return
	}

	f = FileWrapper{bs: bs}
	return
}

func (f FileWrapper) containsSignifier() bool {
	return bytes.Contains(f.bs, []byte(dataSignifier))
}

func (f FileWrapper) getData() []byte {
	if !f.containsSignifier() {
		return []byte{}
	}

	split := bytes.Split(f.bs, []byte(dataSignifier))

	if len(split) <= 2 {
		return []byte{}
	} else {
		return split[2]
	}
}

func (f *FileWrapper) addData(data string) {
	f.bs = append(f.bs, dataSignifier+data...)
}

func (f *FileWrapper) setData(data string) {
	if f.containsSignifier() {
		dataWithSignifier := append([]byte(dataSignifier), f.getData()...)
		f.bs, _ = bytes.CutSuffix(f.bs, dataWithSignifier)
	}

	f.addData(data)
}

func (f FileWrapper) writeFile(fileName string, perm fs.FileMode) error {
	return os.WriteFile(fileName, f.bs, perm)
}

func main() {
	data := flag.String("data", "", "data to add at end of binary")
	flag.Parse()

	fileName := os.Args[0]
	f, err := newFileWrapper(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if *data != "" {
		info, err := os.Stat(fileName)
		if err != nil {
			log.Fatal(err)
		}

		f.setData(*data)

		if err = os.Remove(fileName); err != nil {
			log.Fatal(err)
		}

		if err := f.writeFile(fileName, info.Mode()); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data:")
	fmt.Println(string(f.getData()))
}
