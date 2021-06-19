package code

import (
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var svcContextTemplate = `package {{.PackageName}} 

import (
	{{.ImportPackages}}
)

type SvcContext struct {
	Config config.Config
}

func NewSvcContext(c config.Config) *SvcContext {
	return &SvcContext{
		Config: c,
	}
}`

type SvcContext struct {
	ctx *svc.SvcContext
	template string
}

func NewSvcContext(ctx *svc.SvcContext) SvcContext {
	return SvcContext{
		ctx:      ctx,
		template: svcContextTemplate,
	}
}

func (m SvcContext) Generate(g *types.GeneratorContext) error {
	subPaths := strings.Split(g.SubDir, "/")
	g.ReplaceMap["{{.PackageName}}"] = subPaths[len(subPaths)-1]
	g.ReplaceMap["{{.ImportPackages}}"] = m.getPackageImports(g)
	g.Template = m.template
	g.FileName = "route"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}

func (m SvcContext) getPackageImports(g *types.GeneratorContext) string {
	packages := []string{
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/config"),
	}

	return strings.Join(packages, "\r\n")
}
