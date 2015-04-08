package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"path"
	"sort"
	"strconv"
	"strings"
)

//--------------------------------------------------------------------------------
// generate source code from template

// interpret/replace simple commands found in templates
func (rg *RpcGen) compileJob(buf *bytes.Buffer, f Func, job Job) error {
	argNames := f.ArgNames
	argTypes := f.ArgTypes
	retTypes := f.ReturnTypes
	// ident is either a keyword or a variable name
	spl := strings.Split(job.ident, ".")
	ident := spl[0]
	switch ident {
	case "name":
		fmt.Fprintf(buf, f.Name)
	case "args":
		field := spl[1]
		if i, err := strconv.Atoi(field); err == nil {
			argNames = []string{f.ArgNames[i]}
			argTypes = []string{f.ArgTypes[i]}
		}
		switch spl[1] {
		case "def":
			fmt.Fprintf(buf, joinArgTypes(argNames, argTypes))
		case "ident":
			if len(f.ArgNames) == 0 {
				fmt.Fprintf(buf, "")
			} else {
				fmt.Fprintf(buf, strings.Join(argNames, ", "))
			}
		case "name":
			if len(f.ArgNames) == 0 {
				fmt.Fprintf(buf, "nil")
			} else {
				fmt.Fprintf(buf, "[]string{\""+strings.Join(argNames, "\" , \"")+"\"}")
			}
		}
	case "response":
		if len(spl) > 1 {
			field := spl[1]
			if i, err := strconv.Atoi(field); err == nil {
				retTypes = []string{f.ReturnTypes[i]}
			}
		}
		fmt.Fprintf(buf, strings.Join(retTypes, ", "))
	case "lowername":
		fmt.Fprintf(buf, "\""+CamelToLower(f.Name)+"\"")
	default:
		// check if the ident is registered
		// and if so call the function
		if def, ok := rg.funcdefs[job.ident]; ok {
			fmt.Fprintf(buf, defToCall(def, argNames))
		} else {
			return fmt.Errorf("Unknown identifier %s", job.ident)
		}
	}
	return nil
}

// implement a template for a given function
func (rg *RpcGen) makeMethod(buf *bytes.Buffer, f Func) error {
	for i, t := range rg.txt {
		// write the preceding text
		fmt.Fprintf(buf, t)

		// compile a job to txt
		if i < len(rg.jobs) {
			if err := rg.compileJob(buf, f, rg.jobs[i]); err != nil {
				return err
			}
		}
	}
	fmt.Fprintf(buf, "\n\n")
	return nil
}

//--------------------------------------------------------------------------------
// manage the client interface

// parse the template. for each function, implement the template
func (rg *RpcGen) implementInterface(clientType string, stringFuncs []*Func) ([]byte, error) {
	tmp := rg.templates[clientType]
	p := Parser(tmp)
	if err := p.run(); err != nil {
		return nil, err
	}
	txt, jobs := p.results()
	//fmt.Println(jobs)
	rg.SetContext(txt, jobs)
	buf := new(bytes.Buffer)
	for _, f := range stringFuncs {
		rg.makeMethod(buf, *f)
	}
	//fmt.Println(string(buf.Bytes()))
	return buf.Bytes(), nil
}

// create an interface definition containing all defined methods
func populateInterface(baseDef string, funcs map[string]*ast.Object, imports map[string]string, pkgName, pkgPath string) ([]*Func, string, map[string]string) {
	stringFuncs := make([]*Func, len(funcs))

	neededImps := make(map[string]string)

	// sort functions alphabetically
	funcNames := []string{}
	for n, _ := range funcs {
		funcNames = append(funcNames, n)
	}
	sort.Strings(funcNames)

	// pull off the final }
	baseDef = baseDef[:len(baseDef)-1]
	// extract each functions string info,
	// add to stringFuncs and append to interface def
	i := 0 // using append on stringFuncs was breaking ...
	for _, name := range funcNames {
		obj := funcs[name]
		baseDef += "\t" + name + "("
		thisFunc := objectToStringFunc(name, obj)
		updateFunctionAndImport(&thisFunc, imports, neededImps, pkgName, pkgPath)
		for i, n := range thisFunc.ArgNames {
			t := thisFunc.ArgTypes[i]
			baseDef += n + " " + t + ", "
		}
		stringFuncs[i] = &thisFunc
		if len(thisFunc.ArgNames) > 0 { // argList (?)
			baseDef = baseDef[:len(baseDef)-2]
		}
		baseDef += ") ("
		baseDef += strings.Join(thisFunc.ReturnTypes, ", ")
		baseDef += ")\n"
		//baseDef += fmt.Sprintf(" (*%s.Response%s, error)\n", pkg, name)
		i += 1
	}
	return stringFuncs, baseDef + "\n}\n", neededImps
}

//--------------------------------------------------------------------------------
// stringify/parse/manipulate function definitions

// clean string based representation of a go function
type Func struct {
	Name        string
	ArgNames    []string
	ArgTypes    []string
	ReturnTypes []string
}

