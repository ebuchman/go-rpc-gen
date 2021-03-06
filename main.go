package main

import (
	"bytes"
	"flag"
	"fmt"
	gofmt "go/format"
	goparser "go/parser"
	gotoken "go/token"
	"os"
	"path"
	"strings"
)

var (
	GoPath = os.Getenv("GOPATH")

	interfaceF = flag.String("interface", "", "interface type to define the rpc methods on")
	typeF      = flag.String("type", "", "comma separated list of types that should implement the interface")
	pkgNameF   = flag.String("pkg", "", "package containing functions providing the core functionality for the rpc")
	dirF       = flag.String("dir", "", "relative directory of package containing functions")
	outF       = flag.String("out", "client_methods.go", "output file for client methods")
	outPkgF    = flag.String("out-pkg", "", "name of the package for which code is to be generated")
	excludeF   = flag.String("exclude", "", "comma separated list of files to exclude public functions from (relative to pkg)")
	//templatesF = flag.String("templates", ".", "file/s in which the template functions are located")
)

func main() {

	flag.Parse()

	iface := *interfaceF
	types := strings.Split(*typeF, ",")
	pkgName := *pkgNameF
	dir := *dirF
	outFile := *outF
	excludes := strings.Split(*excludeF, ",")
	outPkg := *outPkgF
	//	templateFiles := strings.Split(*templatesF, ",")

	fset := gotoken.NewFileSet() // positions are relative to fset

	// get the core functions to be exposed
	corePkgs, err := goparser.ParseDir(fset, dir, returnFilter(excludes), 0)
	if err != nil {
		panic(err)
	}
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
	stringFuncs, interfaceDef, neededImports := populateInterface(interfaceDef, coreFuncs, imports, pkgName, corePkgImportPath)

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

	// parse the generated source text for the sake of gofmt
	data := buf.Bytes()
	node, err := goparser.ParseFile(fset, "", data, goparser.ParseComments)
	if err != nil {
		panic(err)
	}

	// gofmt and write to file
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gofmt.Node(f, fset, node)
}
