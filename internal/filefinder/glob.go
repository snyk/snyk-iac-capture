package filefinder

import (
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func hasMeta(path string) bool {
	magicChars := `?*[]`
	return strings.ContainsAny(path, magicChars)
}

func glob(pathPattern string) ([]string, error) {
	var files []string

	base, pattern := doublestar.SplitPattern(pathPattern)
	err := doublestar.GlobWalk(os.DirFS(base), pattern, func(p string, d fs.DirEntry) error {
		// Ensure paths aren't actually directories
		// For example when the directory matches the filefinder pattern like it's a file
		if !d.IsDir() {
			files = append(files, path.Join(base, p))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
