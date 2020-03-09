/*
Package fs provides methods for interacting with the filesystem
*/
package fs

import (
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

func (f *caFs) FindFile(fileName string) (interface{}, error) {
	res, err := afero.ReadFile(f.fs, fileName)
	if err != nil {
		return nil, nil
	}

	return res, nil
}

func (f *caFs) CreateDir(dirName string) error {
	err := f.fs.Mkdir(dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (f *caFs) CreateFile(fileName string) error { return nil }
func (f *caFs) RemoveDir(dirName string) error {
	err := f.fs.RemoveAll(dirName)
	if err != nil {
		return err
	}
	return nil
}
func (f *caFs) RemoveFile(fileName string) error { return nil }
func (f *caFs) RenameDir(dirName string) error   { return nil }
func (f *caFs) RenameFile(fileName string) error { return nil }
