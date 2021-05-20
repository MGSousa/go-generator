package generator

import "os"

type Binary struct {
	Asset      func(name string) ([]byte, error)
	AssetInfo  func(name string) (os.FileInfo, error)
	AssetNames func() []string
	Gzip       bool
}
