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
	"log"
	"os"
	"path/filepath"
	"testing"
)

// test is not comprehensive yet, manually observe output
// go clean -testcache && go test . -v

func TestSelf(t *testing.T) {
	t.Logf("starting test...")

	p, e := GetPackage("github.com/clusterkit/go-pkgtool")
	if e != nil {
		t.Fatalf("%s", e)
	}

	for n, f := range p.Funs {
		t.Logf("Indexed Filename: %s", f)
		fn := p.GetFun(n)
		t.Logf("Name: %s", fn.Name)
		t.Logf("Exported: %#v", fn.Exported)
		t.Logf("Filename: %s", fn.Filename)
		t.Logf("Package: %s", fn.Package)
		t.Logf("Params: %s", fn.Params)
		t.Logf("Results: %s", fn.Results)
	}

	e = os.MkdirAll(filepath.Join(".", "test-output"), os.ModePerm)
	if e != nil {
		log.Fatalf("%s", e)
	}

	p.Rename("pkgtool", "somethingelse")
	e = p.Copy("getpackage.go", "./test-output/mynewfile.txt")
	if e != nil {
		log.Fatalf("%s", e)
	}

	e = p.Generate("./test-output/testgen.txt", nil, `package {{ .package.Name }}{{ .NL }}`)
	if e != nil {
		t.Fatalf("%s", e)
	}
}
