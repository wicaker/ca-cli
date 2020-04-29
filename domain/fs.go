package domain

import "os"

// FsService /
type FsService interface {
	FindDir(dirName string) (interface{}, error)
	FindFile(fileName string) (interface{}, error)
	CreateDir(dirName string) error
	RemoveDir(dirName string) error
	ReadDir(dirName string) ([]os.FileInfo, error)
}
