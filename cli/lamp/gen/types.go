package gen

import (
	"github.com/seepre/go-lamp/cli/lamp/code"
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
)

const (
	CODE_SRC = "github.com/seepre"
)


var svcCtx = svc.NewSvc()

var CodeGenerators = []types.GeneratorItem{
	{
		Dir:       "/",
		Generator: code.NewMain(svcCtx),
	},
	{
		Dir:       "/internal/config",
		Generator: code.NewConfig(svcCtx),
	},
	{
		Dir:       "/internal/handler",
		Generator: code.NewRoute(svcCtx),
	},
	{
		Dir:       "/internal/handler/demo",
		Generator: code.NewDemoHandler(svcCtx),
	},
	{
		Dir:       "/internal/service",
		Generator: code.NewDemoService(svcCtx),
	},
	{
		Dir:       "/",
		Generator: code.NewMod(svcCtx),
	},
}


