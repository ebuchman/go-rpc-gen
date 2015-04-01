package main

import (
	"fmt"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
)

// asserts the map has only one item and returns it
func onePkg(pkgs map[string]*ast.Package) *ast.Package {
	if len(pkgs) > 1 {
		panic("More than one package found. Exiting")
	}
	pkg := new(ast.Package)
	var pkgname string
	for n, p := range pkgs {
		pkg = p
		pkgname = n
	}
	fmt.Println("pkgname:", pkgname)
	_ = pkgname
	return pkg
}

// get the $GOPATH relative path from the dir
func goImportPathFromDir(dir string) (string, error) {
	dir, _ = filepath.Abs(dir)
	prefix := path.Join(GoPath, "src")
	if strings.HasPrefix(dir, prefix) {
		return dir[len(prefix)+1:], nil
	} else {
		return "", fmt.Errorf("%s not on the $GOPATH", dir)
	}
}

// join argument names and types (for function def)
func joinArgTypes(names, types []string) string {
	union := make([]string, len(names))
	for i, n := range names {
		union[i] = fmt.Sprintf("%s %s", n, types[i])
	}
	return strings.Join(union, ", ")
}

// convert camel case to lowercase (CamelCase => camel_case)
func CamelToLower(s string) string {
	lower := ""
	for i := 0; i < len(s); i++ {
		if unicode.IsUpper(rune(s[i])) && i != 0 {
			lower += "_"
		}
		lower += string(unicode.ToLower(rune(s[i])))
	}
	return lower
}

// convert ast.Expr to its go type as it would be in source code
func typeToString(typ ast.Expr) string {
	switch t := typ.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.String()
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	}
	panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(typ)))
}

// return a list of all comments
func getComments(pkg *ast.Package) []*ast.Comment {
	// is this a comment or what
	comments := []*ast.Comment{}
	fs := pkg.Files
	for _, f := range fs {
		for _, c := range f.Comments {
			for _, cc := range c.List {
				comments = append(comments, cc)
			}
		}
	}
	return comments
}

// get list of imports
func getImports(pkg *ast.Package, pkgPath string) map[string]string {
	allImps := make(map[string]string)
	fs := pkg.Files
	for _, f := range fs {
		for _, imp := range f.Imports {
			impPath := imp.Path
			impPathVal := strings.Trim(impPath.Value, "\"")
			name := path.Base(impPathVal)
			impName := imp.Name
			if impName != nil {
				name = impName.Name
			}
			allImps[name] = impPathVal
			fmt.Println("set imp", name, impPathVal)
		}
	}
	return allImps
}

// check if the type is builtin or defined
func isBuiltin(arg string) bool {
	arg, _ = stripPointerArray(arg)
	switch arg {
	case "bool",
		"byte",
		"complex128",
		"complex64",
		"error",
		"float32",
		"float64",
		"int",
		"int16",
		"int32",
		"int64",
		"int8",
		"rune",
		"string",
		"uint",
		"uint16",
		"uint32",
		"uint64",
		"uint8",
		"uintptr":
		return true
	default:
		return false
	}
}

// strip "*" and/or "[]" from the front of a string
func stripPointerArray(imp string) (string, string) {
	pre := ""
	for i := 0; i < len(imp); i++ {
		if strings.Contains("[]*", imp[i:i+1]) {
			pre += imp[i : i+1]
		}
	}
	return imp[len(pre):], pre
}

// returns a list of all exported functions in a pkg
func getFuncs(pkg *ast.Package) map[string]*ast.Object {
	objs := make(map[string]*ast.Object)
	fs := pkg.Files
	for _, f := range fs {
		// Print the scope
		for n, o := range f.Scope.Objects {
			if o.Kind == ast.Fun && ast.IsExported(n) {
				objs[n] = o
			}
		}

	}
	return objs
}

// returns a filter function for parsing a directory
func returnFilter(excludes []string) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		name := info.Name()
		var excluded bool
		for _, ex := range excludes {
			if name == ex {
				excluded = true
				break
			}
		}
		return !info.IsDir() && !excluded && path.Ext(name) == ".go" && !strings.HasSuffix(name, "_test.go")
	}
}
