package code

import (
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var routeTemplate = `package {{.PackageName}} 

import (
	"net/http"

	"github.com/seepre/go-lamp/httpz"
	"github.com/seepre/go-lamp/httpz/router"
	{{.ImportPackages}}
)

func RegisterRoute(serve *httpz.Serve, ctx *svc.SvcContext) {
	serve.AddRoutes([]router.Route{
		{
			Method: http.MethodPost,
			Path: "/v1.0/sayHello",
			Handler: demo.HelloWorld(ctx),
		},
	})
}`

type Route struct {
	ctx *svc.SvcContext
	template string
}

func NewRoute(ctx *svc.SvcContext) Route {
	return Route{
		ctx:      ctx,
		template: routeTemplate,
	}
}

func (m Route) Generate(g *types.GeneratorContext) error {
	subPaths := strings.Split(g.SubDir, "/")
	g.ReplaceMap["{{.PackageName}}"] = subPaths[len(subPaths)-1]
	g.ReplaceMap["{{.ImportPackages}}"] = m.getPackageImports(g)
	g.Template = m.template
	g.FileName = "route"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}

func (m Route) getPackageImports(g *types.GeneratorContext) string {
	packages := []string{
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/svc"),
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/handler/demo"),
	}

	return strings.Join(packages, "\r\n")
}
