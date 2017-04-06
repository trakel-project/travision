package rpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"encoding/hex"

	"git.hyperchain.cn/yeyc/hyperkit/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type Rpc struct {
	client client
}

func NewRpc(endpoint string, timeout time.Duration) (*Rpc, error) {
	client, err := newHTTPClient(endpoint, timeout)
	if err != nil {
		return nil, err
	}
	return &Rpc{client: client}, nil
}

func (r *Rpc) call(method string, params []interface{}) (json.RawMessage, error) {
	req := &jsonRequest{
		Method:  method,
		Version: JSONRPCVersion,
		Id:      1,
		Params:  params,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Rpc: %v", err)
	}

	data, err := r.client.Send(body)
	if err != nil {
		//err e.g:Post http://localhost:8089: dial tcp 127.0.0.1:8089: getsockopt: connection refused
		return nil, err
	}

	var resp jsonResponse
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("Rpc: %v", err)
	}
	if resp.Result != nil {
		return resp.Result, nil
	}
	return nil, fmt.Errorf("Rpc Server: %v", resp.Error.Message)
}

/*---------------------------------- node -------------------------------------*/

func (r *Rpc) GetNodes() ([]Node, error) {
	data, err := r.call("node_getNodes", []interface{}{})
	if err != nil {
		return nil, err
	}
	var nodes []Node
	if err = json.Unmarshal(data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

/*---------------------------------- block -------------------------------------*/

func (r *Rpc) GetLatestBlock() (*Block, error) {
	data, err := r.call("block_latestBlock", []interface{}{})
	if err != nil {
		return nil, err
	}

	var block BlockRaw
	if err = json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return block.ToBlock()
}

// @param to, int or "latest"
func (r *Rpc) GetBlocks(from int, to interface{}) ([]Block, error) {
	obj := make(map[string]interface{})
	obj["from"], obj["to"] = from, to
	data, err := r.call("block_getBlocks", []interface{}{obj})
	if err != nil {
		return nil, err
	}

	var blocksRaw []BlockRaw
	if err = json.Unmarshal(data, &blocksRaw); err != nil {
		return nil, err
	}
	var blocks []Block
	for _, b := range blocksRaw {
		if block, err := b.ToBlock(); err != nil {
			return nil, err
		} else {
			blocks = append(blocks, *block)
		}
	}

	return blocks, nil
}

func (r *Rpc) GetBlockByNumber(number int) (*Block, error) {
	data, err := r.call("block_getBlockByNumber", []interface{}{number})
	if err != nil {
		return nil, err
	}
	var blocksRaw BlockRaw
	if err = json.Unmarshal(data, &blocksRaw); err != nil {
		return nil, err
	}
	return blocksRaw.ToBlock()
}

func (r *Rpc) GetBlockByHash(hash string) (*Block, error) {
	data, err := r.call("block_getBlockByHash", []interface{}{hash})
	if err != nil {
		return nil, err
	}
	var blocksRaw BlockRaw
	if err = json.Unmarshal(data, &blocksRaw); err != nil {
		return nil, err
	}
	return blocksRaw.ToBlock()
}

/*--------------------------------- transaction --------------------------*/
// @param txhash, string in form: "0x....."
func (r *Rpc) GetTransactionByHash(txhash string) (*Transaction, error) {
	data, err := r.call("tx_getTransactionByHash", []interface{}{txhash})
	if err != nil {
		return nil, err
	}

	var tx TransactionRaw
	if err = json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}
	return tx.ToTransaction()
}

// @param txhash, string in form: "0x....."
func (r *Rpc) GetTxReceipt(txhash string) (*TxReceipt, error) {
	data, err := r.call("tx_getTransactionReceipt", []interface{}{txhash})
	if err != nil {
		return nil, err
	}

	var txr TxReceipt
	if err = json.Unmarshal(data, &txr); err != nil {
		return nil, err
	}
	return &txr, nil
}

// @param timeout, int for time.Second
func (r *Rpc) GetTxReceiptPolling(txhash string, timeout int) (*TxReceipt, error) {
	if timeout > 5 || timeout <= 0 {
		timeout = 5
	}
	for i := 0; i < timeout; i++ {
		time.Sleep(time.Second)
		txr, err := r.GetTxReceipt(txhash)
		if err != nil {
			return nil, err
		}
		if txr.TxHash != "" {
			return txr, nil
		}
	}
	return &TxReceipt{TxHash: txhash}, fmt.Errorf("Polling TxReceipt timeout after %v seconds", timeout)
}

func (r *Rpc) GetSighHash(from, to, nonce, value, payload, timestamp interface{}) (string, error) {
	obj := make(map[string]interface{})
	obj["from"], obj["to"], obj["nonce"], obj["value"] = from, to, nonce, value
	obj["payload"], obj["timestamp"] = payload, timestamp
	data, err := r.call("tx_getSignHash", []interface{}{obj})
	if err != nil {
		return "", err
	}

	var hash string
	if err = json.Unmarshal(data, &hash); err != nil {
		return "", err
	}
	return hash, nil
}

/*-------------------------------- contract ----------------------------------*/
func (r *Rpc) CompileContract(code string) (*CompileResult, error) {
	data, err := r.call("contract_compileContract", []interface{}{code})
	if err != nil {
		return nil, err
	}

	var cr CompileResult
	if err = json.Unmarshal(data, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

func (r *Rpc) DeployContract(from, payload, privateKey string) (txhash string, err error) {
	var (
		timestamp = time.Now().UnixNano()
		nonce     = rand.Int63()
		hash      string
		sig       []byte
	)
	if hash, err = r.GetSighHash(from, nil, nonce, nil, payload, timestamp); err != nil {
		return "", err
	}
	if sig, err = secp256k1.Sign(common.FromHex(hash), common.FromHex(privateKey)); err != nil {
		return "", err
	}

	obj := make(map[string]interface{})
	obj["from"], obj["nonce"], obj["timestamp"], obj["payload"], obj["signature"] = from, nonce, timestamp, payload, common.ToHex(sig)
	data, err := r.call("contract_deployContract", []interface{}{obj})
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(data, &txhash); err != nil {
		return "", err
	}
	return txhash, nil
}

func (r *Rpc) InvokeContract(from, to, payload, privateKey string, Const bool) (txhash string, err error) {
	var (
		timestamp = time.Now().UnixNano()
		nonce     = rand.Int63()
		hash      string
		sig       []byte
	)
	if hash, err = r.GetSighHash(from, to, nonce, nil, payload, timestamp); err != nil {
		return "", err
	}
	if sig, err = secp256k1.Sign(common.FromHex(hash), common.FromHex(privateKey)); err != nil {
		return "", err
	}

	obj := make(map[string]interface{})
	obj["from"], obj["to"], obj["nonce"], obj["timestamp"] = from, to, nonce, timestamp
	obj["payload"], obj["signature"], obj["simulate"] = payload, common.ToHex(sig), Const
	data, err := r.call("contract_invokeContract", []interface{}{obj})
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(data, &txhash); err != nil {
		return "", err
	}
	return txhash, nil
}

// Deploy Contract Failed: return nil, err
// After Deploy Success, use the txhash to get TxReceipt in polling way with timeout
// If fail to get TxReceipt, return txhash, err
// If success to get TxReceipt, return ContractAddress, nil
func (r *Rpc) Deploy(from, payload, privateKey string) (string, error) {
	txhash, err := r.DeployContract(from, payload, privateKey)
	if err != nil {
		return "", err
	}

	txr, err := r.GetTxReceiptPolling(txhash, 5)
	if err != nil {
		return txhash, err
	}

	return txr.ContractAddress, nil
}

// Encode constructor params, append them to payload(bin)
func (r *Rpc) DeployWithArgs(from, payload, privateKey, abiString string, args ...interface{}) (string, error) {
	ABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return "", err
	}

	packed, err := ABI.PackJSON("", args...)
	if err != nil {
		return "", err
	}

	payload = payload + hex.EncodeToString(packed)

	return r.Deploy(from, payload, privateKey)
}

// Invoke Contract Failed: return nil, err
// After Invoke Success, use the txhash to get TxReceipt in polling way with timeout
// If fail to get TxReceipt, return &TxReceipt{TxHash: txhash}, err
// If success to get TxReceipt, unpack ret to readable form
// MultiReturn unpack to []interface{} form
// SingleReturn unpack to interface{} form
// NoReturn return nil
// Unpack Failed: return *TxReceipt, nil
func (r *Rpc) Invoke(from, to, privateKey, abiString, method string, Const bool, args ...interface{}) (interface{}, error) {
	ABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}

	packed, err := ABI.PackJSON(method, args...)
	if err != nil {
		return nil, err
	}

	txhash, err := r.InvokeContract(from, to, common.ToHex(packed), privateKey, Const)
	if err != nil {
		return nil, err
	}

	txr, err := r.GetTxReceiptPolling(txhash, 5)
	if err != nil {
		return txhash, err
	}

	switch {
	case len(ABI.Methods[method].Outputs) == 1:
		var v interface{}
		if err = ABI.Unpack(&v, method, common.FromHex(txr.Ret)); err != nil {
			return txhash, err
		}
		return v, nil
	case len(ABI.Methods[method].Outputs) > 1:
		var v []interface{}
		if err = ABI.Unpack(&v, method, common.FromHex(txr.Ret)); err != nil {
			return txhash, err
		}
		return v, nil
	default:
		return txhash, nil
	}
}
