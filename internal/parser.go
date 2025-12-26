package internal

import (
	"os"
	"path/filepath"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
)

type ParseResult struct {
	Enums   []EnumInfo
	Structs []StructInfo
}

type Parser struct {
	parseResult *ParseResult
}

func NewParser() *Parser {
	return &Parser{
		parseResult: &ParseResult{},
	}
}

// goes through all the go files in the directory and parses them
func (p *Parser) FromDir(dir string) error {
	fset := token.NewFileSet()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fileResult, err := p.parseFile(fset, path)
		if err != nil {
			return err
		}
		p.parseResult.Enums = append(p.parseResult.Enums, fileResult.Enums...)
		p.parseResult.Structs = append(p.parseResult.Structs, fileResult.Structs...)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// parses a single go file
func (p *Parser) parseFile(fset *token.FileSet, path string) (*ParseResult, error) {
	result := &ParseResult{}

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					result.Structs = append(result.Structs, *p.parseStruct(typeSpec.Name.Name, file.Name.Name, structType))
				}
			}
		}
	}

	return result, nil
}

// parses a single struct
func (p *Parser) parseStruct(name, pkgName string, structType *ast.StructType) *StructInfo {
	result := &StructInfo{
		Name:    name,
		Package: pkgName,
		Fields:  make([]FieldInfo, 0),
	}
	for _, field := range structType.Fields.List {

		// Embedded field
		if len(field.Names) == 0 {
			continue
		}

		fieldInfo := FieldInfo{
			Name: field.Names[0].Name,
			Type: p.typeToString(field.Type),
		}

		if strings.HasPrefix(fieldInfo.Type, "*") {
			fieldInfo.IsPointer = true
			fieldInfo.IsOptional = true
		}

		baseType := strings.TrimPrefix(fieldInfo.Type, "*")
		baseType = strings.TrimPrefix(baseType, "[]")

		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			fieldInfo.JSONTag = tag
		}

		result.Fields = append(result.Fields, fieldInfo)
	}
	return result
}

// converts a go type to a string
func (p *Parser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + p.typeToString(t.Elt)
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.SelectorExpr:
		return p.typeToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + p.typeToString(t.Key) + "]" + p.typeToString(t.Value)
	default:
		return "any"
	}
}
