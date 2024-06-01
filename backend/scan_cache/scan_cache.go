package scan_cache

import (
	"fmt"
	"manhuagui-downloader/backend/types"
	"manhuagui-downloader/backend/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func ScanCacheDir(cacheDir string, exportDir string, maxDepth int64) ([]types.TreeNode, error) {
	// 将路径中的反斜杠转换为正斜杠
	cacheDir = filepath.ToSlash(cacheDir)

	root := types.TreeNode{
		Key:      cacheDir,
		Children: []types.TreeNode{},
	}
	if err := buildTree(&root, cacheDir, exportDir, 0, maxDepth); err != nil {
		return []types.TreeNode{}, fmt.Errorf("build tree failed: %w", err)
	}

	return root.Children, nil
}

func buildTree(node *types.TreeNode, cacheDir string, exportDir string, depth int64, maxDeep int64) error {
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

		key := path.Join(node.Key, entry.Name())
		isLeaf := isLeafNode(key)
		exported := isExported(key, cacheDir, exportDir)
		disabled := false
		if isLeaf { // 只有叶子节点才能disable
			disabled = exported
		}
		defaultExpand := false
		if depth < 1 { // 默认展开第一层
			defaultExpand = true
		}

		childNode := types.TreeNode{
			Label:          entry.Name(),
			Key:            key,
			Children:       []types.TreeNode{},
			DefaultExpand:  defaultExpand,
			IsLeaf:         isLeaf,
			DefaultChecked: exported,
			Disabled:       disabled,
		}
		//fmt.Printf("childNode: %v\n", childNode)
		if err = buildTree(&childNode, cacheDir, exportDir, depth+1, maxDeep); err != nil {
			return fmt.Errorf("build tree failed: %w", err)
		}

		node.Children = append(node.Children, childNode)
	}

	return nil
}

func isLeafNode(key string) bool {
	// 如果无法读取目录，则认为是叶子节点
	entries, err := os.ReadDir(key)
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

func isExported(key string, cacheDir string, exportDir string) bool {
	relPath, err := filepath.Rel(cacheDir, key)
	if err != nil {
		return false
	}
	pdfPath := path.Join(exportDir, relPath+".pdf")

	return utils.PathExists(pdfPath)
}
