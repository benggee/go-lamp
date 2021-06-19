package code

import (
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var mainTemplate = `package {{.PackageName}} 

import (
	"fmt"

	"github.com/seepre/go-lamp/httpz"
	{{.ImportPackages}}
)

func main() {
	var c config.Config

	config.MustLoad(&c)

	ctx := svc.NewSvcContext(c)

	serve := httpz.MustNewServe(c.HttpConf)
	defer serve.Stop()

	handler.RegisterRoute(serve, ctx)

	fmt.Printf("starting server at: %s...\n", c.Addr)

	serve.Run()
}`

type Main struct {
	ctx *svc.SvcContext
	template string
}

func NewMain(ctx *svc.SvcContext) Main {
	return Main{
		ctx:      ctx,
		template: mainTemplate,
	}
}

func (m Main) Generate(g *types.GeneratorContext) error {
	g.ReplaceMap["{{.PackageName}}"] = "main"
	g.ReplaceMap["{{.ImportPackages}}"] = m.getPackageImports(g)
	g.Template = m.template
	g.FileName = "main"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}

func (m Main) getPackageImports(g *types.GeneratorContext) string {
	packages := []string{
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/config"),
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/handler"),
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/svc"),
	}

	return strings.Join(packages, "\r\n")
}
