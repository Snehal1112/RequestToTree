package treema

import (
	"encoding/json"
	"log"
	"reflect"
)

func NewTree() *Tree {
	return &Tree{}
}

func NewEmptyNode(root bool) *Node {
	if root {
		return &Node{
			name:   RootName,
			isRoot: root,
		}
	}
	return &Node{}
}

const RootName = "root"

type Tree struct {
	node   *Node
	source string
}

func (t *Tree) GetNodeByName(name string) interface{} {
	result := t.getNodeByName(name, t.GetRootNode())
	return result
}

func (t *Tree) getNodeByName(name string, node *Node) interface{} {
	if node.name == name {
		return node
	}
	for _, e := range node.children {
		if result := t.getNodeByName(name, e); result != nil {
			return result
		}
	}
	return nil
}

func (t *Tree) GetRootNode() *Node {
	return t.node
}

type Node struct {
	name        string       `json:"name"`
	leaf        bool         `json:"leaf"`
	value       interface{}  `json:"value"`
	nodeType    reflect.Type `json:"node_type"`
	children    []*Node      `json:"children"`
	parent      *Node        `json:"parent"`
	hasChildren bool         `json:"has_children"`
	isRoot      bool         `json:"is_root"`
}

func (t *Tree) Load(raw []byte) *Tree {
	mapData := make(map[string]interface{})
	if err := json.Unmarshal(raw, &mapData); err != nil {
		log.Println(err.Error())
	}
	nodes := make([]*Node, 0)
	rootNode := NewEmptyNode(true)
	for k, v := range mapData {
		node := NewEmptyNode(false)
		node.createNode(k, reflect.TypeOf(v), v, rootNode)
		nodes = append(nodes, node)
	}
	rootNode.children = nodes
	rootNode.hasChildren = true
	t.source = string(raw)
	t.node = rootNode
	return t
}

func (n *Node) isRootNode() bool {
	return n.isRoot
}

func (n *Node) hasChild() bool {
	return n.hasChildren
}

func (n *Node) GetParentNode() *Node {
	return n.parent
}

func (n *Node) createNode(key string, t reflect.Type, value interface{}, parent *Node) *Node {
	n.name = string(key)
	switch t.Kind().String() {
	case "map", "slice":
		n.hasChildren = true
		n.leaf = false
		break
	case "string", "bool", "int", "float":
	default:
		n.hasChildren = false
		n.leaf = true
	}
	n.value = value
	n.nodeType = t
	n.parent = parent

	if n.hasChildren {
		tyy := reflect.TypeOf(value).Kind().String()
		if tyy == "slice" {
			for _, ee := range value.([]interface{}) {
				ch := ee.(map[string]interface{})
				for ty, e := range ch {
					node := NewEmptyNode(false)
					cn := node.createNode(ty, reflect.TypeOf(e), e, n)
					n.children = append(n.children, cn)
				}
			}
			return n
		}
		ch := value.(map[string]interface{})
		for ty, e := range ch {
			node := NewEmptyNode(false)
			cn := node.createNode(ty, reflect.TypeOf(e), e, n)
			n.children = append(n.children, cn)
		}
	}
	return n
}
