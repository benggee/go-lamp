package svc

import (
	"github.com/seepre/go-lamp/cli/lamp/file"
)

type SvcContext struct {
	FileCtl *file.File
}

func NewSvc() *SvcContext {
	return &SvcContext{
		FileCtl: file.NewFile(),
	}
}