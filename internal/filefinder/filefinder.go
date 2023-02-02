package filefinder

import (
	"fmt"
	"os"
	"path"
)

// FindFiles get as input a path (directory, glob pattern, file) and give back a list of usable files.
// endPattern is the glob pattern to use when given a directory.
func FindFiles(p, endPattern string) ([]string, error) {
	files := []string{
		p,
	}
	var err error

	if hasMeta(p) {
		files, err = glob(p)
		if err != nil {
			return nil, fmt.Errorf("unable to find state in pattern ''%s: %+v", p, err)
		}
	} else {
		fileInfo, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("'%s' p does not exists: %+v", p, err)
		}

		if fileInfo.IsDir() {
			files, err = glob(path.Join(p, endPattern))
			if err != nil {
				return nil, fmt.Errorf("unable to find state in p ''%s: %+v", p, err)
			}
		}
	}

	return files, nil
}
