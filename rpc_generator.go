package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	gofmt "go/format"
	goparser "go/parser"
	gotoken "go/token"
	//"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	GoPath = os.Getenv("GOPATH")

	interfaceF = flag.String("interface", "", "interface type to define the rpc methods on")
	typeF      = flag.String("type", "", "comma separated list of types that should implement the interface")
	pkgF       = flag.String("pkg", "", "package containing functions providing the core functionality for the rpc")
	outF       = flag.String("out", "", "output package for client methods")
	outPkgF    = flag.String("out-pkg", "", "name of the package for which code is to be generated")
	excludeF   = flag.String("exclude", "", "comma separated list of files to exclude public functions from (relative to pkg)")
	templatesF = flag.String("templates", ".", "file/s in which the template functions are located")
)

func main() {

	flag.Parse()

	iface := *interfaceF
	types := strings.Split(*typeF, ",")
	dir := *pkgF
	outFile := *outF
	excludes := strings.Split(*excludeF, ",")
	outPkg := *outPkgF
	//	templateFiles := strings.Split(*templatesF, ",")

	_ = outFile

	fset := gotoken.NewFileSet() // positions are relative to fset

	// get the core functions to be exposed
	corePkgs, err := goparser.ParseDir(fset, dir, returnFilter(excludes), 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(corePkgs)
	corePkg := onePkg(corePkgs)
	coreFuncs := getFuncs(corePkg)
	corePkgImportPath, err := goImportPathFromDir(dir)
	if err != nil {
		panic(err)
	}

	// get the interface to be populated (present in current dir)
	pkgs, err := goparser.ParseDir(fset, ".", nil, goparser.ParseComments)
	if err != nil {
		panic(err)
	}
	pkg := onePkg(pkgs)

	interfaceDef := fmt.Sprintf(`
type %s interface{

}`, iface)

	// init the rpc generator by parsing the templates and definitions
	rpcGen, err := initRpcGen(pkg)
	if err != nil {
		panic(err)
	}
	if len(rpcGen.templates) != len(types) {
		panic(fmt.Sprintf("rpc-gen requires equal numbers of types and templates. Got %d, %d", len(types), len(rpcGen.templates)))
	}

	imports := getImports(corePkg, corePkgImportPath)
	// populate interface and stringify func defs
	stringFuncs, interfaceDef, neededImports := populateInterface(interfaceDef, coreFuncs, imports, dir, corePkgImportPath)

	// add base imports to neededImports
	for k, v := range rpcGen.imports {
		neededImports[k] = v
	}
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "// File generated by github.com/ebuchman/rpc-gen")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "package", outPkg)
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "import(")
	for n, im := range neededImports {
		fmt.Println(n, im)
		if n != path.Base(im) {
			fmt.Fprintln(buf, "\t"+n+" \""+im+"\"")
		} else {
			fmt.Fprintln(buf, "\t\""+im+"\"")

		}
	}
	fmt.Fprintln(buf, ")")
	fmt.Fprintln(buf, "")
	fmt.Fprintf(buf, interfaceDef)

	fmt.Println(string(buf.Bytes()))

	// for each client type, implement the interface
	// using its template and the stringFuncs
	for _, clientType := range types {
		implementation, err := rpcGen.implementInterface(clientType, stringFuncs)
		if err != nil {
			panic(err)
		}
		// write implementation to buffer
		buf.Write(implementation)
	}

	//fmt.Println(string(buf.Bytes()))
	// save buffer to file
	newFile := "client_methods.go"
	data := buf.Bytes()

	node, err := goparser.ParseFile(fset, "", data, 0)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(newFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gofmt.Node(f, fset, node)
	/*if err := ioutil.WriteFile(newFile, data, 0660); err != nil {
		panic(err)
	}*/
}

func defToCall(def string, args []string) string {
	return ""
}

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
				fmt.Fprintf(buf, "nil")
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

// sets the context for implementing a template
func (rg *RpcGen) SetContext(txt []string, jobs []Job) {
	rg.txt = txt
	rg.jobs = jobs
}

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

// clean string based representation of a go function
type Func struct {
	Name        string
	ArgNames    []string
	ArgTypes    []string
	ReturnTypes []string
}

func NewFunc(name string, nargs, nret int) Func {
	return Func{
		Name:        name,
		ArgNames:    make([]string, nargs),
		ArgTypes:    make([]string, nargs),
		ReturnTypes: make([]string, nret), //[]string{"*Response" + name},
	}
}

// convert a function object to a clean string representation of args and returns
func objectToStringFunc(name string, obj *ast.Object) Func {
	fdecl := obj.Decl.(*ast.FuncDecl)
	ftype := fdecl.Type
	argList := ftype.Params.List
	retList := ftype.Results.List
	thisFunc := NewFunc(name, len(argList), len(retList))
	for i, p := range argList {
		t := typeToString(p.Type)
		n := p.Names[0].Name
		thisFunc.ArgNames[i] = n
		thisFunc.ArgTypes[i] = t
	}
	for i, r := range retList {
		thisFunc.ReturnTypes[i] = typeToString(r.Type)
	}
	return thisFunc
}

// create an interface definition containing all defined methods
func populateInterface(baseDef string, funcs map[string]*ast.Object, imports map[string]string, pkgName, pkgPath string) ([]*Func, string, map[string]string) {
	stringFuncs := make([]*Func, len(funcs))

	neededImps := make(map[string]string)

	// pull off the final }
	baseDef = baseDef[:len(baseDef)-1]
	// extract each functions string info,
	// add to stringFuncs and append to interface def
	i := 0 // using append on stringFuncs was breaking ...
	for name, obj := range funcs {
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

type RpcGen struct {
	templates map[string]string
	ifaceDef  string
	funcdefs  map[string]string

	imports map[string]string // default imports for template functions

	txt  []string
	jobs []Job
}

func initRpcGen(pkg *ast.Package) (*RpcGen, error) {
	rpcGen := &RpcGen{
		templates: make(map[string]string),
		funcdefs:  make(map[string]string),
		imports:   make(map[string]string),
	}

	comments := getComments(pkg)
	//funcs := getFuncs(pkgs)
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
			rpcGen.ifaceDef = defn // */ ?

		case "define-func":
			//name := defs[1]
			// the function definition follows the comment
			// in the code.
			//pos := c.End()

			//rpcGen.funcdefs[name] = defn
		case "imports":
			rest = txt[len("imports"):]
			fmt.Println("REST:", rest)
			spl := strings.Split(rest, "\n")
			fmt.Println("SPL:", spl)
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
				fmt.Println(importName, importPath)
				rpcGen.imports[importName] = importPath
			}
		default:
			// expects a name registered by "define-set"
			// ie. serialization routine for a particular type

		}
	}
	return rpcGen, nil
}

// update a function's arg/return types by appending the package name if necessary
// and add packages to neededImps
func updateFunctionAndImport(f *Func, allImps map[string]string, neededImps map[string]string, pkgName, pkgPath string) {
	fmt.Println(f)
	for i, arg := range f.ArgTypes {
		fmt.Println(i, arg)
		f.ArgTypes[i] = updateImport(arg, allImps, &neededImps, pkgName, pkgPath)
	}

	for i, ret := range f.ReturnTypes {
		fmt.Println(i, ret)
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
