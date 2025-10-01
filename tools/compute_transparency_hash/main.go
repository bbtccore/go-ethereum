package main

import (
    "fmt"
    "github.com/ethereum/go-ethereum/core"
)

func main() {
    g := core.DefaultTransparencyGenesisBlock()
    h := g.ToBlock().Hash()
    fmt.Println(h.String())
}

