package admin

import (
	"embed"
	"io/fs"
)

//go:embed all:dist fallback
var embedded embed.FS

func FS() fs.FS {
	dist, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	if _, err := fs.Stat(dist, "index.html"); err == nil {
		return dist
	}
	fallback, err := fs.Sub(embedded, "fallback")
	if err != nil {
		panic(err)
	}
	return fallback
}
