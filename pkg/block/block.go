package block

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/lestoni/sapphire/pkg/node"
)

type Container interface {
	AddNode(node *node.Node) error
	AddNodes(nodes []*node.Node) error
	Verify(hash string) bool
	GetNode(identity string) *node.Node
}

type Block struct {
	Nodes    []*node.Node
	Previous string
	Size     int
	Root     *node.Node
	Weight   int
	Height   int
	lock     sync.Mutex
	Identity string
	MRoot    string
}

const (
	blockSize = 1024000000
)

func NewRoot() (*Block, error) {
	identity, err := calculateIdentity()
	if err != nil {
		return nil, err
	}

	block := &Block{
		Identity: identity,
		Size:     blockSize,
	}

	return block, nil
}

func New(previous string) (*Block, error) {
	if len(previous) == 0 {
		return nil, fmt.Errorf("Previous Reference Block Not Set")
	}

	identity, err := calculateIdentity()
	if err != nil {
		return nil, err
	}

	block := &Block{
		Identity: identity,
		Size:     blockSize,
		Previous: previous,
	}

	return block, nil
}

func (b *Block) AddNode(item *node.Node) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	//fmt.Println("Add Node--", time.Now().String())
	//time.Sleep(time.Second)

	if err := b.validateNodeWeight(item); err != nil {
		return err
	}

	// Validate Block Size
	if err := b.validateWeight(); err != nil {
		return err
	}

	// Compute Identity
	if len(b.Nodes) == 0 {
		b.Root = item
	}

	b.Nodes = append(b.Nodes, item)
	b.updateDimensions()

	return nil
}

func (b *Block) AddNodes(items []*node.Node) error {
	// Validate Block Size

	if err := b.validateWeight(); err != nil {
		return err
	}

	for _, item := range items {
		err := b.AddNode(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Block) Build() error {
	if len(b.Nodes) == 0 {
		return fmt.Errorf("Block is Null")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	fmt.Println("Building")

	// if odd number of leaves, prefill with emptyNode
	// to achieve a balanced binary tree
	if len(b.Nodes)%2 == 1 {
		prefiller := node.NewRoot()
		prefiller.AddContent("")

		b.Nodes = append(b.Nodes, prefiller)
	}

	var leaves []string

	for i := 0; i < len(b.Nodes); i += 2 {
		left := b.Nodes[i]
		right := b.Nodes[i+1]

		combined := getHash(append(left.Content[:], right.Content[:]...))

		leaves = append(leaves, combined)
	}

	b.MRoot = computeMerkleRoot(leaves)

	return nil
}

func (b *Block) validateNodeWeight(node *node.Node) error {
	if len(node.Content) > b.Size {
		return fmt.Errorf("Block Size Exceeded by Node")
	}

	return nil
}

func (b *Block) validateWeight() error {
	if b.Weight >= b.Size {
		return fmt.Errorf("Block Size Exceeded")
	}

	return nil
}

func (b *Block) Verify(hash string) bool {
	return b.MRoot == hash
}

// Walk Tree
func (b *Block) GetNode(identity string) *node.Node {
	return b.Nodes[0]
}

func (b *Block) updateDimensions() {
	for _, node := range b.Nodes {

		b.Weight += len(node.Content)
	}

	b.Height = len(b.Nodes)
}

func computeMerkleRoot(leaves []string) string {
	if len(leaves)%2 == 1 {
		leaves = append(leaves, "")
	}

	var level []string

	for i := 0; i < len(leaves); i += 2 {
		left := leaves[i]
		right := leaves[i+1]

		combined := getHash([]byte(left + right))

		level = append(level, combined)
	}

	if len(level) == 1 {
		return level[0]
	} else {
		return computeMerkleRoot(level)
	}
}

func getHash(data []byte) string {
	firstIteration := sha256.New()
	secondIteration := sha256.New()

	firstIteration.Write(data)
	secondIteration.Write(firstIteration.Sum(nil))

	return fmt.Sprintf("%x", secondIteration.Sum(nil))
}

func calculateIdentity() (string, error) {
	var identity string

	// Identity: current Timestamp, machine id, random salt,
	nBig, err := rand.Int(rand.Reader, big.NewInt(big.MaxExp))
	if err != nil {
		return identity, err
	}

	hash := sha256.New()
	uuid := time.Now().String() + strconv.FormatInt(nBig.Int64(), 16)
	hash.Write([]byte(uuid))

	identity = fmt.Sprintf("%x", hash.Sum(nil))

	return identity, nil
}
