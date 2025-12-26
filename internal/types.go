package internal

// represents an enum (TODO)
type EnumInfo struct {
	Name    string
	Package string
	Values  []EnumValue
}

// represents a value of an enum (TODO)
type EnumValue struct {
	Name  string
	Value int
}

// represents a struct
type StructInfo struct {
	Name    string
	Package string
	Fields  []FieldInfo
}

// represents a field of a struct
type FieldInfo struct {
	Name       string
	Type       string
	JSONTag    string
	IsEnum     bool
	EnumType   string
	IsOptional bool
	IsPointer  bool
}
