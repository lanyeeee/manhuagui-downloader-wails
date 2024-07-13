package export_pdf

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"manhuagui-downloader/backend/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func MergePdfs(pdfDir string, outputPath string) error {
	// 获取目录下的 PDF 文件列表
	pdfEntries, err := getPdfEntries(pdfDir)
	if err != nil {
		return fmt.Errorf("get pdf entries failed: %w", err)
	}
	// 如果目录下没有 PDF 文件，则返回错误
	if len(pdfEntries) == 0 {
		return fmt.Errorf("dir %s has no pdf files", pdfDir)
	}
	// 从 PDF 文件列表中提取 PDF 文件路径
	pdfPaths := make([]string, len(pdfEntries))
	for i, entry := range pdfEntries {
		pdfPaths[i] = path.Join(pdfDir, entry.Name())
	}
	// 合并 PDF 文件
	if err = api.MergeCreateFile(pdfPaths, outputPath, false, nil); err != nil {
		// 删除已经创建的 PDF 文件
		_ = os.Remove(outputPath)
		return fmt.Errorf("merge pdfs failed: %w", err)
	}

	return nil
}

func getPdfEntries(pdfDir string) ([]os.DirEntry, error) {
	// 读取目录下的文件列表
	entries, err := os.ReadDir(pdfDir)
	if err != nil {
		return nil, fmt.Errorf("read dir failed: %w", err)
	}
	// 过滤出 PDF 文件
	var pdfEntries []os.DirEntry
	for _, entry := range entries {
		// 忽略目录或非 PDF 文件
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".pdf" {
			continue
		}

		pdfEntries = append(pdfEntries, entry)
	}
	// 按文件名排序
	sort.Slice(pdfEntries, func(i, j int) bool {
		return utils.FilenameComparer(pdfEntries[i].Name(), pdfEntries[j].Name())
	})

	return pdfEntries, nil
}
