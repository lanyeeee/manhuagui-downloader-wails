package api

import (
	"context"
	"manhuagui-downloader/backend/download"
	"manhuagui-downloader/backend/http_client"
	"manhuagui-downloader/backend/search"
	"manhuagui-downloader/backend/types"
)

type DownloadApi struct {
	ctx context.Context
}

func NewDownloadApi() *DownloadApi {
	return &DownloadApi{}
}

func (d *DownloadApi) Startup(ctx context.Context) {
	d.ctx = ctx
}

func (d *DownloadApi) SearchComicInfo(comicId string, proxyUrl string) types.Response {
	resp := types.Response{}
	err := http_client.UpdateProxy(proxyUrl)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return types.Response{}
	}

	comicInfo, err := search.Info(comicId)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}

	resp.Data = comicInfo
	return resp
}

func (d *DownloadApi) DownloadChapter(chapterUrl string, saveDir string, concurrentCount int64, proxyUrl string) types.Response {
	resp := types.Response{}
	err := http_client.UpdateProxy(proxyUrl)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return types.Response{}
	}

	err = download.ComicChapter(d.ctx, chapterUrl, saveDir, concurrentCount)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
	}

	return resp
}

func (d *DownloadApi) ComicInfoModel() search.ComicInfo {
	return search.ComicInfo{}
}
