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

func (d *DownloadApi) SearchComicById(comicId string, proxyUrl string, cacheDir string) types.Response {
	resp := types.Response{}

	err := http_client.UpdateProxy(proxyUrl)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	comicInfo, err := search.ComicByComicId(comicId, cacheDir)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	resp.Data = comicInfo
	return resp
}

func (d *DownloadApi) SearchComicByKeyword(keyword string, pageNum int, proxyUrl string) types.Response {
	resp := types.Response{}

	err := http_client.UpdateProxy(proxyUrl)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	result, err := search.ComicByKeyword(keyword, pageNum)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	resp.Data = result
	return resp
}

func (d *DownloadApi) DownloadChapter(chapterUrl string, saveDir string, concurrentCount int64, proxyUrl string) types.Response {
	resp := types.Response{}
	err := http_client.UpdateProxy(proxyUrl)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	err = download.ComicChapter(d.ctx, chapterUrl, saveDir, concurrentCount)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}

	return resp
}

func (d *DownloadApi) ComicInfoModel() search.ComicInfo {
	return search.ComicInfo{}
}

func (d *DownloadApi) ComicSearchInfoModel() search.ComicSearchInfo {
	return search.ComicSearchInfo{}
}

func (d *DownloadApi) ComicSearchResultModel() search.ComicSearchResult {
	return search.ComicSearchResult{}
}
