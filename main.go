package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int
}

type Blockchain struct {
	blocks     []*Block
	difficulty int
}

// -------- Block Logic --------

func (b *Block) prepData() []byte {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	nonce := []byte(strconv.Itoa(b.Nonce))

	return bytes.Join(
		[][]byte{
			b.PrevHash,
			b.Data,
			timestamp,
			nonce,
		},
		[]byte{},
	)
}

func (b *Block) Mine(difficulty int) {
	target := strings.Repeat("0", difficulty)

	for {
		hash := sha256.Sum256(b.prepData())
		hashHex := fmt.Sprintf("%x", hash)

		if strings.HasPrefix(hashHex, target) {
			b.Hash = hash[:]
			break
		}
		b.Nonce++
	}
}

func NewBlock(data string, prevHash []byte, difficulty int) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte(data),
		PrevHash:  prevHash,
		Hash:      []byte{},
		Nonce:     0,
	}

	block.Mine(difficulty)
	return block
}

// -------- Blockchain Logic --------

func NewBlockchain(difficulty int) *Blockchain {
	genesis := NewBlock("Genesis Block", []byte{}, difficulty)
	return &Blockchain{
		blocks:     []*Block{genesis},
		difficulty: difficulty,
	}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, bc.difficulty)
	bc.blocks = append(bc.blocks, newBlock)
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.blocks); i++ {
		prev := bc.blocks[i-1]
		curr := bc.blocks[i]

		if !bytes.Equal(curr.PrevHash, prev.Hash) {
			return false
		}

		hash := sha256.Sum256(curr.prepData())
		if !bytes.Equal(curr.Hash, hash[:]) {
			return false
		}
	}
	return true
}

// -------- Main --------

func main() {
	bc := NewBlockchain(3)

	bc.AddBlock("Sending 1 ETH to Charles")
	bc.AddBlock("Sending 3 ETH to Alice")

	for i, block := range bc.blocks {
		fmt.Printf("Block %d\n", i)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}

	fmt.Println("Blockchain valid?", bc.IsValid())
}
