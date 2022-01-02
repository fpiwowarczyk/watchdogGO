package utils

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockOsFs struct {
	mock.Mock
}

func (mock *MockOsFs) Stat(name string) (os.FileInfo, error) {
	args := mock.Called(name)
	return nil, args.Error(1)
}

func (mock *MockOsFs) ExecAndOutput(name string, arg ...string) ([]byte, error) {
	args := mock.Called(name, arg)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *MockOsFs) Open(name string) (file, error) {
	args := mock.Called(name)
	return args.Get(0).(file), args.Error(1)
}
