package generator

import (
	"net/http"
	"os"
)

type Binary struct {
	Asset      func(name string) ([]byte, error)
	AssetInfo  func(name string) (os.FileInfo, error)
	AssetNames func() []string
	AssetFile  func() http.FileSystem
	Gzip       bool
}
