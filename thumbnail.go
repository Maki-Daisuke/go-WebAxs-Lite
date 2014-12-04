package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/go-martini/martini"
)

func HandleThumbnail(res http.ResponseWriter, req *http.Request, params martini.Params) {
	path := ResolvePath(params["_1"])
	if !isImageFile(path) {
		res.WriteHeader(404)
		res.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(res, "File not found")
		return
	}
	switch file := path.(type) {
	default:
		panic("why isn't this *FilePath!?")
	case *FilePath:
		if req.URL.Query().Get("nonblocking") != "" {
			handleNonBlocking(res)
		} else {
			handleThumbnailAux(res, req, file)
		}
	}
}

func handleNonBlocking(res http.ResponseWriter) {
	res.WriteHeader(200)
}

func handleThumbnailAux(res http.ResponseWriter, req *http.Request, file *FilePath) {
	geometry := parseSize(req.URL.Query().Get("size"))
	// convert <INPUT> -resize <SIZE> -background white -gravity center -extent <SIZE> - # output to STDOUT
	cmd := exec.Command("convert", file.Physical(), "-resize", geometry, "-background", "white", "-gravity", "center", "-extent", geometry, "-format", "jpg", "-")
	stdout := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		panic("convert command failed: " + stderr.String())
	}
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "image/jpeg")
	_, err = res.Write(stdout.Bytes())
	if err != nil {
		panic(err.Error())
	}
}

var sizeAlias = map[string]string{
	"S":  "42x42",
	"M":  "85x85",
	"L":  "180x180",
	"LL": "232x232",
	"3L": "500x500",
	"4L": "1024x1024",
}

var reGeometry, _ = regexp.Compile(`^[0-9]+x[0-9]+$`)

func parseSize(size string) (geometry string) {
	if reGeometry.MatchString(size) {
		return size
	} else {
		geometry = sizeAlias[size]
		if geometry != "" {
			return geometry
		} else {
			// Unknown geometry string, use "M" as default
			return sizeAlias["M"]
		}
	}
}

var reImageFile, _ = regexp.Compile(`(?i)\.(?:jpg|jpeg|gif|bmp|png|tif)$`)

func isImageFile(path Path) bool {
	return path != nil && path.IsRegular() && reImageFile.MatchString(path.Clean())
}
