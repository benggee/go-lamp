package gen

import (
	"errors"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"github.com/urfave/cli"
)

func GeneratorCommand(c *cli.Context) error {
	path := c.String("d")
	if len(path) == 0 {
		return errors.New("dir not found, use -d to assign")
	}

	name := c.String("name")
	if len(name) == 0 {
		return errors.New("name invalid, use -name to assign")
	}

	for _, v := range CodeGenerators {
		g := types.GeneratorContext{
			Src:         CODE_SRC,
			RootDir:     path,
			ProjectName: name,
			SubDir:      v.Dir,
			ReplaceMap:  make(map[string]string),
		}
		if err := v.Generator.Generate(&g); err != nil {
			return err
		}
	}

	return nil
}
