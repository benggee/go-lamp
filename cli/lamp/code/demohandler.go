package code

import (
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var demoHandlerTemplate = `package {{.PackageName}}

import (
	"net/http"

	"github.com/seepre/go-lamp/httpz"
	{{.ImportPackages}}
)

func HelloWorld(ctx *svc.SvcContext) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var req service.SayHelloRequest
		if err := httpz.Parse(r, &req); err != nil {
			httpz.Error(w, err)
			return
		}

		s := service.NewDemoService(r.Context(), ctx)
		resp, err := s.SayHello(req)
		if err != nil {
			httpz.Error(w, err)
			return
		}
		httpz.SuccessJson(w, resp)
	}
}`

type DemoHandler struct {
	ctx *svc.SvcContext
	template string
}

func NewDemoHandler(ctx *svc.SvcContext) DemoHandler {
	return DemoHandler{
		ctx:      ctx,
		template: demoHandlerTemplate,
	}
}

func (m DemoHandler) Generate(g *types.GeneratorContext) error {
	subPaths := strings.Split(g.SubDir, "/")
	g.ReplaceMap["{{.PackageName}}"] = subPaths[len(subPaths)-1]
	g.ReplaceMap["{{.ImportPackages}}"] = m.getPackageImports(g)
	g.Template = m.template
	g.FileName = "demo"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}

func (m DemoHandler) getPackageImports(g *types.GeneratorContext) string {
	packages := []string{
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/svc"),
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/service"),
	}

	return strings.Join(packages, "\r\n")
}