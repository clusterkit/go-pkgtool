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
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func (p *Package) Generate(dest string, data interface{}, tplSource string) error {
	// new template renderer
	renderer := template.New(dest) //.Option("missingkey=error")
	renderer.Funcs(sprig.TxtFuncMap())

	// parse the template
	renderer, err := renderer.Parse(tplSource)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	datamap := make(map[string]interface{})
	datamap["package"] = p.Info
	datamap["data"] = data
	datamap["bad_objects"] = p.Bads // for error handling
	datamap["packages"] = p.Pkgs    // package
	datamap["constants"] = p.Cons   // constant
	datamap["types"] = p.Typs       // type
	datamap["variables"] = p.Vars   // variable
	datamap["functions"] = p.Funs   // function or method
	datamap["labels"] = p.Lbls      // label
	datamap["NL"] = "\n"

	// render the template to output
	err = renderer.Execute(file, datamap)
	if err != nil {
		return err
	}

	return nil
}
