package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Node struct {
	Content  []byte
	Parent   string
	Identity string
}

// Implementation of merkel tree

func New(parent string) *Node {
	return &Node{Parent: parent}
}

func NewRoot() *Node {
	node := &Node{}

	return node
}

func (n *Node) AddContent(content interface{}) error {
	encoded, err := getBytes(content)
	if err != nil {
		return err
	}

	n.Content = encoded

	n.computeIdentity()

	return nil
}

func (n *Node) computeIdentity() {
	hash := sha256.New()
	hash.Write(n.Content)

	n.Identity = fmt.Sprintf("%x", hash.Sum(nil))
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
