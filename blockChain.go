package main

type BlockChain struct {
	_blocks []*Block
}

func (this *BlockChain) AddBlock(data []byte) {
	plastBlock := this._blocks[len(this._blocks)-1]
	prevHash := plastBlock._hash
	pblock := NewBlock(prevHash, data)
	this._blocks = append(this._blocks, pblock)
}

func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	blockChain := BlockChain{
		_blocks: []*Block{genesisBlock},
	}
	return &blockChain
}

func GenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("genesisBlock"))
}
