package export_pdf

import (
	"errors"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"manhuagui-downloader/backend/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
)

type CreatePdfsRequest struct {
	Tasks           []CreatePdfTask `json:"tasks"`
	ConcurrentCount int64           `json:"concurrentCount"`
}

type CreatePdfTask struct {
	ImgDir     string `json:"imgDir"`
	OutputPath string `json:"outputPath"`
	OptionKey  string `json:"optionKey"`
}

type createPdfResult struct {
	completedOptionKey string
	err                error
}

func CreatePdfs(ctx context.Context, request CreatePdfsRequest) error {
	concurrentCount := request.ConcurrentCount
	// 创建一个通道，用于传输创建PDF的结果
	pdfResultCh := make(chan createPdfResult, concurrentCount)
	// 启动一个生产者goroutine，创建PDF
	go func() {
		wg := sync.WaitGroup{}
		sem := semaphore.NewWeighted(concurrentCount)
		for i := range request.Tasks {
			wg.Add(1)
			go createPdf(ctx, &request.Tasks[i], pdfResultCh, sem, &wg)
		}
		wg.Wait()
		close(pdfResultCh)
	}()
	// 当前goroutine作为消费者，等待生产者goroutine创建PDF
	totalTaskCount := len(request.Tasks) // 总共需要创建的 PDF 数量
	completedTaskCount := 0              // 已经成功创建的 PDF 数量
	for result := range pdfResultCh {
		// 某个任务创建 PDF 失败，返回错误
		if result.err != nil {
			return fmt.Errorf("create pdf failed: %w", result.err)
		}
		// 某个任务创建 PDF 成功，更新进度
		completedTaskCount++

		msg := fmt.Sprintf("(%d/%d)", completedTaskCount, totalTaskCount)
		percentage := float64(completedTaskCount) / float64(totalTaskCount) * 100

		runtime.EventsEmit(ctx, "create_pdf", result.completedOptionKey, msg, percentage)
	}

	return nil
}

// TODO: 支持控制纸张大小
func createPdf(ctx context.Context, task *CreatePdfTask, pdfResultCh chan<- createPdfResult, sem *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	// 获取图片文件列表
	imgEntries, err := getImgEntries(task.ImgDir)
	if err != nil {
		pdfResultCh <- createPdfResult{"", fmt.Errorf("get img entries failed: %w", err)}
		return
	}
	// 如果没有图片文件，则返回错误
	if len(imgEntries) == 0 {
		err = errors.New(fmt.Sprintf("dir '%s' has no image files", task.ImgDir))
		pdfResultCh <- createPdfResult{"", err}
		return
	}
	// 从文件列表中获取图片文件路径
	imgPaths := make([]string, len(imgEntries))
	for i, entry := range imgEntries {
		imgPaths[i] = path.Join(task.ImgDir, entry.Name())
	}

	if err = sem.Acquire(ctx, 1); err != nil {
		pdfResultCh <- createPdfResult{"", fmt.Errorf("acquire semaphore failed: %w", err)}
		return
	}
	defer sem.Release(1)
	// 获取导出的目录，并创建目录
	dir, _ := filepath.Split(task.OutputPath)
	if err = os.MkdirAll(dir, 0777); err != nil {
		pdfResultCh <- createPdfResult{"", fmt.Errorf("create dir failed: %w", err)}
		return
	}
	// 创建 PDF 导出选项
	imp, err := api.Import("form:A4, pos:c, scale:1.0", types.POINTS)
	if err != nil {
		pdfResultCh <- createPdfResult{"", fmt.Errorf("pdfcup import failed: %w", err)}
		return
	}
	// 将图片文件导入到PDF文件中
	if err = api.ImportImagesFile(imgPaths, task.OutputPath, imp, nil); err != nil {
		// 删除导出的文件
		_ = os.Remove(task.OutputPath)
		pdfResultCh <- createPdfResult{"", fmt.Errorf("pdfcup import images file failed: %w", err)}
		return
	}

	pdfResultCh <- createPdfResult{task.OptionKey, nil}
}

func getImgEntries(imgDir string) ([]os.DirEntry, error) {
	// 读取目录下的文件列表
	entries, err := os.ReadDir(imgDir)
	if err != nil {
		return []os.DirEntry{}, fmt.Errorf("read dir failed: %w", err)
	}
	// 过滤出图片文件
	imgEntries := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		// 忽略目录
		if entry.IsDir() {
			continue
		}
		// 忽略非图片文件
		if path.Ext(entry.Name()) != ".jpg" && path.Ext(entry.Name()) != ".jpeg" && path.Ext(entry.Name()) != ".png" {
			continue
		}

		imgEntries = append(imgEntries, entry)
	}
	// 按文件名排序
	sort.Slice(imgEntries, func(i, j int) bool {
		return utils.FilenameComparer(imgEntries[i].Name(), imgEntries[j].Name())
	})

	return imgEntries, nil
}
