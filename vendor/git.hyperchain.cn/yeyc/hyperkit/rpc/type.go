package rpc

import "strconv"

type Node struct {
	Status    uint
	Ip        string
	Port      uint
	Id        uint
	Isprimary bool
	Delay     uint
}

type BlockRaw struct {
	Number       string
	Hash         string
	ParentHash   string
	WriteTime    uint64
	AvgTime      string
	Txcounts     string
	MerkleRoot   string
	Transactions []TransactionRaw
}

type Block struct {
	Number       uint64
	Hash         string
	ParentHash   string
	WriteTime    uint64
	AvgTime      uint64
	Txcounts     uint64
	MerkleRoot   string
	Transactions []Transaction
}

type TransactionRaw struct {
	Hash        string
	BlockNumber string
	BlockHash   string
	TxIndex     string
	From        string
	To          string
	Amount      string
	Timestamp   uint64
	ExecuteTime string
	Invalid     bool
	InvalidMsg  string
}

type Transaction struct {
	Hash        string
	BlockNumber uint64
	BlockHash   string
	TxIndex     uint64
	From        string
	To          string
	Amount      uint64
	Timestamp   uint64
	ExecuteTime uint64
	Invalid     bool
	InvalidMsg  string
}

type TxReceipt struct {
	TxHash          string
	PostState       string
	ContractAddress string
	Ret             string
}

type CompileResult struct {
	Abi   []string
	Bin   []string
	Types []string
}

func (b *BlockRaw) ToBlock() (*Block, error) {
	var (
		Number       uint64
		AvgTime      uint64
		Txcounts     uint64
		Transactions []Transaction
		err          error
	)
	if Number, err = strconv.ParseUint(b.Number, 0, 64); err != nil {
		return nil, err
	}
	if AvgTime, err = strconv.ParseUint(b.AvgTime, 0, 64); err != nil {
		return nil, err
	}
	if Txcounts, err = strconv.ParseUint(b.Txcounts, 0, 64); err != nil {
		return nil, err
	}
	for _, t := range b.Transactions {
		if transaction, err := t.ToTransaction(); err != nil {
			return nil, err
		} else {
			Transactions = append(Transactions, *transaction)
		}
	}
	return &Block{
		Number:       Number,
		Hash:         b.Hash,
		ParentHash:   b.ParentHash,
		WriteTime:    b.WriteTime,
		AvgTime:      AvgTime,
		Txcounts:     Txcounts,
		MerkleRoot:   b.MerkleRoot,
		Transactions: Transactions,
	}, nil
}

func (t *TransactionRaw) ToTransaction() (*Transaction, error) {
	var (
		BlockNumber uint64
		TxIndex     uint64
		Amount      uint64
		ExecuteTime uint64
		err         error
	)
	if BlockNumber, err = strconv.ParseUint(t.BlockNumber, 0, 64); err != nil {
		return nil, err
	}
	if TxIndex, err = strconv.ParseUint(t.TxIndex, 0, 64); err != nil {
		return nil, err
	}
	if Amount, err = strconv.ParseUint(t.Amount, 0, 64); err != nil {
		return nil, err
	}
	if ExecuteTime, err = strconv.ParseUint(t.ExecuteTime, 0, 64); err != nil {
		return nil, err
	}
	return &Transaction{
		Hash:        t.Hash,
		BlockNumber: BlockNumber,
		BlockHash:   t.BlockHash,
		TxIndex:     TxIndex,
		From:        t.From,
		To:          t.To,
		Amount:      Amount,
		Timestamp:   t.Timestamp,
		ExecuteTime: ExecuteTime,
		Invalid:     t.Invalid,
		InvalidMsg:  t.InvalidMsg,
	}, nil
}
