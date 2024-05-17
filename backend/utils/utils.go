package utils

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Sanitize 去掉非法字符，保证目录或文件名合法
func Sanitize(dirOrFileName string) string {
	illegalChars := `<>:"/\|?*`
	var stringBuilder strings.Builder

	for _, r := range dirOrFileName {
		// 如果字符是非法字符，跳过
		if !strings.ContainsRune(illegalChars, r) && !unicode.IsControl(r) {
			stringBuilder.WriteRune(r)
		}
	}

	// 去掉开头和结尾的空格
	return strings.TrimSpace(stringBuilder.String())
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FilenameComparer(a, b string) bool {
	// 分割数字和非数字
	splitRegexp := regexp.MustCompile(`(\d+\.\d+|\d+|\D+)`)

	aMatches := splitRegexp.FindAllString(a, -1)
	bMatches := splitRegexp.FindAllString(b, -1)

	// 比较每个分割后的部分
	for i := 0; i < len(aMatches) && i < len(bMatches); i++ {
		aMatch, bMatch := aMatches[i], bMatches[i]

		aNumber, aErr := strconv.ParseFloat(aMatch, 64)
		bNumber, bErr := strconv.ParseFloat(bMatch, 64)
		// 如果两部分都是数字，则按数字大小进行比较
		if aErr == nil && bErr == nil {
			if aNumber != bNumber {
				return aNumber < bNumber
			}
			// 如果数字相等，继续比较
			continue
		}

		// 如果一个是数字而另一个不是，则数字视为较小
		if aErr == nil {
			return true
		}
		if bErr == nil {
			return false
		}

		// 如果都不是数字，则按字符串比较
		if aMatch != bMatch {
			return aMatch < bMatch
		}
	}

	// 如果所有匹配的部分都相等，那么比较它们的长度
	return len(aMatches) < len(bMatches)
}
