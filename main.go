package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-martini/martini"
)

func main() {
	fsp := NewFileSystemPoint("share", os.Args[1], os.TempDir())
	fsp.readableUsers[""] = true
	fsp.readableUsers["admin"] = true
	Mount("share", fsp)
	//	fmt.Println(mapping[0])
	m := martini.Classic()
	m.Get("/rpc/version", HandleVersion)
	m.Post("/rpc/version", HandleVersion)
	m.Get("/rpc/ls/**", HandleLs)
	m.Get("/rpc/cat/**", HandleCat)
	m.Get("/rpc/user_config", HandleUserConfig)
	m.Run()
}

func HandleLs(res http.ResponseWriter, params martini.Params) {
	path := ResolvePath(params["_1"])
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

func HandleCat(res http.ResponseWriter, req *http.Request, params martini.Params) {
	path := ResolvePath(params["_1"])
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

func HandleUserConfig() string {
	return `{"webaxs_version":"3.0-Lite", "lang":"ja", "name":":anonymous"}`
}

func HandleVersion() string {
	return "3.1"
}
