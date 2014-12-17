package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	fsp := NewFileSystemPoint("share", os.Args[1], os.TempDir())
	fsp.readableUsers[""] = true
	fsp.readableUsers["admin"] = true
	Mount("share", fsp)
	m := mux.NewRouter()
	m.HandleFunc("/rpc/version", HandleVersion)
	m.HandleFunc(`/rpc/ls{path:/[^?]*}`, HandleLs)
	m.HandleFunc(`/rpc/cat{path:/[^?]*}`, HandleCat)
	m.HandleFunc(`/rpc/thumbnail{path:/.*}`, HandleThumbnail)
	m.HandleFunc("/rpc/user_config", HandleUserConfig)
	n := negroni.Classic()
	n.UseHandler(m)
	n.Run(":3000")
}

func HandleLs(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	path := ResolvePath(vars["path"])
	if path == nil || !path.Exists() {
		res.WriteHeader(404)
		fmt.Fprint(res, "File not found")
	} else if path.IsDir() {
		res.Header().Add("Content-Type", "application/json")
		json, _ := json.Marshal(path.Children(NewUser("admin")))
		res.Write(json)
	} else {
		res.Header().Add("Content-Type", "application/json")
		json, _ := json.Marshal(path.StatInfo(NewUser("admin")))
		res.Write(json)
	}
}

func HandleCat(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	path := ResolvePath(vars["path"])
	if path == nil {
		res.WriteHeader(404)
	} else {
		rs, err := path.Reader()
		if err != nil {
			res.WriteHeader(404)
		} else {
			http.ServeContent(res, req, path.Clean(), path.ModTime(), rs)
		}
	}
}

func HandleUserConfig(res http.ResponseWriter, _ *http.Request) {
	io.WriteString(res, `{"webaxs_version":"3.0-Lite", "lang":"ja", "name":":anonymous"}`)
}

func HandleVersion(res http.ResponseWriter, _ *http.Request) {
	io.WriteString(res, "3.1")
}
