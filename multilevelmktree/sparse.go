/*
Package implements a sparse Merkle tree data structure.

The SparseMerkleTree struct represents a sparse Merkle tree and contains the
root node, depth, and a map of leaves. The MerklePathItem struct represents an
item in the Merkle tree path.

The package provides the following functions and methods:

Functions:
- getHashEmptyForDepth(depth int) *big.Int: Calculates the hash value for an
  empty node at a given depth.
- getPaddedBinaryString(i int, depth int) string: Returns a binary string
  representation of an integer, padded with leading zeros to a specified length.
- NewDeterministicSparseMerkleTree(depth int) *SparseMerkleTree: Creates a new
  deterministic sparse Merkle tree with non-null leaves.

Methods:
- NewSparseMerkleTree(depth int) *SparseMerkleTree: Creates a new sparse Merkle
  tree with empty leaves.
- (smt *SparseMerkleTree) Insert(key string, value *big.Int): Inserts a leaf
  with the given key and value into the tree.
- (smt *SparseMerkleTree) GenerateMerklePath(key string) ([]*MerklePathItem,
  error): Generates a Merkle tree path for the leaf with the given key.
- VerifyMerklePath(leafHash *big.Int, path []*MerklePathItem, expectedRoot
  *big.Int) bool: Verifies a Merkle tree path against the expected root hash.

Methods (internal):
- (node *MerkleNode) getLeftChild(depth int) *MerkleNode: Returns the left child
  node of the current node.
- (node *MerkleNode) getRightChild(depth int) *MerkleNode: Returns the right
  child node of the current node.
- hashChildren(left, right *MerkleNode, depth int) *big.Int: Computes the hash
  value of two child nodes.
- getPathBit(key string, depth int) int: Retrieves the bit value of the key at
  the specified depth.

Note: The code provided here has been enhanced with comments explaining the
purpose of each function and method.
*/

package multilevelmktree

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

// SparseMerkleTree represents a sparse Merkle tree.
type SparseMerkleTree struct {
	Root   *MerkleNode
	Depth  int
	Leaves map[string]*big.Int
}

// MerklePathItem represents an item in the Merkle tree path.
type MerklePathItem struct {
	SiblingHash *big.Int
	IsRight     bool
}

var zeroLeaf, _ = poseidon.Hash([]*big.Int{big.NewInt(0)})

// getHashEmptyForDepth calculates the hash value for an empty node at a given
// depth.
func getHashEmptyForDepth(depth int) *big.Int {
	h := zeroLeaf
	for i := 0; i < depth; i++ {
		h, _ = poseidon.Hash([]*big.Int{h, h})
	}
	return h
}

// NewSparseMerkleTree creates a new sparse Merkle tree with empty leaves.
func NewSparseMerkleTree(depth int) *SparseMerkleTree {
	emptyLeaves := make(map[string]*big.Int)
	root := &MerkleNode{Data: getHashEmptyForDepth(depth)}
	return &SparseMerkleTree{Root: root, Depth: depth, Leaves: emptyLeaves}
}

// Insert inserts a leaf with the given key and value into the tree.
func (smt *SparseMerkleTree) Insert(key string, value *big.Int) {
	smt.Leaves[key] = value
	smt.Root = smt.insertIntoNode(smt.Root, key, value, 0, smt.Depth)
}

// insertIntoNode inserts a leaf into the given node at the specified depth.
func (smt *SparseMerkleTree) insertIntoNode(node *MerkleNode, key string, value *big.Int, depth, maxDepth int) *MerkleNode {
	if node == nil {
		node = &MerkleNode{Data: getHashEmptyForDepth(maxDepth - depth)}
	}

	if depth == maxDepth {
		return &MerkleNode{Data: value}
	}

	pathBit := getPathBit(key, depth)
	if pathBit == 0 {
		node.Left = smt.insertIntoNode(node.getLeftChild(depth+1), key, value, depth+1, maxDepth)
	} else {
		node.Right = smt.insertIntoNode(node.getRightChild(depth+1), key, value, depth+1, maxDepth)
	}

	node.Data = hashChildren(node.Left, node.Right, maxDepth-depth)
	return node
}

