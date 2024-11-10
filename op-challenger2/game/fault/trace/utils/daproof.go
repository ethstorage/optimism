package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left, Right *MerkleNode
	Hash        []byte
}

// MerkleTree represents a non-2^n Merkle tree
type MerkleTree struct {
	Root   *MerkleNode
	Leaves []*MerkleNode
}

// NewMerkleNode creates a new Merkle node, hashing the left and right children if both are present
func NewMerkleNode(left, right *MerkleNode, data common.Hash) *MerkleNode {
	node := &MerkleNode{}
	if left == nil && right == nil {
		// Leaf node
		node.Hash = data[:]
	} else {
		// Internal node
		var combined []byte
		if right != nil {
			combined = append(left.Hash[:], right.Hash[:]...)
		} else {
			combined = left.Hash[:] // Only one child, carry the hash up
		}
		node.Hash = crypto.Keccak256(combined)
	}
	node.Left = left
	node.Right = right
	return node
}

// NewMerkleTree constructs a Merkle tree from given leaf data without padding
func NewMerkleTree(data []common.Hash) *MerkleTree {
	var leaves []*MerkleNode
	for _, datum := range data {
		leaves = append(leaves, NewMerkleNode(nil, nil, datum))
	}
	tree := &MerkleTree{Leaves: leaves}
	tree.Root = buildTree(leaves)
	return tree
}

// buildTree recursively builds the Merkle tree without padding
func buildTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var nextLevel []*MerkleNode
	for i := 0; i < len(nodes); i += 2 {
		if i+1 < len(nodes) {
			// Pair of nodes
			parent := NewMerkleNode(nodes[i], nodes[i+1], common.Hash{})
			nextLevel = append(nextLevel, parent)
		} else {
			// Unpaired node, carry it up to the next level
			nextLevel = append(nextLevel, nodes[i])
		}
	}
	return buildTree(nextLevel)
}

// GenerateProof generates a Merkle proof for a leaf node at a given index
func (mt *MerkleTree) GenerateProof(index uint32) []byte {
	var proof []byte
	if index < 0 || index >= uint32(len(mt.Leaves)) {
		return proof // Invalid index
	}

	node := mt.Leaves[index]
	for node != mt.Root {
		parent := findParent(mt.Root, node)
		sibling := getSibling(parent, node)
		if sibling != nil {
			proof = append(proof, sibling.Hash...)
		}
		node = parent
	}
	return proof
}

// findParent finds the parent node of a given child in the Merkle tree
func findParent(root, child *MerkleNode) *MerkleNode {
	if root == nil || (root.Left == child || root.Right == child) {
		return root
	}
	if node := findParent(root.Left, child); node != nil {
		return node
	}
	return findParent(root.Right, child)
}

// getSibling returns the sibling of a given node in the tree
func getSibling(parent, node *MerkleNode) *MerkleNode {
	if parent == nil || (parent.Left != node && parent.Right != node) {
		return nil
	}
	if parent.Left == node {
		return parent.Right
	}
	return parent.Left
}

func GenerateProofForSubValues(subValues []common.Hash, index uint32) []byte {
	tree := NewMerkleTree(subValues)
	return tree.GenerateProof(index)
}
