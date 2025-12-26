package gotots

import "github.com/sairash/gotots/internal"

type Generator struct {
	gen *internal.Generator
}

func New() *Generator {
	return &Generator{
		gen: internal.New(),
	}
}

func (g *Generator) FromDir(dir string) *Generator {
	g.gen.FromDir(dir)
	return g
}

func (g *Generator) ToFile(file string) *Generator {
	g.gen.ToFile(file)
	return g
}

func (g *Generator) Generate() error {
	return g.gen.Generate()
}
