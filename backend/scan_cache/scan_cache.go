package scan_cache

import (
	"fmt"
	"manhuagui-downloader/backend/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type TreeNode struct {
	Label    string     `json:"label"`
	Key      string     `json:"key"`
	Children []TreeNode `json:"children"`
	IsLeaf   bool       `json:"isLeaf"`
}

func ScanCacheDir(cacheDir string, maxDepth int64) ([]TreeNode, error) {
	// 将路径中的反斜杠转换为正斜杠
	cacheDir = filepath.ToSlash(cacheDir)

	root := TreeNode{
		Label:    path.Base(cacheDir),
		Key:      cacheDir,
		Children: []TreeNode{},
	}

	err := buildTree(&root, 0, maxDepth)
	if err != nil {
		return []TreeNode{}, fmt.Errorf("build tree failed: %w", err)
	}

	return root.Children, nil
}

func buildTree(node *TreeNode, depth int64, maxDeep int64) error {
	defer func() { node.IsLeaf = isLeafNode(node) }()

	if depth > maxDeep {
		return nil
	}

	entries, err := os.ReadDir(node.Key)
	if err != nil {
		return fmt.Errorf("read dir failed: %w", err)
	}

	//给 entries 按照更合理的文件名排序
	sort.Slice(entries, func(i, j int) bool { return utils.FilenameComparer(entries[i].Name(), entries[j].Name()) })

	for _, entry := range entries {
		// 忽略非目录和隐藏文件
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := path.Join(node.Key, entry.Name())
		childNode := TreeNode{
			Label:    entry.Name(),
			Key:      childPath,
			Children: []TreeNode{},
		}
		err = buildTree(&childNode, depth+1, maxDeep)
		if err != nil {
			return fmt.Errorf("build tree failed: %w", err)
		}

		node.Children = append(node.Children, childNode)
	}

	return nil
}

func isLeafNode(node *TreeNode) bool {
	// 如果无法读取目录，则认为是叶子节点
	entries, err := os.ReadDir(node.Key)
	if err != nil {
		return true
	}

	// 如果有子目录，则不是叶子节点
	for _, entry := range entries {
		if entry.IsDir() {
			return false
		}
	}

	return true
}
