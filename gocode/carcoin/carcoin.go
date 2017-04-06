package carcoin

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"travision/gocode/colorLog"
	"travision/gocode/contract"

	"sync"

	"git.hyperchain.cn/yeyc/hyperkit/rpc"
)

// global varaible for contract invoke
var ownerAddr = contract.OwnerAddr
var ownerPriKey = contract.OwnerPriKey
var contractAddr = contract.CarcoinAddr
var contractABI = contract.CarcoinABI
var methodName = []string{"tsTotalNumOfRecord", "tsRecords"}
var testAccount = "0x38ba6f023d4a600f8a1cc566bba49190f3d8467e"
var hrpc *rpc.Rpc

const (
	toDegree      = 0.000001
	degreeToMeter = 0.1
	toCoin        = 0.01
)

type address string
type recordAttributes struct {
	From    address
	To      address
	Amount  int64
	Comment string
}

// Container is
type Container struct {
	Total uint64
	Data  map[uint64]recordAttributes
	Mu    sync.Mutex
}

func init() {
	colorLog.Info("Initializing....")
	var err error
	hrpc, err = rpc.NewRpc(contract.IP, time.Second*10)
	if err != nil {
		colorLog.Error("Connect to server fail, error is %q", err)
		os.Exit(1)
	}
	colorLog.Success("Connect to server success")
}

// Initial is
func (c *Container) Initial() error {
	if err := c.UpdateTotal(); err != nil {
		return err
	}
	if c.Data == nil {
		c.Data = make(map[uint64]recordAttributes)
	}
	return nil
}

// UpdateTotal returns the number of orders in blockchain.
func (c *Container) UpdateTotal() error {
	ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[0], true, fmt.Sprintf("%d", 1))
	//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[0], ret, ret)
	if err != nil {
		colorLog.Error("Invoke to %q failed, error is %q", methodName[0], err)
		return err
	}
	//colorLog.Success("Invoke to %q success", methodName[0])
	c.Total = ret.(*big.Int).Uint64()
	return nil
}

// Insert is
func (c *Container) Insert(idx uint64) error {
	colorLog.Info("Start query data of %d", idx)
	defer colorLog.Info("End query data of %d", idx)
	var wg sync.WaitGroup
	wg.Add(1)
	var temp recordAttributes
	go func() error {
		defer wg.Done()
		ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[1], true, fmt.Sprintf("%d", idx))
		//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", method, ret, ret)
		if err != nil {
			colorLog.Error("Invoke to %q failed, error is %q", methodName[1], err)
			return err
		}
		//colorLog.Success("Invoke to %q success", method)
		retArray := ret.([]interface{})
		temp.From = address(fmt.Sprintf("0x%x", retArray[0]))
		temp.To = address(fmt.Sprintf("0x%x", retArray[1]))
		temp.Amount = retArray[2].(*big.Int).Int64()
		temp.Comment = fmt.Sprintf("%s", retArray[3])
		return nil
	}()
	wg.Wait()
	c.Mu.Lock()
	c.Data[idx] = temp
	c.Mu.Unlock()
	return nil
}

// Print is
func (c *Container) Print(idx uint64) {
	fmt.Println("--------------------------------------")
	fmt.Printf("From: %s\n", c.Data[idx].From)
	fmt.Printf("To: %s\n", c.Data[idx].To)
	fmt.Printf("数额: %d\n", c.Data[idx].Amount)
	fmt.Printf("备注: %s\n", c.Data[idx].Comment)
	fmt.Println("--------------------------------------")
}

// Check checks whether the orders indexed in [i, j) exists in local data file
// and return the non-existed order ID
func (c *Container) Check(i uint64, j uint64) []uint64 {
	var re []uint64
	for ; i < j; i++ {
		if i >= c.Total {
			return re
		}
		if _, ok := c.Data[i]; !ok {
			re = append(re, i)
		}
	}
	return re
}
