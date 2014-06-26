package upnp

import (
	"bytes"
)

type Node struct {
	Name    string
	Content string
	Attr    map[string]string
	Child   []Node
}

func (n *Node) AddChild(node Node) {
	n.Child = append(n.Child, node)
}
func (n *Node) BuildXML() string {
	buf := bytes.NewBufferString("<")
	buf.WriteString(n.Name)
	for key, value := range n.Attr {
		buf.WriteString(" ")
		buf.WriteString(key + "=" + value)
	}
	buf.WriteString(">" + n.Content)

	for _, node := range n.Child {
		buf.WriteString(node.BuildXML())
	}
	buf.WriteString("</" + n.Name + ">")
	return buf.String()
}