// GenerateMerklePath generates a Merkle tree path for the leaf with the given key.
func (smt *SparseMerkleTree) GenerateMerklePath(key string) ([]*MerklePathItem, error) {
	if _, exists := smt.Leaves[key]; !exists {
		return nil, fmt.Errorf("no leaf exists at key: %s", key)
	}

	path := make([]*MerklePathItem, smt.Depth)
	current := smt.Root
	for depth := 0; depth < smt.Depth; depth++ {
		pathBit := getPathBit(key, depth)
		if pathBit == 0 {
			path[depth] = &MerklePathItem{
				SiblingHash: current.getRightChild(depth + 1).Data,
				IsRight:     true,
			}
			current = current.getLeftChild(depth + 1)
		} else {
			path[depth] = &MerklePathItem{
				SiblingHash: current.getLeftChild(depth + 1).Data,
				IsRight:     false,
			}
			current = current.getRightChild(depth + 1)
		}
	}

	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

// VerifyMerklePath verifies a Merkle tree path against the expected root hash.
func VerifyMerklePath(leafHash *big.Int, path []*MerklePathItem, expectedRoot *big.Int) bool {
	currentHash := leafHash
	for _, item := range path {
		siblingHash := item.SiblingHash

		if item.IsRight {
			currentHash, _ = poseidon.Hash([]*big.Int{currentHash, siblingHash})
			fmt.Println("currentHash", currentHash, "siblingHash", siblingHash)
		} else {
			currentHash, _ = poseidon.Hash([]*big.Int{siblingHash, currentHash})
		}
	}

	return currentHash.Cmp(expectedRoot) == 0
}

// getLeftChild returns the left child node of the current node.
func (node *MerkleNode) getLeftChild(depth int) *MerkleNode {
	if node.Left == nil {
		return &MerkleNode{Data: getHashEmptyForDepth(depth), Left: nil, Right: nil}
	}
	return node.Left
}

// getRightChild returns the right child node of the current node.
func (node *MerkleNode) getRightChild(depth int) *MerkleNode {
	if node.Right == nil {
		return &MerkleNode{Data: getHashEmptyForDepth(depth), Left: nil, Right: nil}
	}
	return node.Right
}

// hashChildren computes the hash value of two child nodes.
func hashChildren(left, right *MerkleNode, depth int) *big.Int {
	leftData := getHashEmptyForDepth(depth - 1)
	rightData := getHashEmptyForDepth(depth - 1)

	if left != nil {
		leftData = left.Data
	}

	if right != nil {
		rightData = right.Data
	}

	hash, _ := poseidon.Hash([]*big.Int{leftData, rightData})
	return hash
}

// getPathBit retrieves the bit value of the key at the specified depth.
func getPathBit(key string, depth int) int {
	if len(key) == 0 {
		return 0
	}
	i, _ := strconv.Atoi(key[depth : depth+1])
	return i
}

// getPaddedBinaryString returns a binary string representation of an integer,
// padded with leading zeros to a specified length.
func getPaddedBinaryString(i int, depth int) string {
	binStr := strconv.FormatInt(int64(i), 2)
	for len(binStr) < depth {
		binStr = "0" + binStr
	}
	return binStr
}

// NewDeterministicSparseMerkleTree creates a new deterministic sparse Merkle tree with non-null leaves.
func NewDeterministicSparseMerkleTree(depth int) *SparseMerkleTree {
	numLeaves := int(math.Pow(2, float64(depth)))
	smt := NewSparseMerkleTree(depth)
	for i := 0; i < numLeaves; i++ {
		key := getPaddedBinaryString(i, depth)
		leaf := big.NewInt(int64(i))
		smt.Insert(key, leaf)
	}

	return smt
}
