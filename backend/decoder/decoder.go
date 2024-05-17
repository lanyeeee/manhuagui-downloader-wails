package decoder

import (
	"encoding/json"
	"fmt"
	lzstring "github.com/daku10/go-lz-string"
	"regexp"
	"strconv"
	"strings"
)

type DecodeResult struct {
	Bid      int      `json:"bid"`
	Bname    string   `json:"bname"`
	Bpic     string   `json:"bpic"`
	Cid      int      `json:"cid"`
	Cname    string   `json:"cname"`
	Files    []string `json:"files"`
	Finished bool     `json:"finished"`
	Len      int      `json:"len"`
	Path     string   `json:"path"`
	Status   int      `json:"status"`
	BlockCc  string   `json:"block_cc"`
	NextId   int      `json:"nextId"`
	PrevId   int      `json:"prevId"`
	Sl       struct {
		E int    `json:"e"`
		M string `json:"m"`
	} `json:"sl"`
}

func Decode(htmlContent *string) (DecodeResult, error) {
	function, a, c, data, err := decodeHtmlContent(htmlContent)
	if err != nil {
		return DecodeResult{}, fmt.Errorf("decode html content failed: %w", err)
	}

	dict := generateDict(a, c, data)

	js := generateJs(function, dict)

	return generateDecodeResult(js)
}

func decodeHtmlContent(htmlContent *string) (string, int, int, []string, error) {
	re := regexp.MustCompile(`^.*}\('(.*)',(\d*),(\d*),'([\w|+/=]*)'.*$`)
	matches := re.FindStringSubmatch(*htmlContent)

	function := matches[1]

	a, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("convert a to int failed: %w", err)
	}

	c, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("convert c to int failed: %w", err)
	}

	decompress, err := lzstring.DecompressFromBase64(matches[4])
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("decompress data failed: %w", err)
	}
	data := strings.Split(decompress, "|")

	return function, a, c, data, nil
}

func generateJs(function string, dict map[string]string) string {
	re := regexp.MustCompile(`(\b\w+\b)`)
	splits := re.Split(function, -1)
	matches := re.FindAllString(function, -1)
	var pieces []string
	for i := 0; i < len(splits); i++ {
		if i < len(matches) {
			pieces = append(pieces, splits[i], matches[i])
		} else {
			pieces = append(pieces, splits[i])
		}
	}
	js := ""
	for _, x := range pieces {
		val, ok := dict[x]
		if ok {
			js += val
		} else {
			js += x
		}
	}
	return js
}

func generateDict(a int, c int, data []string) map[string]string {
	var itr func(value int, num int) string
	itr = func(value int, num int) string {
		const d = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		if value <= 0 {
			return ""
		}
		return itr(value/num, num) + string(d[value%a])

	}

	tr := func(value int, num int, a int) string {
		tmp := itr(value, num)
		if tmp == "" {
			return "0"
		}
		return tmp
	}

	var e func(c int) string
	e = func(c int) string {
		return func() string {
			if c < a {
				return ""
			}
			return e(c / a)
		}() + func() string {
			if c%a > 35 {
				return string(rune(c%a + 29))
			}
			return tr(c%a, 36, a)
		}()
	}

	dict := make(map[string]string)
	for c -= 1; c+1 > 0; c-- {
		if data[c] == "" {
			dict[e(c)] = e(c)
		} else {
			dict[e(c)] = data[c]
		}
	}

	return dict
}

func generateDecodeResult(js string) (DecodeResult, error) {
	re := regexp.MustCompile(`^.*\((\{.*})\).*$`)
	matches := re.FindStringSubmatch(js)

	var result DecodeResult
	err := json.Unmarshal([]byte(matches[1]), &result)
	if err != nil {
		return DecodeResult{}, fmt.Errorf("unmarshal decode result failed: %w", err)
	}

	return result, nil
}
