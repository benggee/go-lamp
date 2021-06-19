package code

import (
	"github.com/seepre/go-lamp/cli/lamp/svc"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"strings"
)

var configTemplate = `package {{.PackageName}}

import (
	"github.com/seepre/go-lamp/httpz"
)

type Config struct {
	httpz.HttpConf
}

func MustLoad(c *Config) {
	c.Addr = ":8080"
}`

type Config struct {
	ctx *svc.SvcContext
	template string
}

func NewConfig(ctx *svc.SvcContext) Config {
	return Config{
		ctx:      ctx,
		template: configTemplate,
	}
}

func (m Config) Generate(g *types.GeneratorContext) error {
	subPaths := strings.Split(g.SubDir, "/")
	g.ReplaceMap["{{.PackageName}}"] = subPaths[len(subPaths)-1]
	g.Template = m.template
	g.FileName = "config"
	g.FileExt = ".go"


	return m.ctx.FileCtl.Execute(g)
}