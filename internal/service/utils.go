package service

import (
	"path"
	"path/filepath"
)

func picDri(id string, small bool) string {
	p, _ := filepath.Abs("./runtime/pic")
	if small {
		return path.Join(p, "small", id+".png")
	}
	return path.Join(p, "original", id+".png")
}
