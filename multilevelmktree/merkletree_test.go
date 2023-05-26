package multilevelmktree

import (
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

func TestNewMerkleNode(t *testing.T) {
	// Test case for leaf node
	leafNode := NewMerkleNode(nil, nil, big.NewInt(1))

	if leafNode.Data.Cmp(big.NewInt(1)) != 0 {
		t.Error("Expected leaf node data to be 1, got ", leafNode.Data)
	}

	// Test case for non-leaf node
	left := NewMerkleNode(nil, nil, big.NewInt(1))
	right := NewMerkleNode(nil, nil, big.NewInt(2))

	// Hash of 1 and 2
	input := []*big.Int{big.NewInt(1), big.NewInt(2)}
	expected, _ := poseidon.Hash(input)

	nonLeafNode := NewMerkleNode(left, right, nil)

	if nonLeafNode.Data.Cmp(expected) != 0 {
		t.Error("Expected non-leaf node data to be hash of 1 and 2, got ", nonLeafNode.Data)
	}
}

func TestNewMerkleTree(t *testing.T) {
	// Test case for Merkle tree

	merkleTree := NewDeterministicMerkleTree(4, 1)

	if merkleTree == nil {
		t.Error("Expected new Merkle tree, got nil")
	}

	if merkleTree.Root == nil {
		t.Error("Expected root node, got nil")
	} else {
		if merkleTree.Root.Data == nil {
			t.Error("Expected root node data, got nil")
		}

		i := new(big.Int)
		i.SetString("12849909573197439023386719626541092579807164430016488237755007164956786115756", 10)
		if merkleTree.Root.Data.Cmp(i) != 0 {
			t.Error("Expected root node data to be", i, "got", merkleTree.Root.Data)
		}
	}
}
