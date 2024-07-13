package download

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sync/semaphore"
	"io"
	"manhuagui-downloader/backend/decoder"
	"manhuagui-downloader/backend/http_client"
	"manhuagui-downloader/backend/utils"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type downloadResult struct {
	imgData *[]byte
	err     error
}

func ComicChapter(ctx context.Context, chapterUrl string, saveDir string, concurrentCount int64) error {
	// 处理saveDir，去掉特殊字符
	dir, file := path.Split(saveDir)
	file = utils.Sanitize(file)
	// 获取保存目录和临时目录
	saveDir = path.Join(dir, file)
	tempDir := path.Join(dir, "."+file)
	// 创建临时目录
	if err := os.MkdirAll(tempDir, 0777); err != nil {
		return fmt.Errorf("create temp dir failed: %w", err)
	}
	// 解析章节页面，获取图片列表
	decodeResult, err := requestDecodeResult(chapterUrl)
	if err != nil {
		// TODO: 处理EOF错误(IP被封)
		return fmt.Errorf("request decode result failed: %w", err)
	}
	// 没有图片就直接返回
	if len(decodeResult.Files) == 0 {
		return nil
	}
	// 通过解析结果获取图片url列表
	imgUrls := make([]string, len(decodeResult.Files))
	for i, fileName := range decodeResult.Files {
		// 去掉.webp后缀，剩下的就是.jpg的文件名
		fileName = strings.TrimSuffix(fileName, ".webp")
		imgUrls[i] = "https://i.hamreus.com" + decodeResult.Path + fileName
	}
	// 创建一个通道，用于传输下载的图片数据
	downloadResultCh := make(chan downloadResult, concurrentCount)
	// 启动一个生产者goroutine，下载图片
	go func() {
		wg := sync.WaitGroup{}
		sem := semaphore.NewWeighted(concurrentCount)
		// 并发下载
		for i, imgUrl := range imgUrls {
			indexLength := len(strconv.Itoa(len(imgUrls)))
			filename := fmt.Sprintf("%0*d", indexLength, i) + ".jpg"
			dstPath := path.Join(tempDir, filename)
			// 已经存在的文件不再下载
			if utils.PathExists(dstPath) {
				fmt.Printf("%s 已存在，跳过下载\n", dstPath)
				continue
			}

			wg.Add(1)
			go downloadImage(ctx, imgUrl, dstPath, downloadResultCh, sem, &wg)
		}
		wg.Wait()
		close(downloadResultCh)
	}()

	// 当前goroutine是消费者，接收下载的图片数据
	start := time.Now()
	downloadedBytes := 0          // 已下载的字节数
	imgDownloadedCount := 0       // 已下载的图片数量
	totalImgCount := len(imgUrls) // 总共需要下载的图片数量
	for result := range downloadResultCh {
		// 某个图片下载失败
		if result.err != nil {
			return fmt.Errorf("download image failed: %w", result.err)
		}
		// 某个图片下载成功，更新进度
		imgDownloadedCount++
		downloadedBytes += len(*result.imgData)
		elapsed := time.Since(start)
		mbPerSecond := float64(downloadedBytes) / 1024 / 1024 / elapsed.Seconds()

		msg := fmt.Sprintf("%s (%d/%d)：%.2f MB/s", file, imgDownloadedCount, totalImgCount, mbPerSecond)
		percentage := float64(imgDownloadedCount) / float64(totalImgCount) * 100

		runtime.EventsEmit(ctx, "download", msg, percentage)
	}

	// 如果saveDir已存在，删除
	if err = os.RemoveAll(saveDir); err != nil {
		return fmt.Errorf("remove save dir failed: %w", err)
	}
	// 将临时目录改名为saveDir
	if err = os.Rename(tempDir, saveDir); err != nil {
		return fmt.Errorf("rename temp dir to save dir failed: %w", err)
	}

	return nil
}

func downloadImage(ctx context.Context, imgUrl string, dstPath string, downloadResultCh chan<- downloadResult, sem *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := sem.Acquire(ctx, 1); err != nil {
		downloadResultCh <- downloadResult{imgData: nil, err: fmt.Errorf("acquire semaphore failed: %w", err)}
		return
	}
	defer sem.Release(1)
	// 最多重试3次
	const MaxRetry = 3
	var imgData *[]byte
	var err error
	// 带重试下载图片
	for i := 0; i < MaxRetry; i++ {
		imgData, err = requestImageData(imgUrl)
		// 下载成功则退出循环
		if err == nil {
			break
		}
		// TODO: 处理EOF错误(IP被封)
		// 下载失败则等待1秒后重试
		time.Sleep(1 * time.Second)
	}
	// 下载失败
	if err != nil {
		fmt.Printf("下载图片 %s 失败，错误信息：\n%s\n", imgUrl, err)
		downloadResultCh <- downloadResult{imgData: nil, err: fmt.Errorf("download image failed: %w", err)}
		return
	}
	// 下载失败
	if imgData == nil {
		downloadResultCh <- downloadResult{imgData: nil, err: fmt.Errorf("download image failed: imgData is nil")}
		return
	}
	// 下载成功，保存图片
	if err = os.WriteFile(dstPath, *imgData, 0644); err != nil {
		fmt.Printf("保存图片 %s 失败，错误信息：\n%s\n", dstPath, err)
		downloadResultCh <- downloadResult{imgData: nil, err: fmt.Errorf("save image failed: %w", err)}
		return
	}

	downloadResultCh <- downloadResult{imgData: imgData, err: nil}
}

func requestImageData(imgUrl string) (*[]byte, error) {
	req, err := http.NewRequest("GET", imgUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Referer", "https://www.manhuagui.com/")

	resp, err := http_client.HttpClientInst().Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	return &respBody, nil
}

func requestDecodeResult(chapterUrl string) (decoder.DecodeResult, error) {
	resp, err := http_client.HttpClientInst().Get(chapterUrl)
	if err != nil {
		return decoder.DecodeResult{}, fmt.Errorf("do request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return decoder.DecodeResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decoder.DecodeResult{}, fmt.Errorf("read response body failed: %w", err)
	}

	htmlContent := string(respBody)

	result, err := decoder.Decode(&htmlContent)
	if err != nil {
		return decoder.DecodeResult{}, fmt.Errorf("decode failed: %w", err)
	}

	return result, nil
}
