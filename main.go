package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sync"

	merkletree "github.com/pycckuu/merkle-tree-generation/src"
	"github.com/schollz/progressbar/v3"
)

type Output struct {
	HLevel   int      `json:"hLevel"`
	LLevel   int      `json:"lLevel"`
	PreImage int      `json:"preimage"`
	Root     string   `json:"root"`
	Branches []string `json:"branches"`
}

// getMerkleRoots computes the Merkle tree roots for each branch concurrently
func getMerkleRoots(hLevel, lLevel int, preImage int) []*big.Int {
	n := int(math.Pow(2, float64(hLevel)))
	increment := int(math.Pow(2, float64(lLevel)))
	branches := make([]*big.Int, n)

	bar := progressbar.Default(int64(n))

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			merkleTree := merkletree.NewDeterministicMerkleTree(lLevel, (i+preImage)*increment)
			branches[i] = merkleTree.Root.Data
			bar.Add(1)
		}(i)
	}

	wg.Wait()

	return branches
}

// outputJSON formats the output as JSON and prints to stdout
func outputJSON(branches []*big.Int, root *big.Int, hLevel, lLevel int, preImage int) {
	branchesHex := make([]string, len(branches))
	for i, branch := range branches {
		branchesHex[i] = fmt.Sprintf("0x%064s", branch.Text(16))
	}
	rootHex := fmt.Sprintf("0x%064s", root.Text(16))

	output := Output{
		Branches: branchesHex,
		HLevel:   hLevel,
		PreImage: preImage,
		Root:     rootHex,
		LLevel:   lLevel,
	}

	outputJSON, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s\n", outputJSON)

	// Open output file
	fileName := fmt.Sprintf("output_hLevel_%d_lLevel_%d_preImage_%d.json", hLevel, lLevel, preImage)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// Write JSON data to the file
	_, err = file.Write(outputJSON)
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}

	fmt.Println("Output written to", fileName)
}

func main() {
	// Define the flags
	hLevelPtr := flag.Int("hLevel", 4, "An integer value for the hLevel")
	lLevelPtr := flag.Int("lLevel", 16, "An integer value for the lLevel")
	preimagePtr := flag.Int("preImage", 0, "An integer value for the preimage")

	// Parse the flags
	flag.Parse()

	hLevel := *hLevelPtr
	lLevel := *lLevelPtr
	preImage := *preimagePtr

	branches := getMerkleRoots(hLevel, lLevel, preImage)
	root := merkletree.NewMerkleTreeWithLeaves(branches).Root.Data

	outputJSON(branches, root, hLevel, lLevel, preImage)
}
