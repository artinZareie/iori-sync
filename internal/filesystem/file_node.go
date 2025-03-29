package filesystem

import "encoding/json"

// This structure is meant to be used as a tree node.
type FileNode struct {
	File     File
	Children []*FileNode
}

func (fn *FileNode) AddChild(child *FileNode) {
	fn.Children = append(fn.Children, child)
}

func (fn *FileNode) ToFileJSON() FileJSON {
	return fn.File.ToFileJSON()
}

func (fn *FileNode) ToJSON() string {
	// Create a nested json structure from the tree rooted at fn.
	type jsonNode struct {
		File     FileJSON   `json:"file"`
		Children []jsonNode `json:"children,omitempty"`
	}

	var buildNode func(node *FileNode) jsonNode
	buildNode = func(node *FileNode) jsonNode {
		result := jsonNode{
			File: node.ToFileJSON(),
		}
		
		if len(node.Children) > 0 {
			result.Children = make([]jsonNode, 0, len(node.Children))
			for _, child := range node.Children {
				result.Children = append(result.Children, buildNode(child))
			}
		}
		
		return result
	}
	
	root := buildNode(fn)
	jsonData, err := json.Marshal(root)
	if err != nil {
		return "{}"
	}
	
	return string(jsonData)
}
