package merkletree

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
	leaves := make([]*big.Int, numLeaves)

	for i := 0; i < numLeaves; i++ {
		hashedLeaf, _ := poseidon.Hash([]*big.Int{big.NewInt(int64(i + startIndex))})
		leaves[i] = hashedLeaf
	}

	return NewMerkleTreeWithLeaves(leaves)
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
