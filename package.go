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
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// something to hold the source data in
type Package struct {
	Info         *build.Package
	TokenFileSet *token.FileSet
	Ast          map[string]*ast.File

	// objects get indexed by map[object_name]file_name
	Bads map[string]string // for error handling
	Pkgs map[string]string // package
	Cons map[string]string // constant
	Typs map[string]string // type
	Vars map[string]string // variable
	Funs map[string]string // function or method
	Lbls map[string]string // label

	// to prevent double inits
	init bool
}

// sets initial values
func (p *Package) Init() *Package {
	if p.init {
		return p
	}

	p.init = true
	p.TokenFileSet = &token.FileSet{}
	p.Ast = make(map[string]*ast.File)
	p.Bads = make(map[string]string)
	p.Pkgs = make(map[string]string)
	p.Cons = make(map[string]string)
	p.Typs = make(map[string]string)
	p.Vars = make(map[string]string)
	p.Funs = make(map[string]string)
	p.Lbls = make(map[string]string)

	return p
}

//
func (p *Package) ResolvePath(path string) (string, string) {
	path = filepath.ToSlash(path)
	var filename string

	// convert import path to go.mod path
	if strings.HasPrefix(path, p.Info.ImportPath) {
		filename = strings.TrimPrefix(path, p.Info.ImportPath)
		path = p.Info.Dir + filename
		return filename, path
	}

	// check this file exists in package dir
	if !strings.Contains(path, "/") {
		filename := path // filename and path are scoped down
		path := filepath.Join(p.Info.Dir, path)
		if _, err := os.Stat(path); err == nil {
			return filename, path
		}
	}

	path = filepath.Join(p.Info.Dir, path)
	pathlist := strings.Split(path, "/")
	filename = pathlist[len(pathlist)-1]
	return filename, path
}

// parses file and adds it info to the package
func (p *Package) ParseFile(path string) error {
	filename, path := p.ResolvePath(path)

	// prevent someone from reusing filenames
	if _, exists := p.Ast[filename]; exists {
		return fmt.Errorf("filename already parsed: %s", filename)
	}

	// get an ast of the file
	parsed, err := parser.ParseFile(p.TokenFileSet, path, nil, 0)
	if err != nil {
		return err
	}

	// build object indexes
	for name := range parsed.Scope.Objects {
		switch parsed.Scope.Objects[name].Kind {
		case ast.Pkg: // package
			p.Pkgs[name] = filename
		case ast.Con: // constant
			p.Cons[name] = filename
		case ast.Typ: // type
			p.Typs[name] = filename
		case ast.Var: // variable
			p.Vars[name] = filename
		case ast.Fun: // function or method
			p.Funs[name] = filename
		case ast.Lbl: // label
			p.Lbls[name] = filename
		default: // for error handling | ast.Bad
			p.Bads[name] = filename
		}
	}

	// add file's ast
	p.Ast[filename] = parsed

	return nil
}

// changes the package name across the package
func (p *Package) Rename(old_name string, new_name string) {
	// rename the package itself for whoever is consuming
	if p.Info.Name == old_name {
		p.Info.Name = new_name
	}

	// rename package names in the ast
	for _, astFile := range p.Ast {
		if astFile.Name.Name == old_name {
			astFile.Name.Name = new_name
		}
	}
}

// copy source to path
func (p *Package) Copy(src string, dest string) error {
	if p.Ast[src] == nil {
		return fmt.Errorf("source not found: %s", src)
	}

	// for debugging: file := new(strings.Builder)

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	printer.Fprint(file, p.TokenFileSet, p.Ast[src])

	return nil
}

// find the package name for a file
func (p *Package) GetPackageName(filename string) string {
	astfile := p.Ast[filename]
	if astfile == nil {
		return ""
	}
	return astfile.Name.Name
}
