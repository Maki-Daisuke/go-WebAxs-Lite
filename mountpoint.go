package main

import (
	"errors"
	"io"
	"path/filepath"
	"sort"
	"time"
)

// Path abstracts file path treated by WebAxs.
type Path interface {
	String() string // Returns original path
	Clean() string  // Returns cleaned-up path
	Exists() bool
	IsDir() bool
	IsRegular() bool
	IsReadableBy(User) bool
	IsWritableBy(User) bool
	IsCreatableBy(User) bool
	ModTime() time.Time
	Create(io.Reader) error
	MkdirAll() error
	RemoveAll() error
	Reader() (io.ReadSeeker, error)
	Writer() (io.ReadWriteSeeker, error)
	Move(Path) error
	Copy(Path) error
	//	TmpDir() string
	StatInfo(u User) *StatInfo   // Returns StatInfo from perspective of User u
	Children(u User) []*StatInfo // List of Child nodes. Returns nil if this is not a directory
}

// MountPoint abstracts file system that WebAxs handles. It accepts path string
// and returns abstract Path object correnponding to the path.
// path      : whole path string to resolve (already cleaned-up)
// mountPoint: path string where this MountPoint is mounted
// original  : whole path string without cleaning-up
type MountPoint interface {
	Path(path, mountPoint, original string) Path
}

type mappingItemT struct {
	path       string
	mountPoint MountPoint
}

type mappingT []*mappingItemT

var mapping = mappingT(make([]*mappingItemT, 0, 10))

func (m mappingT) Len() int {
	return len(m)
}

func (m mappingT) Less(i, j int) bool {
	return len(m[i].path) >= len(m[j].path)
}

func (m mappingT) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Mount mounts physical path at specified URL path
func Mount(path string, mp MountPoint) {
	path = CleanPath(path)
	mapping = append(mapping, &mappingItemT{path, mp})
	sort.Sort(mapping)
}

// ResolvePath convert a path string into Path object.
// It returns nil when it cannot resolve the path with any of mounted MountPoints.
func ResolvePath(path string) Path {
	clean := CleanPath(path)
	if clean == "/" { // Returns root (/) directory
		return NewRootPath(path)
	}
	for _, m := range mapping {
		if matchPrefix(m.path, clean) {
			return m.mountPoint.Path(clean, m.path, path)
		}
	}
	return nil // Kind of error, where no mapping is found
}

func matchPrefix(prefix, str string) bool {
	if len(prefix) > len(str) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if prefix[i] != str[i] {
			return false
		}
	}
	return true
}

type RootPath string

func NewRootPath(orig string) *RootPath {
	var r RootPath = RootPath(orig)
	return &r
}

func (r *RootPath) String() string {
	return string(*r)
}

func (r *RootPath) Clean() string {
	return "/"
}

func (r *RootPath) Exists() bool {
	return true
}

func (r *RootPath) IsDir() bool {
	return true
}

func (r *RootPath) IsRegular() bool {
	return false
}

func (r *RootPath) IsReadableBy(_ User) bool {
	return true
}

func (r *RootPath) IsWritableBy(_ User) bool {
	return false
}

func (r *RootPath) IsCreatableBy(_ User) bool {
	return false
}

func (r *RootPath) ModTime() time.Time {
	return time.Now()
}

func (r *RootPath) Create(_ io.Reader) error {
	return errors.New("Permission denied") // TODO: return appropriate error
}

func (r *RootPath) MkdirAll() error {
	return errors.New("Permission denied") // TODO: return appropriate error
}

func (r *RootPath) RemoveAll() error {
	return errors.New("Permission denied") // TODO: return appropriate error
}

func (r *RootPath) Reader() (io.ReadSeeker, error) {
	return nil, errors.New("") // TODO: return appropriate error
}

func (r *RootPath) Writer() (io.ReadWriteSeeker, error) {
	return nil, errors.New("") // TODO: return appropriate error
}

func (r *RootPath) Move(_ Path) error {
	return errors.New("") // TODO: return appropriate error
}

func (r *RootPath) Copy(Path) error {
	return errors.New("") // TODO: return appropriate error
}

func (r *RootPath) StatInfo(_ User) *StatInfo {
	return &StatInfo{
		Name:       filepath.Base(r.String()),
		Path:       "/",
		IsDir:      true,
		IsWritable: false,
		Size:       0,
		ATime:      time.Now().Unix(),
		MTime:      time.Now().Unix(),
		CTime:      time.Now().Unix(),
	}
}

func (r *RootPath) Children(user User) []*StatInfo {
	ls := make([]*StatInfo, 0, len(mapping))
	for _, v := range mapping {
		p := ResolvePath(v.path)
		if !p.IsReadableBy(user) {
			continue
		}
		ls = append(ls, p.StatInfo(user))
	}
	return ls
}
