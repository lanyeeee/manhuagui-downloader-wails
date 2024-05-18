package types

type TreeNode struct {
	Label          string     `json:"label"`
	Key            string     `json:"key"`
	Children       []TreeNode `json:"children"`
	IsLeaf         bool       `json:"isLeaf"`
	Disabled       bool       `json:"disabled"`
	DefaultChecked bool       `json:"defaultChecked"`
	DefaultExpand  bool       `json:"defaultExpand"`
}
