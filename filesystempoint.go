package main

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// FileSystemPoint implements MountPoint interface based on real file system.
type FileSystemPoint struct {
	name          string          // name of this shared directory
	path          string          // absolute path of this shared directory
	tmpdir        string          // temporary directory for this shared directory
	readableUsers map[string]bool // true if user can read
	writableUsers map[string]bool // true if user can write
}

// FilePath abstracts file path treated by WebAxs.
type FilePath struct {
	name       string           // base (the last part) of this path
	original   string           // original path
	clean      string           // clean path
	physical   string           // physical path on the file system corresponding to this path
	mountPoint *FileSystemPoint // MountPoint which this path belongs to
	fileInfo   os.FileInfo      // FileInfo cache
	children   []string         // child entries (nil if this path is not a directory)
}

func NewFileSystemPoint(name, path, tmpdir string) *FileSystemPoint {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	tmpdir, err = filepath.Abs(tmpdir)
	if err != nil {
		panic(err)
	}
	return &FileSystemPoint{
		name:          name,
		path:          filepath.Clean(path),
		tmpdir:        filepath.Clean(tmpdir),
		readableUsers: make(map[string]bool),
		writableUsers: make(map[string]bool),
	}
}

func (fsp *FileSystemPoint) Path(path, mountPoint, original string) Path {
	return &FilePath{
		name:       filepath.Base(original),
		original:   original,
		clean:      path,
		physical:   filepath.Join(fsp.path, path[len(mountPoint):]),
		mountPoint: fsp,
	}
}

func (fsp *FileSystemPoint) isReadableBy(u User) bool {
	return fsp.readableUsers[u.Name()]
}

func (fsp *FileSystemPoint) isWritableBy(u User) bool {
	return fsp.writableUsers[u.Name()]
}
func (p *FilePath) String() string {
	return p.original
}

func (p *FilePath) Clean() string {
	return p.clean
}

func (p *FilePath) Physical() string {
	return p.physical
}

func (p *FilePath) Exists() bool {
	if p.fileInfo == nil {
		fi, err := os.Lstat(p.physical)
		if err == nil {
			p.fileInfo = fi
		} else if os.IsNotExist(err) {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func (p *FilePath) IsDir() bool {
	if !p.Exists() {
		return false
	}
	return p.fileInfo != nil && p.fileInfo.IsDir()
}

func (p *FilePath) IsRegular() bool {
	if !p.Exists() {
		return false
	}
	return p.fileInfo != nil && p.fileInfo.Mode().IsRegular()
}

func (p *FilePath) IsReadableBy(u User) bool {
	return p.mountPoint.isReadableBy(u)
}

func (p *FilePath) IsWritableBy(u User) bool {
	return p.mountPoint.isWritableBy(u)
}

func (p *FilePath) IsCreatableBy(u User) bool {
	return p.mountPoint.isWritableBy(u)
}

func (p *FilePath) ModTime() time.Time {
	if !p.Exists() {
		return time.Time{}
	}
	return p.fileInfo.ModTime()
}

func (p *FilePath) Create(r io.Reader) error {
	panic("Not implemented yet")
}

func (p *FilePath) MkdirAll() error {
	return os.MkdirAll(p.physical, 0777)
}

func (p *FilePath) RemoveAll() error {
	panic("Not implemented yet")
}

func (p *FilePath) Reader() (io.ReadSeeker, error) {
	return os.Open(p.physical)
}

func (p *FilePath) Writer() (io.ReadWriteSeeker, error) {
	return os.OpenFile(p.physical, os.O_RDWR, 0666)
}

func (p *FilePath) Move(dst Path) error {
	panic("Not implemented yet")
}

func (p *FilePath) Copy(dst Path) error {
	panic("Not implemented yet")
}

func (p *FilePath) TmpDir() string {
	return p.mountPoint.tmpdir
}

func (p *FilePath) StatInfo(u User) *StatInfo {
	if !p.Exists() {
		return nil
	}
	return &StatInfo{
		Name:       p.name,
		Path:       p.clean,
		IsDir:      p.IsDir(),
		IsWritable: p.IsWritableBy(u),
		Size:       p.fileInfo.Size(),
		ATime:      0, // TODO: assign appropriate value
		MTime:      p.fileInfo.ModTime().Unix(),
		CTime:      0, // TODO: assign appropriate value
	}
}

func (p *FilePath) Children(u User) []*StatInfo {
	if !p.IsDir() {
		return nil
	}
	file, err := os.Open(p.physical)
	if err != nil {
		panic(err)
	}
	fis, err := file.Readdir(0)
	if err != nil && err != io.EOF {
		panic(err)
	}
	ls := make([]*StatInfo, 0, len(fis))
	for _, fi := range fis {
		ls = append(ls, &StatInfo{
			Name:       fi.Name(),
			Path:       JoinPath(p.Clean(), fi.Name()),
			IsDir:      fi.IsDir(),
			IsWritable: p.IsWritableBy(u),
			Size:       fi.Size(),
			ATime:      0, // TODO: assign appropriate value
			MTime:      fi.ModTime().Unix(),
			CTime:      0, // TODO: assign appropriate value
		})
	}
	return ls
}
