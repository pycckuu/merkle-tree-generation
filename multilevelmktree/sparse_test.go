package multilevelmktree

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSparseMerkleTree(t *testing.T) {
	smt := NewSparseMerkleTree(2)
	assert.NotNil(t, smt)
	assert.NotNil(t, smt.Root)
	assert.Equal(t, 2, smt.Depth)
	assert.Empty(t, smt.Leaves)

	tests := []struct {
		key          string
		value        *big.Int
		expectedRoot string
	}{
		{
			key:          "00",
			value:        big.NewInt(0),
			expectedRoot: "18366138217714291923534712849449091358386817997964088830897385671725623871073",
		},
		{
			key:          "01",
			value:        big.NewInt(1),
			expectedRoot: "8029606767784791880250783890079025673413177318829731918696248003381461757603",
		},
		{
			key:          "10",
			value:        big.NewInt(2),
			expectedRoot: "16218640559429857690153995608944582618510520484403789436977896806380962629939",
		},
		{
			key:          "11",
			value:        big.NewInt(3),
			expectedRoot: "3720616653028013822312861221679392249031832781774563366107458835261883914924",
		},
	}

	initRoot := new(big.Int)
	initRoot.SetString("2186774891605521484511138647132707263205739024356090574223746683689524510919", 10)
	if smt.Root.Data.Cmp(initRoot) != 0 {
		t.Error("Expected root node data to be", initRoot, "got", smt.Root.Data)
	}

	for _, test := range tests {
		smt.Insert(test.key, test.value)
		expectedRoot := new(big.Int)
		expectedRoot.SetString(test.expectedRoot, 10)
		if smt.Root.Data.Cmp(expectedRoot) != 0 {
			t.Error("Expected root node data to be", expectedRoot, "got", smt.Root.Data)
		}
	}
}

func TestInsert(t *testing.T) {
	smt := NewSparseMerkleTree(3)

	key := "000"
	value := big.NewInt(5)

	smt.Insert(key, value)

	assert.Equal(t, value, smt.Leaves[key])
}

func TestGetPaddedBinaryString(t *testing.T) {
	assert.Equal(t, "000", getPaddedBinaryString(0, 3))
	assert.Equal(t, "001", getPaddedBinaryString(1, 3))
	assert.Equal(t, "011", getPaddedBinaryString(3, 3))
	assert.Equal(t, "111", getPaddedBinaryString(7, 3))
}

func TestNewDeterministicSparseMerkleTree(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(3)
	assert.NotNil(t, smt)
	assert.NotNil(t, smt.Root)
	assert.Equal(t, 3, smt.Depth)
	assert.NotEmpty(t, smt.Leaves)
	assert.Len(t, smt.Leaves, 8)
}

// This test will depend on the poseidon.Hash function behavior.
func TestMerkleNodeHashes(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(3)

	// Test the root hash
	expectedRootHash := smt.Root.Data
	actualRootHash := hashChildren(smt.Root.Left, smt.Root.Right, smt.Depth)

	assert.Equal(t, expectedRootHash, actualRootHash)
}

func TestGenerateMerklePath(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(4)

	keys := []string{
		"0000",
		"0001",
		"0010",
		"0100",
		"1000",
		"1111",
	}

	for _, key := range keys {
		_, err := smt.GenerateMerklePath(key)
		assert.NoError(t, err, "Should not return an error for key %s", key)
	}

	// Test with non-existing key
	_, err := smt.GenerateMerklePath("10101")
	assert.Error(t, err, "Should return an error for non-existing key")
}

func TestSparseMerkleTree(t *testing.T) {
	depth := 4
	smt := NewDeterministicSparseMerkleTree(depth)
	fmt.Println(smt.Root.Data)
	fmt.Println(smt.Leaves)

	// i := 0
	for i := 0; i < (1 << depth); i++ {
		key := getPaddedBinaryString(i, depth)
		value := smt.Leaves[key]

		path, _ := smt.GenerateMerklePath(key)
		// print []MerklePathItem
		for _, item := range path {
			fmt.Println(item)
		}

		valid := VerifyMerklePath(value, path, smt.Root.Data)
		assert.True(t, valid, "The Merkle path should be valid for all leaves")
	}
}
