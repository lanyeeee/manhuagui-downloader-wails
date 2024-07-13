package api

import (
	"context"
	"fmt"
	"manhuagui-downloader/backend/export_pdf"
	"manhuagui-downloader/backend/scan_cache"
	"manhuagui-downloader/backend/types"
)

type ExportApi struct {
	ctx context.Context
}

func NewExportApi() *ExportApi {
	return &ExportApi{}
}

func (e *ExportApi) Startup(ctx context.Context) {
	e.ctx = ctx
}

func (e *ExportApi) ScanCacheDir(cacheDir string, exportDir string, maxDepth int64) types.Response {
	resp := types.Response{}
	treeOption, err := scan_cache.ScanCacheDir(cacheDir, exportDir, maxDepth)
	if err != nil {
		resp.Code = -1
		resp.Msg = fmt.Sprintf("ScanCacheDir: %s", err.Error())
		return resp
	}

	resp.Data = treeOption
	return resp
}

func (e *ExportApi) CreatePdfs(request export_pdf.CreatePdfsRequest) types.Response {
	resp := types.Response{}
	err := export_pdf.CreatePdfs(e.ctx, request)

	if err != nil {
		resp.Code = -1
		resp.Msg = fmt.Sprintf("CreatePdfs: %s", err.Error())
	}

	return resp
}

func (e *ExportApi) MergePdfs(pdfDir string, outputPath string) types.Response {
	resp := types.Response{}
	err := export_pdf.MergePdfs(pdfDir, outputPath)

	if err != nil {
		resp.Code = -1
		resp.Msg = fmt.Sprintf("MergePdfs: %s", err.Error())
	}

	return resp
}

func (e *ExportApi) TreeOptionModel() types.TreeNode {
	return types.TreeNode{}
}
