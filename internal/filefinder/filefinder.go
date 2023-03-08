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
			return nil, fmt.Errorf("unable to find state in pattern '%s': %+v", p, err)
		}
	} else {
		fileInfo, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("'%s' path does not exist: %+v", p, err)
		}

		if fileInfo.IsDir() {
			files, err = glob(path.Join(p, endPattern))
			if err != nil {
				return nil, fmt.Errorf("unable to find state in path '%s': %+v", p, err)
			}
		}
	}

	return files, nil
}
