package internal

import (
	"fmt"
	"os"
	"strings"
)

type Generator struct {
	inputDir   string
	outputFile string
	parser     *Parser
}

func New() *Generator {
	return &Generator{
		parser:     NewParser(),
		inputDir:   "",
		outputFile: "",
	}
}

// sets the input directory
func (g *Generator) FromDir(dir string) *Generator {
	g.inputDir = dir
	return g
}

// sets the output file
func (g *Generator) ToFile(file string) *Generator {
	g.outputFile = file
	return g
}

// generates the TypeScript code
func (g *Generator) Generate() error {
	if g.inputDir == "" {
		return fmt.Errorf("input directory not set")
	}

	if g.outputFile == "" {
		return fmt.Errorf("output file not set")
	}

	err := g.parser.FromDir(g.inputDir)
	if err != nil {
		return fmt.Errorf("failed to parse directory: %w", err)
	}

	ts := g.generateTypeScript()
	err = os.WriteFile(g.outputFile, []byte(ts), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// generates the TypeScript code
func (g *Generator) generateTypeScript() string {
	var sb strings.Builder

	sb.WriteString("/* Do not change, this code is generated from Golang structs */\n\n")

	for _, structInfo := range g.parser.parseResult.Structs {
		sb.WriteString(g.generateStruct(structInfo))
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// generates the TypeScript code for a struct
func (g *Generator) generateStruct(structInfo StructInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("export interface %s {\n", structInfo.Name))
	for _, field := range structInfo.Fields {
		tsType := g.goTypeToTS(field)
		optionalMarker := ""
		if field.IsOptional {
			optionalMarker = "?"
		}

		sb.WriteString(fmt.Sprintf("    %s%s: %s;\n", field.Name, optionalMarker, tsType))
	}

	sb.WriteString("};\n")

	return sb.String()
}

// converts a go type to a TypeScript type
func (g *Generator) goTypeToTS(field FieldInfo) string {
	goType := field.Type

	isArray := strings.HasPrefix(goType, "[]")
	if isArray {
		goType = strings.TrimPrefix(goType, "[]")
	}

	isPointer := strings.HasPrefix(goType, "*")
	if isPointer {
		goType = strings.TrimPrefix(goType, "*")
	}

	if strings.HasPrefix(goType, "map[") {
		return g.mapTypeToTS(goType)
	}

	var tsType string

	isKnownStruct := false
	for _, s := range g.parser.parseResult.Structs {
		if s.Name == goType {
			isKnownStruct = true
			break
		}
	}

	if isKnownStruct {
		tsType = goType
	} else {
		tsType = g.basicTypeToTS(goType)
	}

	if isArray {
		tsType = tsType + "[]"
	}

	if isPointer {
		tsType = tsType + " | null"
	}

	return tsType
}

// converts a basic go type to a TypeScript type
func (g *Generator) basicTypeToTS(goType string) string {
	switch goType {
	case "string":
		return "string"

	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte", "rune":
		return "number"

	case "bool":
		return "boolean"

	case "interface{}", "any":
		return "any"

	case "time.Time", "Time":
		return "string"

	case "time.Duration", "Duration":
		return "number"

	case "uuid.UUID", "UUID",
		"google.uuid.UUID", // github.com/google/uuid
		"gofrs.uuid.UUID",  // github.com/gofrs/uuid
		"satori.uuid.UUID": // github.com/satori/go.uuid
		return "string"

	case "json.RawMessage", "RawMessage":
		return "any"

	case "sql.NullString", "NullString":
		return "string | null"
	case "sql.NullInt64", "NullInt64", "sql.NullInt32", "NullInt32", "sql.NullInt16", "NullInt16":
		return "number | null"
	case "sql.NullFloat64", "NullFloat64":
		return "number | null"
	case "sql.NullBool", "NullBool":
		return "boolean | null"
	case "sql.NullTime", "NullTime":
		return "string | null"

	case "decimal.Decimal", "Decimal":
		return "string"

	case "big.Int", "big.Float", "big.Rat":
		return "string"

	case "net.IP", "IP":
		return "string"
	case "net.URL", "url.URL", "URL":
		return "string"

	case "[]byte":
		return "string"

	default:
		if idx := strings.LastIndex(goType, "."); idx != -1 {
			return goType[idx+1:]
		}
		return goType
	}
}

// converts a map go type to a TypeScript type
func (g *Generator) mapTypeToTS(goType string) string {
	if !strings.HasPrefix(goType, "map[") {
		return "Record<string, any>"
	}

	depth := 0
	keyEnd := -1
	for i := 4; i < len(goType); i++ {
		if goType[i] == '[' {
			depth++
		} else if goType[i] == ']' {
			if depth == 0 {
				keyEnd = i
				break
			}
			depth--
		}
	}

	if keyEnd == -1 {
		return "Record<string, any>"
	}

	keyType := goType[4:keyEnd]
	valueType := goType[keyEnd+1:]

	tsKeyType := "string"
	switch keyType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		tsKeyType = "number"
	}

	tsValueType := g.basicTypeToTS(valueType)

	return fmt.Sprintf("Record<%s, %s>", tsKeyType, tsValueType)
}
