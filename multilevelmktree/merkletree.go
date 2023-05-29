package multilevelmktree

import (
	"math"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  *big.Int
}

type MerkleTree struct {
	Root *MerkleNode
}

func NewMerkleNode(left, right *MerkleNode, data *big.Int) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		mNode.Data = data
	} else {
		// Hash the concatenation of the left and right data
		input := []*big.Int{left.Data, right.Data}
		hashed, _ := poseidon.Hash(input)

		mNode.Data = hashed
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}

func NewDeterministicMerkleTree(depth int, startIndex int) *MerkleTree {
	numLeaves := int(math.Pow(2, float64(depth)))
	var numBranches int
	if depth > 6 {
		numBranches = int(math.Pow(2, float64(int64(depth-6)))) // Assuming 64 branches
	} else {
		numBranches = 1
	}

	branchRoots := make([]*big.Int, 0, numBranches)

	for i := 0; i < numBranches; i++ {
		// For each branch, generate the leaves and build the Merkle tree
		branchLeaves := make([]*big.Int, 0, numLeaves/numBranches)
		for j := 0; j < numLeaves/numBranches; j++ {
			leaf, _ := poseidon.Hash([]*big.Int{big.NewInt(int64((i * numLeaves / numBranches) + j + startIndex))})
			branchLeaves = append(branchLeaves, leaf)
		}

		branch := NewMerkleTreeWithLeaves(branchLeaves)
		branchRoots = append(branchRoots, branch.Root.Data)
	}

	return NewMerkleTreeWithLeaves(branchRoots)
}

func NewMerkleTreeWithLeaves(leaves []*big.Int) *MerkleTree {
	nodes := make([]MerkleNode, 0, len(leaves))

	for _, leaf := range leaves {
		node := NewMerkleNode(nil, nil, leaf)
		nodes = append(nodes, *node)
	}

	depth := int(math.Log2(float64(len(leaves))))
	for i := 0; i < depth; i++ {
		newLevel := make([]MerkleNode, 0, len(nodes)/2)

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}
