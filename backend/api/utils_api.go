package api

import (
	"context"
	"github.com/rapid7/go-get-proxied/proxy"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type UtilsApi struct {
	ctx context.Context
}

func NewUtilsApi() *UtilsApi {
	return &UtilsApi{}
}

func (u *UtilsApi) Startup(ctx context.Context) {
	u.ctx = ctx
}

func (u *UtilsApi) GetCpuNum() int {
	return runtime.NumCPU()
}

func (u *UtilsApi) GetUserDownloadPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	downloadPath := path.Join(homeDir, "Downloads")
	downloadPath = filepath.ToSlash(downloadPath)

	return downloadPath, nil
}

func (u *UtilsApi) GetUserProxy() string {
	proxies := proxy.NewProvider("").GetProxies("", "")
	if len(proxies) == 0 {
		return ""
	}

	return proxies[0].URL().String()
}
