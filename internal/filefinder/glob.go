/*
 * © 2023 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
