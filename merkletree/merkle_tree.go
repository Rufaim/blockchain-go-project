package merkletree

import (
	"crypto/sha256"
	"math"
)

func hashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

//MerkleTree represents a root node of a tree
type MerkleTree *MerkleNode

//MerkleTreeHash returns a root hash after building a tree
func MerkleTreeHash(data [][]byte) []byte {
	return NewMerkleTree(data).Data
}

// MerkleNode represents a Merkle tree node
// where the latter is a binary tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode) *MerkleNode {
	data := append(left.Data, right.Data...)

	return &MerkleNode{
		Left:  left,
		Right: right,
		Data:  hashData(data),
	}
}

func NewMerkleTree(data [][]byte) MerkleTree {
	nodes := make([]*MerkleNode, 0, len(data))

	for _, d := range data {
		nodes = append(nodes, &MerkleNode{Data: hashData(d)})
	}

	for i := 0; i < int(math.Ceil(math.Log2(float64(len(data))))); i++ {
		var newNodes []*MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1])
			newNodes = append(newNodes, node)
		}
		nodes = newNodes
	}
	return MerkleTree(nodes[0])
}
