package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

var typeProcessors = []typeProcessor{
	processConstructor,
	processBuilder,
	processGetter,
	processSetter,
	processDestructor,
	processWither,
	processStringer,
}

func fileFilter(dir string, fi fs.FileInfo) bool {
	file, err := os.Open(filepath.Join(dir, fi.Name()))
	if err != nil {
		errorLogger.Printf("Failed to open file %s", fi.Name())
		return false
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return true
			}
			errorLogger.Printf("Failed to read file %s", file.Name())
			return false
		}
		if strings.Index(line, ignore) == 0 {
			debugLogger.Printf("Ignoring file %s", file.Name())
			return false
		}
	}
}

func processDirRecursive(dir string) error {
	err := processDir(dir)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.Index(entry.Name(), ".") != 0 {
			err := processDirRecursive(filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processDir(dir string) error {
	debugLogger.Printf("Start parsing directory %s", dir)
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, dir, func(fi fs.FileInfo) bool {
		return fileFilter(dir, fi)
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pkg := range packages {
		debugLogger.Printf("Start parsing package %s", pkg.Name)
		for fileName, file := range pkg.Files {
			writer, fileImports, err := processFile(fset, pkg.Name, fileName, file)
			if err != nil {
				return err
			}
			if writer == nil || fileImports == nil {
				debugLogger.Printf("Ignoring file %s", fileName)
				continue
			}
			err = writer.Write(fileImports)
			if err != nil {
				errorLogger.Printf("Failed to write %s", getGeneratedFileName(fileName[:len(fileName)-3]))
				return err
			}
			infoLogger.Printf("%s is created", getGeneratedFileName(fileName[:len(fileName)-3]))
			return nil
		}
	}
	return nil
}

func processFile(fset *token.FileSet, pkgName string, fileName string, file *ast.File) (*fileWriter, map[string]*impData, error) {
	fileComments := getCommentLines(file.Doc)
	_, found := hasComment(fileComments, "ignore")
	if found {
		return nil, nil, nil
	}
	debugLogger.Printf("Start parsing file %s", fileName)
	fileImports := map[string]*impData{}
	for _, imp := range file.Imports {
		v := imp.Path.Value[1 : len(imp.Path.Value)-1]
		if imp.Name == nil {
			name := path.Base(v)
			fileImports[name] = &impData{
				HasName: false,
				Name:    name,
				Path:    v,
			}
		} else {
			name := imp.Name.Name
			fileImports[name] = &impData{
				HasName: true,
				Name:    name,
				Path:    v,
			}
		}
	}
	writer := NewFileWriter(pkgName, fileName[:len(fileName)-3])
	imports := map[string]bool{}
	for _, decl := range file.Decls {
		if decl == nil {
			continue
		}
		g, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		_, err := processNode(g, imports, writer)
		if err != nil {
			return writer, fileImports, err
		}
	}
/*
	var processErr error = nil
	ast.Inspect(file, func(n ast.Node) bool {
		if processErr != nil {
			return false
		}
		if n == nil {
			return true
		}
		result, err := processNode(n, imports, writer)
		if processErr != nil {
			processErr = err
		}
		return result
	})
	if processErr != nil {
		return writer, fileImports, processErr
	}*/
	for name, imp := range imports {
		writer.imports[name] = imp
	}
	return writer, fileImports, nil
}

func processNode(n ast.Node, imports map[string]bool, writer *fileWriter) (bool, error) {
	g, ok := n.(*ast.GenDecl)
	if !ok {
		return true, nil
	}
	for _, spec := range g.Specs {
		t, ok := spec.(*ast.TypeSpec)
		if !ok {
			return false, nil
		}
		s, ok := t.Type.(*ast.StructType)
		if !ok {
			return false, nil
		}
		genericTypes := map[string]string{}
		genericTypeNames := []string{}
		if t.TypeParams != nil {
			for _, typeParam := range t.TypeParams.List {
				fieldName := typeParam.Names[0].String()
				genericType, importNames := readRootExpression(typeParam.Type, imports)
				for imp := range importNames {
					imports[imp] = true
				}
				genericTypes[fieldName] = genericType
				genericTypeNames = append(genericTypeNames, fieldName)
			}
		}
		structName := t.Name.Name
		fields := map[string]string{}
		fieldComments := map[string][]string{}
		fieldNames := []string{}
		for _, field := range s.Fields.List {
			fieldName := field.Names[0].String()
			typeName, importNames := readRootExpression(field.Type, imports)
			for imp := range importNames {
				imports[imp] = true
			}
			fields[fieldName] = typeName
			fieldComments[fieldName] = append(getCommentLines(field.Doc), readTags(field, writer.pkg, structName, fieldName)...)
			fieldNames = append(fieldNames, fieldName)
		}
		comments := getCommentLines(g.Doc)
		data := &typeProcessorData{
			packageName:      writer.pkg,
			structName:       structName,
			fields:           fields,
			fieldComments:    fieldComments,
			fieldNames:       fieldNames,
			typeComments:     comments,
			genericTypes:     genericTypes,
			genericTypeNames: genericTypeNames,
			addImport: func(i string) {
				writer.imports[i] = true
			},
			addCodeWriter: func(cw codeWriter) {
				writer.codeWriters = append(writer.codeWriters, cw)
			},
		}
		for _, typeProcessor := range typeProcessors {
			err := typeProcessor(data)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func getCommentLines(comments *ast.CommentGroup) []string {
	c := []string{}
	if comments == nil || comments.List == nil {
		return c
	}
	prefix := "//gombok:"
	for _, comment := range comments.List {
		if strings.Index(comment.Text, prefix) == 0 {
			c = append(c, comment.Text[len(prefix):])
		}
	}
	return c
}

func readTags(field *ast.Field, packageName, structName, fieldName string) []string {
	tags := []string{}
	if field.Tag == nil {
		return tags
	}
	tagValue := field.Tag.Value
	tagValue = tagValue[1 : len(tagValue)-1]
	for {
		index := strings.Index(tagValue, ":")
		if index < 0 {
			if tagValue != "" {
				errorLogger.Printf("Invalid tag on %s.%s.%s", packageName, structName, fieldName)
			}
			return tags
		}
		key := tagValue[:index]
		if tagValue[index+1] != '"' {
			errorLogger.Printf("Invalid tag on %s.%s.%s", packageName, structName, fieldName)
			return tags
		}
		endIndex := strings.Index(tagValue[index+2:], "\"") + index + 2
		if endIndex == -1 {
			errorLogger.Printf("Invalid tag on %s.%s.%s", packageName, structName, fieldName)
			return tags
		}
		value := tagValue[index+2 : endIndex]
		if strings.Index(key, "gombok") == 0 {
			prop := key[len("gombok"):]
			tags = append(tags, fmt.Sprintf("%s %s", prop, value))
		}
		tagValue = strings.Trim(tagValue[endIndex+1:], " \t")
	}
}

func hasComment(comments []string, commentText string) ([]string, bool) {
	var commands []string = nil
	found := false
	for _, comment := range comments {
		if strings.Index(comment, commentText) == 0 {
			if found {
				errorLogger.Printf("Duplicated comment %s", commentText)
				continue
			}
			rest := comment[len(commentText):]
			if rest == "" {
				commands, found = []string{}, true
			} else {
				commands, found = strings.Split(rest[1:], " "), true
			}
		}
	}
	return commands, found
}

func readRootExpression(expr ast.Expr, imports map[string]bool) (string, map[string]bool) {
	return readExpression(expr, nil, "", imports)
}

func readChildExpression(expr ast.Expr, parent ast.Expr, imports map[string]bool) (string, map[string]bool) {
	return readExpression(expr, parent, "", imports)
}

func readExpression(expr ast.Expr, parent ast.Expr, name string, imports map[string]bool) (string, map[string]bool) {
	switch e := expr.(type) {
	case *ast.Ident:
		return fmt.Sprintf("%s%s", name, e.Name), imports
	case *ast.SelectorExpr:
		i, _ := readChildExpression(e.X, e, imports)
		_, found := imports[i]
		if !found {
			imports[i] = true
		}
		n, _ := readChildExpression(e.Sel, e, imports)
		return fmt.Sprintf("%s%s.%s", name, i, n), imports
	case *ast.StarExpr:
		n, _ := readChildExpression(e.X, e, imports)
		return fmt.Sprintf("%s*%s", name, n), imports
	case *ast.ArrayType:
		n, _ := readChildExpression(e.Elt, e, imports)
		return fmt.Sprintf("%s[]%s", name, n), imports
	case *ast.MapType:
		k, _ := readChildExpression(e.Key, e, imports)
		v, _ := readChildExpression(e.Value, e, imports)
		return fmt.Sprintf("%smap[%s]%s", name, k, v), imports
	case *ast.UnaryExpr:
		x, _ := readChildExpression(e.X, e, imports)
		return fmt.Sprintf("%s%s%s", name, e.Op.String(), x), imports
	case *ast.BinaryExpr:
		x, _ := readChildExpression(e.X, e, imports)
		y, _ := readChildExpression(e.Y, e, imports)
		return fmt.Sprintf("%s%s %s %s", name, x, e.Op.String(), y), imports
	case *ast.ChanType:
		var ch string
		if e.Dir == ast.SEND {
			ch = "chan<-"
		} else if e.Dir == ast.RECV {
			ch = "<-chan"
		} else {
			ch = "chan"
		}
		v, _ := readChildExpression(e.Value, e, imports)
		return fmt.Sprintf("%s%s %s", name, ch, v), imports
	case *ast.InterfaceType:
		methods := []string{}
		for _, f := range e.Methods.List {
			v, _ := readField(f, e, imports)
			methods = append(methods, v)
		}
		if len(methods) == 0 {
			return fmt.Sprintf("%sinterface{}", name), imports
		}
		if len(methods) == 1 {
			return fmt.Sprintf("%sinterface{ %s }", name, methods[0]), imports
		}
		// return fmt.Sprintf("%sinterface{ %s }", name, strings.Join(methods, "; ")), imports
		return fmt.Sprintf("%sinterface {\n\t%s\n}", name, strings.Join(methods, "\n\t")), imports
	case *ast.StructType:
		fields := []string{}
		for _, f := range e.Fields.List {
			v, _ := readField(f, e, imports)
			fields = append(fields, v)
		}
		if len(fields) == 0 {
			return fmt.Sprintf("%sstruct{}", name), imports
		}
		if len(fields) == 1 {
			return fmt.Sprintf("%sstruct{ %s }", name, fields[0]), imports
		}
		// return fmt.Sprintf("%sinterface{ %s }", name, strings.Join(fields, "; ")), imports
		return fmt.Sprintf("%sstruct {\n\t%s\n}", name, strings.Join(fields, "\n\t")), imports
	case *ast.FuncType:
		params := []string{}
		for _, f := range e.Params.List {
			v, _ := readField(f, e, imports)
			params = append(params, v)
		}
		results := []string{}
		if e.Results != nil {
			for _, f := range e.Results.List {
				v, _ := readField(f, e, imports)
				results = append(results, v)
			}
		}
		genericFieldNames := []string{}
		if e.TypeParams != nil {
			for _, field := range e.TypeParams.List {
				fieldName := field.Names[0].String()
				typeName, importNames := readRootExpression(field.Type, imports)
				for imp := range importNames {
					imports[imp] = true
				}
				genericFieldNames = append(genericFieldNames, fmt.Sprintf("%s %s", fieldName, typeName))
			}
		}
		genericText := ""
		if len(genericFieldNames) != 0 {
			genericText = fmt.Sprintf("[%s]", strings.Join(genericFieldNames, ", "))
		}
		if len(results) == 0 {
			switch parent.(type) {
			case *ast.InterfaceType:
				return fmt.Sprintf("%s%s(%s)", name, genericText, strings.Join(params, ", ")), imports
			default:
				return fmt.Sprintf("%sfunc%s(%s)", name, genericText, strings.Join(params, ", ")), imports
			}
		}
		joinedResults := strings.Join(results, ", ")
		switch parent.(type) {
		case *ast.InterfaceType:
			if strings.Index(joinedResults, " ") == -1 {
				return fmt.Sprintf("%s(%s) %s", name, strings.Join(params, ", "), joinedResults), imports
			} else {
				return fmt.Sprintf("%s(%s) (%s)", name, strings.Join(params, ", "), joinedResults), imports
			}
		default:
			if strings.Index(joinedResults, " ") == -1 {
				return fmt.Sprintf("%sfunc%s(%s) %s", name, genericText, strings.Join(params, ", "), joinedResults), imports
			} else {
				return fmt.Sprintf("%sfunc%s(%s) (%s)", name, genericText, strings.Join(params, ", "), joinedResults), imports
			}
		}
	default:
		errorLogger.Printf("fix type %s", reflect.TypeOf(expr).String())
		panic(fmt.Errorf("Error"))
	}
}

func readField(field *ast.Field, parent ast.Expr, imports map[string]bool) (string, map[string]bool) {
	v, _ := readChildExpression(field.Type, parent, imports)
	names := field.Names
	if len(names) == 0 {
		return v, imports
	} else {
		n := []string{}
		for _, name := range names {
			n = append(n, name.Name)
		}
		switch parent.(type) {
		case *ast.InterfaceType:
			return fmt.Sprintf("%s%s", strings.Join(n, ", "), v), imports
		default:
			return fmt.Sprintf("%s %s", strings.Join(n, ", "), v), imports
		}
	}
}
