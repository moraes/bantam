// Copyright 2013 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bantam

import (
	"bytes"
	"fmt"
)

// Node is the basic interface for expression nodes.
type Node interface {
	String() string
}

// ----------------------------------------------------------------------------

// AssignNode represents an assignment expression like "a = b".
type AssignNode struct {
	Name  string
	Right Node
}

func NewAssignNode(name string, right Node) *AssignNode {
	return &AssignNode{Name: name, Right: right}
}

func (n *AssignNode) String() string {
	return fmt.Sprintf("(%s = %s)", n.Name, n.Right)
}

// ----------------------------------------------------------------------------

// BinaryNode represents a binary arithmetic expression like "a + b".
type BinaryNode struct {
	Left     Node
	Operator TokenType
	Right    Node
}

func NewBinaryNode(left Node, operator TokenType, right Node) *BinaryNode {
	return &BinaryNode{Left: left, Operator: operator, Right: right}
}

func (n *BinaryNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left, n.Operator, n.Right)
}

// ----------------------------------------------------------------------------

// FunctionNode represents a function call like "a(b, c, d)".
type FunctionNode struct {
	Function Node
	Args     *ListNode
}

func NewFunctionNode(function Node, args *ListNode) *FunctionNode {
	return &FunctionNode{Function: function, Args: args}
}

func (n *FunctionNode) String() string {
	b := new(bytes.Buffer)
	for k, v := range n.Args.Nodes {
		fmt.Fprint(b, v)
		if k < len(n.Args.Nodes)-1 {
			b.WriteString(", ")
		}
	}
	return fmt.Sprintf("%s(%s)", n.Function, b)
}

// ----------------------------------------------------------------------------

// ListNode holds a sequence of nodes.
type ListNode struct {
	Nodes []Node // The element nodes in lexical order.
}

func NewListNode() *ListNode {
	return &ListNode{}
}

func (n *ListNode) Append(node Node) {
	n.Nodes = append(n.Nodes, node)
}

func (n *ListNode) String() string {
	b := new(bytes.Buffer)
	for _, v := range n.Nodes {
		fmt.Fprint(b, v)
	}
	return b.String()
}

func listNode(n Node) *ListNode {
	list, ok := n.(*ListNode)
	if !ok {
		list = &ListNode{}
		list.Append(n)
	}
	return list
}

// ----------------------------------------------------------------------------

// NameNode represents a simple variable name expression like "abc".
type NameNode struct {
	Name string
}

func NewNameNode(name string) *NameNode {
	return &NameNode{Name: name}
}

func (n *NameNode) String() string {
	return n.Name
}

// ----------------------------------------------------------------------------

// TernaryNode represents a ternary expression like "a ? b : c".
type TernaryNode struct {
	Condition Node
	List      *ListNode
	ElseList  *ListNode
}

func NewTernaryNode(condition Node, list, elseList *ListNode) *TernaryNode {
	return &TernaryNode{Condition: condition, List: list, ElseList: elseList}
}

func (n *TernaryNode) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", n.Condition, n.List, n.ElseList)
}

// ----------------------------------------------------------------------------

// UnaryNode represents a prefix unary arithmetic expression like "!a" or "-b".
type UnaryNode struct {
	Operator TokenType
	Right    Node
}

func NewUnaryNode(operator TokenType, right Node) *UnaryNode {
	return &UnaryNode{Operator: operator, Right: right}
}

func (n *UnaryNode) String() string {
	return fmt.Sprintf("(%s%s)", n.Operator, n.Right)
}

// ----------------------------------------------------------------------------

// UnaryPostfixNode represents a postfix unary arithmetic expression like "a++".
type UnaryPostfixNode struct {
	Left     Node
	Operator TokenType
}

func NewUnaryPostfixNode(left Node, operator TokenType) *UnaryPostfixNode {
	return &UnaryPostfixNode{Left: left, Operator: operator}
}

func (n *UnaryPostfixNode) String() string {
	return fmt.Sprintf("(%s%s)", n.Left, n.Operator)
}
