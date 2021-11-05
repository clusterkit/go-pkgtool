/*
   Copyright The ClusterKit Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package pkgtool

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

// gather package information
func GetPackage(dir string) (*Package, error) {
	pkg := (&Package{}).Init()

	// Disable CGO for parsing? Not sure.
	// build.Default.CgoEnabled = false

	// we want to work out of the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get current directory: %v", err)
	}

	// get package details
	pkg.Info, err = build.Default.Import(filepath.ToSlash(dir), cwd, build.ImportComment)
	if err != nil {
		if _, ok := err.(*build.NoGoError); !ok {
			return nil, fmt.Errorf("unable to import %q: %v", dir, err)
		}
	}

	// perhaps the package exists but we requested an empty directory
	if pkg.Info == nil {
		pkg.Info, err = build.Default.Import(filepath.ToSlash(dir), cwd, build.FindOnly)
		if err != nil {
			return nil, err
		}
	}

	// build list of files to parse
	fileList := pkg.Info.GoFiles
	fileList = append(fileList, pkg.Info.TestGoFiles...)

	// parse files into the ast
	for _, gofile := range fileList {
		err := pkg.ParseFile(gofile)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}
