/*
Package fs provides methods for interacting with the filesystem
*/
package fs

import (
	"os"

	"github.com/wicaker/cacli/domain"

	"github.com/spf13/afero"
)

type caFs struct {
	fs afero.Fs
}

// NewFsService will create new a caFs object representation of domain.FsService interface
func NewFsService() domain.FsService {
	fs := afero.NewOsFs()
	return &caFs{
		fs,
	}
}

// FindDir is a method for find a directory based on directory name
func (f *caFs) FindDir(dirName string) (interface{}, error) {
	res, err := afero.ReadDir(f.fs, dirName)
	if err != nil {
		return nil, nil
	}

	return res, nil
}

// FindFile is a method for find a file based on file name
func (f *caFs) FindFile(fileName string) (interface{}, error) {
	res, err := afero.ReadFile(f.fs, fileName)
	if err != nil {
		return nil, nil
	}

	return res, nil
}

// CreateDir is a method for create a directory based on directory name
func (f *caFs) CreateDir(dirName string) error {
	err := f.fs.Mkdir(dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}

// RemoveDir is a method for remove a directory based on directory name
func (f *caFs) RemoveDir(dirName string) error {
	err := f.fs.RemoveAll(dirName)
	if err != nil {
		return err
	}
	return nil
}

// ReadDir is a method for read a directory based on directory name
func (f *caFs) ReadDir(dirName string) ([]os.FileInfo, error) {
	res, err := afero.ReadDir(f.fs, dirName)
	if err != nil {
		return nil, nil
	}

	return res, nil
}
