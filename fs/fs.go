package fs

import (
	"github.com/wicaker/cacli/domain"

	"github.com/spf13/afero"
)

type caFs struct {
	fs afero.Fs
}

// NewFsService creates /
func NewFsService() domain.FsService {
	fs := afero.NewOsFs()
	return &caFs{
		fs,
	}
}

func (f *caFs) FindDir(dirName string) (interface{}, error) {
	res, err := afero.ReadDir(f.fs, dirName)
	if err != nil {
		return nil, nil
	}

	return res, nil
}

func (f *caFs) FindFile(fileName string) error { return nil }

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