func NewFunc(name string) Func {
	return Func{
		Name:        name,
		ArgNames:    []string{},
		ArgTypes:    []string{},
		ReturnTypes: []string{},
	}
}

// convert a function object to a clean string representation of args and returns
func objectToStringFunc(name string, obj *ast.Object) Func {
	fdecl := obj.Decl.(*ast.FuncDecl)
	ftype := fdecl.Type
	argList := ftype.Params.List
	retList := ftype.Results.List
	fmt.Println(name, argList, retList)
	thisFunc := NewFunc(name) //, len(argList), len(retList))
	for _, p := range argList {
		t := typeToString(p.Type)
		for _, n := range p.Names {
			thisFunc.ArgNames = append(thisFunc.ArgNames, n.Name)
			thisFunc.ArgTypes = append(thisFunc.ArgTypes, t)
		}
	}
	for _, r := range retList {
		thisFunc.ReturnTypes = append(thisFunc.ReturnTypes, typeToString(r.Type))
	}
	return thisFunc
}

// update a function's arg/return types by appending the package name if necessary
// and add packages to neededImps
func updateFunctionAndImport(f *Func, allImps map[string]string, neededImps map[string]string, pkgName, pkgPath string) {
	for i, arg := range f.ArgTypes {
		f.ArgTypes[i] = updateImport(arg, allImps, &neededImps, pkgName, pkgPath)
	}

	for i, ret := range f.ReturnTypes {
		f.ReturnTypes[i] = updateImport(ret, allImps, &neededImps, pkgName, pkgPath)
	}

}

// updateFunctionAndImport for a single arg/ret
func updateImport(arg string, allImps map[string]string, neededImps *map[string]string, pkgName, pkgPath string) string {
	// if there's a `.`, figure out the
	// needed import
	spl := strings.Split(arg, ".")
	if len(spl) > 1 {
		imp := spl[0] // may have *s and []s
		imp, _ = stripPointerArray(imp)
		(*neededImps)[imp] = allImps[imp]
		return arg
	} else {

		// determine if type not built in
		// (ie. was defined in this package)
		if !isBuiltin(arg) {
			(*neededImps)[pkgName] = pkgPath
			argNoStar, prefix := stripPointerArray(arg)
			newType := prefix + pkgName + "." + argNoStar
			return newType
		}
	}
	return arg
}

//--------------------------------------------------------------------------------
// main RpcGen object

type RpcGen struct {
	templates map[string]string
	ifaceDef  string
	funcdefs  map[string]string

	imports map[string]string // default imports for template functions

	txt  []string
	jobs []Job
}

// sets the context for implementing a template after parsing
func (rg *RpcGen) SetContext(txt []string, jobs []Job) {
	rg.txt = txt
	rg.jobs = jobs
}

// initialize the rpc generator from a pkg by parsing comments
func initRpcGen(pkg *ast.Package) (*RpcGen, error) {
	rpcGen := &RpcGen{
		templates: make(map[string]string),
		funcdefs:  make(map[string]string),
		imports:   make(map[string]string),
	}

	comments := getComments(pkg)
	for _, c := range comments {
		txt := c.Text[2:]
		if !strings.HasPrefix(txt, "rpc-gen:") {
			continue
		}

		txt = txt[len("rpc-gen:"):]
		txtspl := strings.SplitN(txt, " ", 2)
		rest := ""
		if len(txtspl) == 2 {
			rest = txtspl[1]
		}
		def := txtspl[0]

		defs := strings.Split(def, ":")
		typ := defs[0]
		switch typ {
		case "template":
			txt = txt[len("template:"):]
			// next token up to a space should be the client type
			name := ""
			i := 0
			for ; txt[i:i+1] != " "; i++ {
				name += txt[i : i+1]
			}
			txt = txt[i : len(txt)-2]
			//fmt.Println("TEMPLATE:", name, txt)
			rpcGen.templates[name] = txt
		case "define-set":
			// TODO
		case "define-interface":
			// the interface is all in a comment
			spl := strings.SplitN(rest, "\n", 2)
			//name := spl[0]
			defn := spl[1]
			rpcGen.ifaceDef = defn //

		case "define-func":
			//name := defs[1]
			// the function definition follows the comment
			// in the code.
			//pos := c.End()

			//rpcGen.funcdefs[name] = defn
		case "imports":
			rest = txt[len("imports"):]
			spl := strings.Split(rest, "\n")
			for _, s := range spl[1:] {
				if strings.HasPrefix(s, "*") {
					break
				}
				sp := strings.Split(s, " ")
				importName, importPath := "", ""
				if len(sp) > 1 {
					importName = sp[0]
					importPath = sp[1]
				} else {
					importPath = sp[0]
					importName = path.Base(importPath)
				}
				rpcGen.imports[importName] = importPath
			}
		default:
			// expects a name registered by "define-set"
			// ie. serialization routine for a particular type

		}
	}
	return rpcGen, nil
}
