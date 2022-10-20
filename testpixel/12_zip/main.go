package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
)

type Opener interface {
	Open(name string) (fs.File, error)
}

type Loader struct {
	opener Opener
	root   string
}

type Option func(*Loader)

func WithZip(name string) Option {
	return func(l *Loader) {
		l.opener, _ = zip.OpenReader(name)
	}
}

func NewLoader(root string, opts ...Option) *Loader {
	l := &Loader{
		root: root,
	}
	for _, opt := range opts {
		opt(l)
	}

	if l.opener == nil {
		l.opener = os.DirFS(root)
		l.root = ""
	}

	return l
}

func (l *Loader) Read(name string) (io.Reader, error) {
	return l.opener.Open(fmt.Sprintf("%s%s", l.root, name))
}

func main() {
	ld := NewLoader("assets/", WithZip("example.zip"))

	ld2 := NewLoader("assets/")
	mapPath := "my3.tmx"
	r, _ := ld.Read(mapPath)
	fmt.Println("zip opener", r)

	r2, _ := ld2.Read(mapPath)
	fmt.Println("dir opener", r2)
}

// func main() {
// 	reader, err := zip.OpenReader("example.zip")
// 	if err != nil {
// 		msg := "Failed to open: %s"
// 		log.Fatalf(msg, err)
// 	}
// 	defer reader.Close()

// 	dirPath := "assets"
// 	mapPath := "my3.tmx"

// 	r, _ := reader.Open(fmt.Sprintf("%s/%s", dirPath, mapPath))
// 	fmt.Println("zip opener", r)

// 	reader2 := os.DirFS("assets")
// 	r, _ = reader2.Open(mapPath)
// 	fmt.Println("dir opener", r)
// }
