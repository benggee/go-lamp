package file

import (
	"errors"
	"fmt"
	"github.com/seepre/go-lamp/cli/lamp/types"
	"go/format"
	"os"
	"strings"
)

type File struct {
}

func NewFile() *File {
	return &File{}
}


func (f *File) Execute(g *types.GeneratorContext) error {
	if len(g.FileName) == 0 || len(g.FileExt) == 0 {
		return errors.New("file name or ext is invalid.")
	}

	for k, v := range g.ReplaceMap {
		g.Template = strings.ReplaceAll(g.Template, k, v)
	}

	ret, err := format.Source([]byte(g.Template))
	if err == nil {
		g.Template = string(ret)
	}

	return f.writeToFile(fmt.Sprintf("%s%s", g.FileName, g.FileExt),g.Template, g.RootDir, g.ProjectName, g.SubDir)
}


func (f *File) writeToFile(fileName, data string, pathItems ...string ) error {
	var (
		file *os.File
		err error
	)
	defer func() {
		file.Close()
	}()

	file, err = f.createFile(fileName, pathItems...)
	if err != nil {
		return err
	}

	_, err = file.WriteString(data)
	return err
}

func (f *File) createFile(fileName string, pathItems ...string) (*os.File, error) {
	path, err  := f.buildPath(pathItems...)
	if err != nil {
		return nil, err
	}

	if err = f.exist(path); err != nil {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	filePath := fmt.Sprintf("%s%s", path, fileName)
	if err = f.exist(filePath); err != nil {
		return os.Create(filePath)
	}

	return os.OpenFile(filePath, 1, os.ModePerm)
}

func (f *File) exist(path string) error {
	_, err := os.Stat(path)
	return err
}

func (f *File) buildPath(pathItems ...string) (string, error) {
	if len(pathItems) == 0 {
		return "", errors.New("path invalid")
	}

	tmpPaths := make([]string, 0)

	for _, v := range pathItems {
		newPathItem := f.cleanPath(v)
		if len(newPathItem) == 0 {
			continue
		}
		tmpPaths = append(tmpPaths, newPathItem)
	}

	return strings.Join(tmpPaths, "")+"/", nil
}

func (f *File) cleanPath(pathItem string) string {
	if pathItem == "./" {
		return pathItem
	}

	pathBytes := []byte(pathItem)
	for len(pathBytes) > 0 {
		if pathBytes[len(pathBytes)-1] == '/' {
			pathBytes = pathBytes[:len(pathBytes)-1]
		} else {
			break
		}
	}

	if len(pathBytes) == 0 {
		return ""
	}

	for len(pathBytes) > 1 {
		if pathBytes[1] == '/' {
			pathBytes = pathBytes[1:]
		} else {
			break
		}
	}

	return string(pathBytes)
}