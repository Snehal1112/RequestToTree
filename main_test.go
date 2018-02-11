package treema

import (
	"bytes"
	"testing"
)

var req = `
{
  "provider": "auth" ,
  "resource": "user_authentication" ,
  "request_type": "item" ,
  "action": "read" ,
  "message_action": {
    "action_type": "source" ,
    "destination_action_type": "source"
  } ,
  "items": [
    {
      "username": "sd" ,
      "password": "aa" ,
      "token": ""
    }
  ]
}`

func loadTree() *Tree {
	b := bytes.NewBufferString(req)
	tree := NewTree()
	tree.Load(b.Bytes())
	return tree
}

func TestTree_GetNodeByName(t *testing.T) {
	tree := loadTree()
	n := tree.GetNodeByName("action")
	if n == nil || n == "action" {
		t.Errorf("Node not found")
	}
}
func TestTree_GetRootNode(t *testing.T) {
	tree := loadTree()
	rootNode := tree.GetRootNode()

	if rootNode.name != RootName && !rootNode.isRootNode() {
		t.Errorf("Node not found Root node")
	}
}

func TestNode_CreateNode(t *testing.T) {
	tree := loadTree()
	n := tree.GetNodeByName("message_action").(*Node)
	if n.hasChild() == false {
		t.Fatal("problem in creating child nodes")
	}
	if len(n.children) < 2 {
		t.Errorf("problem in creating child nodes")
	}
}

func TestNode_GetParentNode(t *testing.T) {
	tree := loadTree()
	msgActionNode := tree.GetNodeByName("message_action").(*Node)
	msgActionTypeNode := tree.GetNodeByName("action_type").(*Node)
	if msgActionTypeNode.GetParentNode().name != msgActionNode.name {
		t.Errorf("node has wrong parent node")
	}
}

func TestCreate_EmptyNode(t *testing.T) {
	node := NewEmptyNode(false)
	if len(node.name) != 0 {
		t.Errorf("node is not created properly.")
	}
}

func TestCreate_RootNode(t *testing.T) {
	node := NewEmptyNode(true)
	if node.name != RootName && !node.isRootNode() {
		t.Errorf("Root node is not created properly.")
	}
}
