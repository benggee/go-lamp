package code

import (
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
)

var modTemplate = `module {{.Module}}

go 1.16

require (
	github.com/seepre/go-lamp v1.0
)`

type Mod struct {
	ctx *svc.SvcContext
	template string
}

func NewMod(ctx *svc.SvcContext) Mod {
	return Mod{
		ctx:      ctx,
		template: modTemplate,
	}
}

func (m Mod) Generate(g *types.GeneratorContext) error {
	g.Template = m.template
	g.FileName = "go"
	g.FileExt = ".mod"
	g.ReplaceMap["{{.Module}}"] = g.ProjectName

	return m.ctx.FileCtl.Execute(g)
}

