package common

//
//import "strings"
//
//type Node struct {
//	Name      string
//	Childes   []*Node
//	Connector *Connector
//}
//
//func NewNode(name string) *Node {
//	return &Node{
//		Name:      name,
//		Childes:   nil,
//		Connector: nil,
//	}
//}
//
//func (n *Node) Add(name string, connector *Connector) *Node {
//	kinds := strings.Split(name, ".")
//	if len(kinds) == 1 {
//		n.Childes = append(n.Childes, &Node{
//			Name:      name,
//			Childes:   nil,
//			Connector: connector,
//		})
//		return n
//	}
//	if len(kinds) == 2 {
//		for _, child := range n.Childes {
//			if child.Name == kinds[0] {
//				return n.Add(kinds[1], connector)
//			}
//		}
//
//	}
//
//}
