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
	"go/ast"
	"log"
	"os"
)

func (p *Package) GetFun(name string) *FunInfo {
	// need to get the filename
	filename := p.Funs[name]

	// need to get the ast by filename
	astfile := p.Ast[filename]
	if astfile == nil {
		return nil
	}

	// get the object from astfile
	fnObj := astfile.Scope.Lookup(name)
	if fnObj == nil {
		return nil
	}

	// begin building function info
	infoOut := &FunInfo{
		Filename: filename,
		Exported: ast.IsExported(name),
		Package:  p.GetPackageName(filename),
		Name:     name,
		Params:   make([]byte, 0),
		Results:  make([]byte, 0),
	}

	// assert the declaration to function
	fnDecl := fnObj.Decl.(*ast.FuncDecl)

	// extract the full file path
	srcFilePath := p.TokenFileSet.File(astfile.Pos()).Name()
	infoOut.FullPath = srcFilePath

	// open the source file to read contents
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer srcFile.Close()

	// get the params definition
	if fnDecl.Type.Params != nil {
		// allocate bytes
		infoOut.Params = make([]byte, fnDecl.Type.Params.End()-fnDecl.Type.Params.Pos())
		// the astfile data is relative to it's position in token fileset
		// remove the astfile position from the position of item desired
		startPos := int64(fnDecl.Type.Params.Pos() - astfile.Pos())
		// grab as defined in source
		srcFile.ReadAt(infoOut.Params, startPos)
	}

	// get the result definition
	if fnDecl.Type.Results != nil {
		// allocate bytes
		infoOut.Results = make([]byte, fnDecl.Type.Results.End()-fnDecl.Type.Results.Pos())
		// the astfile data is relative to it's position in token fileset
		// remove the astfile position from the position of item desired
		startPos := int64(fnDecl.Type.Results.Pos() - astfile.Pos())
		// grab as defined in source
		srcFile.ReadAt(infoOut.Results, startPos)
	}

	return infoOut
}
