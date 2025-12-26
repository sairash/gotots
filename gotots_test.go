package gotots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatal("New() returned nil")
	}
	if g.gen == nil {
		t.Fatal("New() returned generator with nil internal generator")
	}
}

func TestGeneratorChaining(t *testing.T) {
	g := New()

	result := g.FromDir("testdir")
	if result != g {
		t.Error("FromDir should return the same generator for chaining")
	}

	result = g.ToFile("output.ts")
	if result != g {
		t.Error("ToFile should return the same generator for chaining")
	}
}

func TestGenerateWithoutInputDir(t *testing.T) {
	g := New().ToFile("output.ts")
	err := g.Generate()
	if err == nil {
		t.Error("Generate should fail without input directory")
	}
	if !strings.Contains(err.Error(), "input directory") {
		t.Errorf("Error should mention input directory, got: %v", err)
	}
}

func TestGenerateWithoutOutputFile(t *testing.T) {
	// Create temp dir with a Go file
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type User struct {
	Name string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	g := New().FromDir(tmpDir)
	err = g.Generate()
	if err == nil {
		t.Error("Generate should fail without output file")
	}
	if !strings.Contains(err.Error(), "output file") {
		t.Errorf("Error should mention output file, got: %v", err)
	}
}

func TestGenerateWithNonExistentDir(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "types.ts")

	g := New().FromDir("/nonexistent/path/to/models").ToFile(outputFile)
	err := g.Generate()
	if err == nil {
		t.Error("Generate should fail with non-existent directory")
	}
}

func TestGenerateBasicStruct(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go file with a simple struct
	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type User struct {
	ID   int
	Name string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "export interface User") {
		t.Error("Output should contain 'export interface User'")
	}
	if !strings.Contains(output, "ID: number") {
		t.Error("Output should contain 'ID: number'")
	}
	if !strings.Contains(output, "Name: string") {
		t.Error("Output should contain 'Name: string'")
	}
}

func TestGenerateWithPointerField(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Profile struct {
	Bio     *string
	Website *string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "Bio?:") {
		t.Error("Pointer fields should be optional (Bio?:)")
	}
	if !strings.Contains(output, "string | null") {
		t.Error("Pointer fields should have '| null' type")
	}
}

func TestGenerateWithSliceField(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Post struct {
	Tags     []string
	Comments []int
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "string[]") {
		t.Error("Output should contain 'string[]' for []string")
	}
	if !strings.Contains(output, "number[]") {
		t.Error("Output should contain 'number[]' for []int")
	}
}

func TestGenerateWithMapField(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Config struct {
	Settings map[string]string
	Counts   map[string]int
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "Record<string, string>") {
		t.Error("Output should contain 'Record<string, string>'")
	}
	if !strings.Contains(output, "Record<string, number>") {
		t.Error("Output should contain 'Record<string, number>'")
	}
}

func TestGenerateWithNestedStruct(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Address struct {
	Street string
	City   string
}

type Person struct {
	Name    string
	Address Address
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "export interface Address") {
		t.Error("Output should contain 'export interface Address'")
	}
	if !strings.Contains(output, "export interface Person") {
		t.Error("Output should contain 'export interface Person'")
	}
	if !strings.Contains(output, "Address: Address") {
		t.Error("Output should reference Address type in Person")
	}
}

func TestGenerateSkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create regular Go file
	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type User struct {
	Name string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create test file that should be skipped
	testFile := filepath.Join(tmpDir, "model_test.go")
	err = os.WriteFile(testFile, []byte(`package models

type TestHelper struct {
	Data string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "export interface User") {
		t.Error("Output should contain 'export interface User'")
	}
	if strings.Contains(output, "TestHelper") {
		t.Error("Output should NOT contain 'TestHelper' from test file")
	}
}

func TestGenerateMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create first Go file
	file1 := filepath.Join(tmpDir, "user.go")
	err := os.WriteFile(file1, []byte(`package models

type User struct {
	ID   int
	Name string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create second Go file
	file2 := filepath.Join(tmpDir, "post.go")
	err = os.WriteFile(file2, []byte(`package models

type Post struct {
	Title   string
	Content string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "export interface User") {
		t.Error("Output should contain 'export interface User'")
	}
	if !strings.Contains(output, "export interface Post") {
		t.Error("Output should contain 'export interface Post'")
	}
}

func TestGenerateSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "entities")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create Go file in root
	rootFile := filepath.Join(tmpDir, "user.go")
	err = os.WriteFile(rootFile, []byte(`package models

type User struct {
	Name string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create Go file in subdirectory
	subFile := filepath.Join(subDir, "product.go")
	err = os.WriteFile(subFile, []byte(`package entities

type Product struct {
	Price float64
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "export interface User") {
		t.Error("Output should contain 'export interface User' from root")
	}
	if !strings.Contains(output, "export interface Product") {
		t.Error("Output should contain 'export interface Product' from subdirectory")
	}
}

func TestGenerateAllNumericTypes(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Numbers struct {
	A int
	B int8
	C int16
	D int32
	E int64
	F uint
	G uint8
	H uint16
	I uint32
	J uint64
	K float32
	L float64
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	// All numeric fields should be 'number' type
	lines := strings.Split(output, "\n")
	numberCount := 0
	for _, line := range lines {
		if strings.Contains(line, ": number;") {
			numberCount++
		}
	}
	if numberCount != 12 {
		t.Errorf("Expected 12 number fields, got %d", numberCount)
	}
}

func TestGenerateBoolType(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Flags struct {
	Active  bool
	Deleted bool
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "Active: boolean") {
		t.Error("bool should map to boolean")
	}
	if !strings.Contains(output, "Deleted: boolean") {
		t.Error("bool should map to boolean")
	}
}

func TestGenerateOutputHeader(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Empty struct {
	Field string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.HasPrefix(output, "/* Do not change, this code is generated from Golang structs */") {
		t.Error("Output should start with generation warning comment")
	}
}

func TestGenerateWithJSONTag(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type User struct {
	ID        int    `+"`json:\"id\"`"+`
	FirstName string `+"`json:\"first_name\"`"+`
	LastName  string `+"`json:\"lastName\"`"+`
	Email     string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "id: number") {
		t.Error("Output should use JSON tag 'id' instead of 'ID'")
	}
	if !strings.Contains(output, "first_name: string") {
		t.Error("Output should use JSON tag 'first_name' instead of 'FirstName'")
	}
	if !strings.Contains(output, "lastName: string") {
		t.Error("Output should use JSON tag 'lastName' instead of 'LastName'")
	}
	if !strings.Contains(output, "Email: string") {
		t.Error("Output should keep 'Email' when no JSON tag is present")
	}
}

func TestGenerateWithJSONTagOmitEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Profile struct {
	Bio      string `+"`json:\"bio,omitempty\"`"+`
	Website  string `+"`json:\"website,omitempty\"`"+`
	Location string `+"`json:\"location\"`"+`
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "bio?: string") {
		t.Error("Field with omitempty should be optional (bio?:)")
	}
	if !strings.Contains(output, "website?: string") {
		t.Error("Field with omitempty should be optional (website?:)")
	}
	if !strings.Contains(output, "location: string") {
		t.Error("Field without omitempty should not be optional")
	}
	if strings.Contains(output, "location?:") {
		t.Error("Field without omitempty should not have optional marker")
	}
}

func TestGenerateWithJSONTagSkip(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type Secret struct {
	PublicKey  string `+"`json:\"public_key\"`"+`
	PrivateKey string `+"`json:\"-\"`"+`
	Internal   string `+"`json:\"-\"`"+`
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "public_key: string") {
		t.Error("Output should contain 'public_key' field")
	}
	// Fields with json:"-" should use the original field name (not be skipped entirely)
	if !strings.Contains(output, "PrivateKey: string") {
		t.Error("Field with json:\"-\" should use original field name 'PrivateKey'")
	}
	if !strings.Contains(output, "Internal: string") {
		t.Error("Field with json:\"-\" should use original field name 'Internal'")
	}
}

func TestGenerateWithJSONTagAndPointer(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "model.go")
	err := os.WriteFile(goFile, []byte(`package models

type OptionalFields struct {
	Name     *string `+"`json:\"name\"`"+`
	Age      *int    `+"`json:\"age,omitempty\"`"+`
	Nickname *string
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(tmpDir, "types.ts")
	err = New().FromDir(tmpDir).ToFile(outputFile).Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	output := string(content)
	// Pointer fields are already optional
	if !strings.Contains(output, "name?: string | null") {
		t.Error("Pointer field should be optional and nullable with JSON tag name")
	}
	if !strings.Contains(output, "age?: number | null") {
		t.Error("Pointer field with omitempty should be optional and nullable")
	}
	if !strings.Contains(output, "Nickname?: string | null") {
		t.Error("Pointer field without JSON tag should use original name")
	}
}
