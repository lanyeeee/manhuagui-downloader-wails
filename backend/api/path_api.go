package api

import (
	"context"
	"manhuagui-downloader/backend/utils"
	"os"
	"path"
	"path/filepath"
)

type PathApi struct {
	ctx context.Context
}

func NewPathApi() *PathApi {
	return &PathApi{}
}

func (p *PathApi) Startup(ctx context.Context) {
	p.ctx = ctx
}

func (p *PathApi) GetAbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	abs = filepath.ToSlash(abs)
	return abs, err
}

func (p *PathApi) PathExists(path string) bool {
	return utils.PathExists(path)
}

func (p *PathApi) GetRelPath(cacheDir string, path string) (string, error) {
	rel, err := filepath.Rel(cacheDir, path)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	return rel, nil
}

func (p *PathApi) Join(args ...interface{}) string {
	params := args[0].([]interface{})
	ss := make([]string, len(params))

	for _, param := range params {
		ss = append(ss, param.(string))
	}
	return filepath.ToSlash(path.Join(ss...))
}

func (p *PathApi) MkDirAll(path string) error {
	return os.MkdirAll(path, 0777)
}
