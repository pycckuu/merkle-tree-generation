# Merkle Tree Generator

This Go project generates Merkle trees using the Poseidon hash function. The
trees can be created with a given depth or with specific leaves. The project
also includes command-line interface and JSON output functionality.

## Installation

Make sure you have installed [Go](https://golang.org/doc/install) (version 1.15
or later is recommended).

```bash
git clone https://github.com/pycckuu/merkle-tree-generation.git
cd merkle-tree-generation
go build .
```

## Usage
Command Line Interface You can run the program with specific parameters for the
Merkle tree depth and the number of leaves. The parameters can be passed as
command-line arguments:

```bash
./merkle-tree-generation -hLevel=4 -lLevel=16
```
This will generate a Merkle tree with a high-level of 4 and a low-level of 16.
The branches and the root of the tree will be printed to the console in JSON
format and saved to a file.

## JSON Output
The output JSON will have the following format:

```json
{
    "hLevel": 2,
    "lLevel": 16,
    "preimage": 1,
    "root": "0x2c370151f5ef741f065f0c4fc5c302f579cb52383b9d19e6d608bd25c2c76ab2",
    "branches": [
        "0x0c005cdbea16533de8615665f5490da32311c0a32f22e1e353a6d9f8a44419f8",
        "0x2fd4b719d39da0a5853ca2395f4b008303c107ac5ebedf5fdd4dbc73f75f31cb",
        "0x2f4d6625d3a809cb22cbb6da2c936d07dd14cc2fc213a4be397e0f301a9d7340",
        "0x00d8511fa0073158be570458e9da1bbac1e8fa0ea02ce46c451b520662b5836d"
    ]
}
```
Each branch and the root are represented as 32-byte hexadecimal strings.
