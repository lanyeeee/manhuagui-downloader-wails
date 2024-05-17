package api

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"manhuagui-downloader/backend/types"
	"manhuagui-downloader/backend/utils"
	"path/filepath"
)

type SettingsApi struct {
	ctx context.Context
}

func NewSettingsApi() *SettingsApi {
	return &SettingsApi{}
}

func (s *SettingsApi) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *SettingsApi) ChooseDirectory(dirPath string) types.Response {
	resp := types.Response{}
	// 如果目录不存在，则打开默认目录
	if !utils.PathExists(dirPath) {
		dirPath = ""
	}

	option := runtime.OpenDialogOptions{
		DefaultDirectory: dirPath,
		Title:            "选择目录",
	}
	// 打开目录选择对话框
	chosenDir, err := runtime.OpenDirectoryDialog(s.ctx, option)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	resp.Data = filepath.ToSlash(chosenDir)
	return resp
}
