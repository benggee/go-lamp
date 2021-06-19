package types

type GeneratorContext struct {
	Src string
	RootDir string
	ProjectName string
	SubDir string
	ReplaceMap map[string]string
	Template string
	FileName string
	FileExt string
}

type GeneratorItem struct {
	Dir       string
	Generator Generator
}

type Generator interface {
	Generate(g *GeneratorContext) error
}