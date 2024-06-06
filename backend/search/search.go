package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	lzstring "github.com/daku10/go-lz-string"
	"io"
	"manhuagui-downloader/backend/http_client"
	"manhuagui-downloader/backend/types"
	"manhuagui-downloader/backend/utils"
	"net/http"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// ComicInfo 漫画信息，包含 漫画标题 和 章节类型(单话、单行本、番外)
type ComicInfo struct {
	Title        string        `json:"title"`
	ChapterTypes []ChapterType `json:"chapterTypes"`
}

// ChapterType 章节类型，包含 章节类型标题 和 章节分页(第1-10页、第11-20页)
type ChapterType struct {
	Title        string        `json:"title"`
	ChapterPages []ChapterPage `json:"chapterPages"`
}

// ChapterPage 分页信息，包含分页标题(第1-10页、第11-20页) 和 章节列表
type ChapterPage struct {
	Title    string    `json:"title"`
	Chapters []Chapter `json:"chapters"`
}

// Chapter 章节信息，包含 章节标题 和 章节链接
type Chapter struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

// ChapterTreeNodeKey 章节树节点的Key，包含 章节链接 和 保存目录
type ChapterTreeNodeKey struct {
	Href    string `json:"href"`
	SaveDir string `json:"saveDir"`
}

// ComicSearchResult 漫画搜索结果，包含 漫画标题、作者 和 漫画ID
type ComicSearchResult struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	ComicId string   `json:"comicId"`
}

func ComicByComicId(comicId string, cacheDir string) (types.TreeNode, error) {
	resp, err := http_client.HttpClientInst().Get("https://www.manhuagui.com/comic/" + comicId)
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("do request failed: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("read response body failed: %w", err)
	}
	// 处理HTTP错误
	switch resp.StatusCode {
	case http.StatusOK:
		// ignore
	case http.StatusNotFound:
		return types.TreeNode{}, errors.New(fmt.Sprintf("can't find comic with id: %s", comicId))
	default:
		return types.TreeNode{}, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	htmlContent := string(respBody)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("parse html failed: %w", err)
	}

	title, err := getTitle(doc)
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("get title failed: %w", err)
	}
	warningBar := doc.Find("div[class=warning-bar]")
	// 如果是带警告的漫画
	if warningBar.Length() > 0 {
		// 获取id为__VIEWSTATE的input标签的value属性
		val, exists := doc.Find("input[id=__VIEWSTATE]").First().Attr("value")
		if !exists {
			return types.TreeNode{}, errors.New("can't find __VIEWSTATE")
		}
		// 解码得到隐藏的html内容
		hiddenContent, err := lzstring.DecompressFromBase64(val)
		if err != nil {
			return types.TreeNode{}, fmt.Errorf("decompress __VIEWSTATE failed: %w", err)
		}
		// 重新解析隐藏的html内容
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(hiddenContent))
		if err != nil {
			return types.TreeNode{}, fmt.Errorf("parse hidden html failed: %w", err)
		}
	}

	chapterTypes, err := getChapterTypes(doc)
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("get chapter types failed: %w", err)
	}

	comicInfo := ComicInfo{
		Title:        title,
		ChapterTypes: chapterTypes,
	}
	// 构建树
	root, err := buildTree(&comicInfo, cacheDir)
	if err != nil {
		return types.TreeNode{}, fmt.Errorf("build tree failed: %w", err)
	}

	return root, nil
}

func ComicByKeyword(keyword string) ([]ComicSearchResult, error) {
	resp, err := http_client.HttpClientInst().Get(fmt.Sprintf("https://www.manhuagui.com/s/%s.html", keyword))
	if err != nil {
		return []ComicSearchResult{}, fmt.Errorf("do request failed: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []ComicSearchResult{}, fmt.Errorf("read response body failed: %w", err)
	}

	htmlContent := string(respBody)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return []ComicSearchResult{}, fmt.Errorf("parse html failed: %w", err)
	}

	var results []ComicSearchResult
	doc.Find(".book-detail").Each(func(_ int, div *goquery.Selection) {
		var result ComicSearchResult
		// 获取书名和漫画ID
		a := div.Find("dt a").First()
		title, titleExists := a.Attr("title")
		if titleExists {
			result.Title = title
		}
		href, hrefExists := a.Attr("href")
		if hrefExists {
			parts := strings.Split(href, "/")
			result.ComicId = parts[2]
		}

		// 获取作者名
		div.Find("dd.tags span a").Each(func(_ int, s *goquery.Selection) {
			// 跳过非作者链接
			href, hrefExists := s.Attr("href")
			if !hrefExists || !strings.HasPrefix(href, "/author/") {
				return
			}

			author, authorExist := s.Attr("title")
			if authorExist {
				result.Authors = append(result.Authors, author)
			}
		})

		results = append(results, result)
	})
	return results, nil
}

func getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("h1").Text()
	return title, nil
}

func getChapterTypes(doc *goquery.Document) ([]ChapterType, error) {
	var chapterTypes []ChapterType

	doc.Find("h4").Each(func(i int, h4 *goquery.Selection) {
		chapterType := ChapterType{Title: h4.Find("span").Text()}

		// class中包含chapter-page的div表示这个章节类型有分页
		if h4.Next().Is("div[class~=chapter-page]") {
			chapterPageDiv := h4.Next()
			chapterPageDiv.Find("a").Each(func(_ int, a *goquery.Selection) {
				title, exist := a.Attr("title")
				if exist {
					chapterType.ChapterPages = append(chapterType.ChapterPages, ChapterPage{Title: title})
				}
			})

			chapterListDiv := chapterPageDiv.Next()
			chapterListDiv.Find("ul").Each(func(pageIndex int, ul *goquery.Selection) {
				// 每个ul表示一个分页
				chapterType.ChapterPages[pageIndex].Chapters = getChaptersFromUl(ul)
			})

		} else { // 这个章节类型没有分页
			chapterListDiv := h4.Next()
			ul := chapterListDiv.Find("ul").First()
			chapters := getChaptersFromUl(ul)
			page := ChapterPage{Chapters: chapters}
			chapterType.ChapterPages = []ChapterPage{page}
		}

		chapterTypes = append(chapterTypes, chapterType)
	})

	return chapterTypes, nil
}

func getChaptersFromUl(ul *goquery.Selection) []Chapter {
	var chapters []Chapter

	ul.Find("a").Each(func(_ int, a *goquery.Selection) {
		href, hrefExist := a.Attr("href")
		title, titleExist := a.Attr("title")
		if hrefExist && titleExist {
			chapter := Chapter{Title: title, Href: href}
			chapters = append(chapters, chapter)
		}
	})

	slices.Reverse(chapters)
	return chapters
}

func buildTree(comicInfo *ComicInfo, cacheDir string) (types.TreeNode, error) {
	root := types.TreeNode{
		Label:         comicInfo.Title,
		Key:           filepath.ToSlash(path.Join(cacheDir, comicInfo.Title)),
		Children:      []types.TreeNode{},
		DefaultExpand: true,
	}

	for _, chapterType := range comicInfo.ChapterTypes {
		chapterTypeNode := types.TreeNode{
			Label:         chapterType.Title,
			Key:           filepath.ToSlash(path.Join(root.Key, chapterType.Title)),
			Children:      []types.TreeNode{},
			DefaultExpand: true,
		}

		for _, chapterPage := range chapterType.ChapterPages {
			chapterPageNode := types.TreeNode{
				Label:    chapterPage.Title,
				Key:      filepath.ToSlash(path.Join(chapterTypeNode.Key, chapterPage.Title)),
				Children: []types.TreeNode{},
			}

			for _, chapter := range chapterPage.Chapters {
				saveDir := filepath.ToSlash(path.Join(chapterPageNode.Key, chapter.Title))
				saveDirExists := utils.PathExists(saveDir)
				keyJsonBytes, err := json.Marshal(ChapterTreeNodeKey{
					Href:    chapter.Href,
					SaveDir: saveDir,
				})
				if err != nil {
					return types.TreeNode{}, fmt.Errorf("marshal key failed: %w", err)
				}

				chapterNode := types.TreeNode{
					Label:          chapter.Title,
					Key:            string(keyJsonBytes),
					IsLeaf:         true,
					Disabled:       saveDirExists,
					Children:       []types.TreeNode{},
					DefaultChecked: saveDirExists,
				}
				chapterPageNode.Children = append(chapterPageNode.Children, chapterNode)
			}

			chapterTypeNode.Children = append(chapterTypeNode.Children, chapterPageNode)
		}

		// 如果只有一个分页，就不要显示分页了，直接显示章节
		if len(chapterTypeNode.Children) == 1 {
			page := chapterTypeNode.Children[0]
			chapterTypeNode.Children = page.Children
		}

		root.Children = append(root.Children, chapterTypeNode)
	}

	return root, nil
}
