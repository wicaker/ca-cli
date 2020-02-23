package domain

// Fs domain /
type Fs struct {
	FileName string
	DirName  string
}

// FsService /
type FsService interface {
	FindDir(dirName string) (interface{}, error)
	FindFile(fileName string) (interface{}, error)
	CreateDir(dirName string) error
	CreateFile(fileName string) error
	RemoveDir(dirName string) error
	RemoveFile(fileName string) error
	RenameDir(dirName string) error
	RenameFile(fileName string) error
}
