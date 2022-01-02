package utils

// Mocking file system as in here https://talks.golang.org/2012/10things.slide#8 for testing

import (
	"io"
	"os"
	"os/exec"
)

var fs FileSystem = OsFS{}

type FileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	ExecAndOutput(name string, arg ...string) ([]byte, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type OsFS struct{}

func (OsFS) Open(name string) (file, error) {
	return os.Open(name)
}

func (OsFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsFS) ExecAndOutput(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}
