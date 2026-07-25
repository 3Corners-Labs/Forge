// Package workbench exposes the compiled exploratory Forge Workbench client.
// It contains presentation assets only; canonical delivery state remains owned
// by the capability APIs projected by the CLI coordinator.
package workbench

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var compiled embed.FS

func Assets() fs.FS {
	assets, err := fs.Sub(compiled, "dist")
	if err != nil {
		panic(err)
	}
	return assets
}
