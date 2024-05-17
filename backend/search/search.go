package search

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"manhuagui-downloader/backend/http_client"
	"net/http"
	"slices"
	"strings"
)

type ComicInfo struct {
	Title        string        `json:"title"`
	ChapterTypes []ChapterType `json:"chapterTypes"`
}

type ChapterType struct {
	Title         string         `json:"title"`
	ChapterPagers []ChapterPager `json:"chapterPagers"`
}

type ChapterPager struct {
	Title    string    `json:"title"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func Info(comicId string) (ComicInfo, error) {
	resp, err := http_client.HttpClientInst().Get("https://www.manhuagui.com/comic/" + comicId)
	if err != nil {
		return ComicInfo{}, fmt.Errorf("do request failed: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ComicInfo{}, fmt.Errorf("read response body failed: %w", err)
	}

	// 处理HTTP错误
	switch resp.StatusCode {
	case http.StatusOK:
		// ignore
	case http.StatusNotFound:
		return ComicInfo{}, errors.New(fmt.Sprintf("can't find comic with id: %s", comicId))
	default:
		return ComicInfo{}, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	htmlContent := string(respBody)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return ComicInfo{}, fmt.Errorf("parse html failed: %w", err)
	}

	title, err := getTitle(doc)
	if err != nil {
		return ComicInfo{}, fmt.Errorf("get title failed: %w", err)
	}

	chapterTypes, err := getChapterTypes(doc)
	if err != nil {
		return ComicInfo{}, fmt.Errorf("get chapter types failed: %w", err)
	}

	comicInfo := ComicInfo{
		Title:        title,
		ChapterTypes: chapterTypes,
	}

	return comicInfo, nil
}

func getChapterTypes(doc *goquery.Document) ([]ChapterType, error) {
	// chapterDiv是包含章节信息的div
	chapterDiv := doc.Find("div[class~=chapter]")
	var chapterTypes []ChapterType

	completeChapterTypeTitles(&chapterTypes, chapterDiv)
	completeChapterTypeChapterPagers(&chapterTypes, chapterDiv)

	return chapterTypes, nil
}

func completeChapterTypeTitles(chapterTypes *[]ChapterType, chapterDiv *goquery.Selection) {
	chapterDiv.Find("h4").Each(func(i int, h4 *goquery.Selection) {
		chapterTypeTitle := h4.Find("span").Text()
		if chapterTypeTitle != "" {
			*chapterTypes = append(*chapterTypes, ChapterType{Title: chapterTypeTitle})
		}
	})
}

func completeChapterTypeChapterPagers(chapterTypes *[]ChapterType, chapterDiv *goquery.Selection) {
	completeChapterPagerTitles(chapterTypes, chapterDiv)
	completeChapterPagerChapters(chapterTypes, chapterDiv)
}

func completeChapterPagerTitles(chapterTypes *[]ChapterType, chapterDiv *goquery.Selection) {
	chapterDiv.Find("div[class~=chapter-page]").Each(func(chapterTypeIndex int, div *goquery.Selection) {
		// 这个div的内容是该章节类型的分页信息，div中含有多个a标签，每个a标签对应一个分页
		div.Find("a").Each(func(_ int, a *goquery.Selection) {
			chapterPagerTitle, exist := a.Attr("title")
			if exist {
				chapterType := &(*chapterTypes)[chapterTypeIndex]
				chapterPager := ChapterPager{Title: chapterPagerTitle}
				chapterType.ChapterPagers = append(chapterType.ChapterPagers, chapterPager)
			}
		})
	})
	// 如果章节类型没有分页信息，则添加一个默认的分页
	for i := range *chapterTypes {
		chapterType := &(*chapterTypes)[i]
		if chapterType.ChapterPagers == nil {
			chapterType.ChapterPagers = []ChapterPager{{Title: "全部"}}
		}
	}
}

func completeChapterPagerChapters(chapterTypes *[]ChapterType, chapterDiv *goquery.Selection) {
	chapterDiv.Find("div[class~=chapter-list]").Each(func(chapterTypeIndex int, div *goquery.Selection) {
		chapterType := &(*chapterTypes)[chapterTypeIndex]
		// 这个div的内容是该张杰类型的章节信息，div中含有多个ul，每个ul对应一个分页
		div.Find("ul").Each(func(chapterPagerIndex int, ul *goquery.Selection) {
			chapterPager := &chapterType.ChapterPagers[chapterPagerIndex]
			// 每个ul中含有多个a标签，每个a标签对应一个章节
			ul.Find("a").Each(func(_ int, a *goquery.Selection) {
				href, exist := a.Attr("href")
				if !exist {
					return
				}

				title, exist := a.Attr("title")
				if !exist {
					return
				}

				chapter := Chapter{Title: title, Href: href}
				chapterPager.Chapters = append(chapterPager.Chapters, chapter)
			})
			slices.Reverse(chapterPager.Chapters)
		})
	})
}

func getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("h1").Text()
	return title, nil
}
