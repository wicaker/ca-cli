package domain

// Fs domain /
type Fs struct {
	FileName string
	DirName  string
}

// FsService /
type FsService interface {
	FindDir(dirName string) (interface{}, error)
	FindFile(fileName string) error
	CreateDir(dirName string) error
	CreateFile() error
	RemoveDir() error
	RemoveFile() error
	RenameDir() error
	RenameFile() error
}
