package api

import (
	"context"
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
