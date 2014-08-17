package main

type StatInfo struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	IsDir      bool   `json:"directory"`
	IsWritable bool   `json:"writable"`
	Size       int64  `json:"size"`
	ATime      int64  `json:"atime"`
	MTime      int64  `json:"mtime"`
	CTime      int64  `json:"ctime"`
}
