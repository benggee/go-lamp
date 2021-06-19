package code

import (
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var demoServiceTemplate = `package {{.PackageName}}

import (
	"context"

	{{.ImportPackages}}
)

type SayHelloRequest struct {
	Content string {{.JsonQuo}}json:"content"{{.JsonQuo}}
}

type SayHelloResponse struct {
	Data string {{.JsonQuo}}json:"data"{{.JsonQuo}}
}

type DemoService struct {
	ctx context.Context
	svcCtx *svc.SvcContext
}

func NewDemoService(ctx context.Context, svcCtx *svc.SvcContext) DemoService {
	return DemoService{
		ctx: ctx,
		svcCtx: svcCtx,
	}
}

func (s *DemoService) SayHello(req SayHelloRequest) (*SayHelloResponse, error) {
	return &SayHelloResponse{
		Data: req.Content,
	}, nil
}`

type DemoService struct {
	ctx *svc.SvcContext
	template string
}

func NewDemoService(ctx *svc.SvcContext) DemoService {
	return DemoService{
		ctx:      ctx,
		template: demoServiceTemplate,
	}
}

func (m DemoService) Generate(g *types.GeneratorContext) error {
	subPaths := strings.Split(g.SubDir, "/")
	g.ReplaceMap["{{.PackageName}}"] = subPaths[len(subPaths)-1]
	g.ReplaceMap["{{.ImportPackages}}"] = m.getPackageImports(g)
	g.ReplaceMap["{{.JsonQuo}}"] = "`"
	g.Template = m.template
	g.FileName = "demo"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}

func (m DemoService) getPackageImports(g *types.GeneratorContext) string {
	packages := []string{
		fmt.Sprintf("\"%s/%s\"", g.ProjectName, "internal/svc"),
	}

	return strings.Join(packages, "\r\n")
}
