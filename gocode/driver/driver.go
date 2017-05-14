package driver

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"travision/gocode/colorLog"
	"travision/gocode/contract"

	"git.hyperchain.cn/yeyc/hyperkit/rpc"
)

// global varaible for contract invoke
var ownerAddr = contract.OwnerAddr
var ownerPriKey = contract.OwnerPriKey
var contractAddr = contract.Address
var contractABI = contract.ABI
var methodName = []string{"tsTotalNumOfDriver", "tsDriverInfo"}
var testAccount = "0x38ba6f023d4a600f8a1cc566bba49190f3d8467e"
var hrpc *rpc.Rpc
var driverState = []string{"接单中", "已抢单", "待接客", "行程中"}
var driverWorkingState = []string{"休息中", "接单中"}

const (
	toDegree      = 0.000001
	degreeToMeter = 0.1
	toCoin        = 0.01
)

type address string
type driverAttributes struct {
	ID           uint64
	CorX         float64
	CorY         float64
	IsWorking    int64
	Name         address
	Info0        string
	Info1        string
	Info2        string
	OrderPool    [8]uint64
	TotalJudge   int64
	AverageScore int64
	State        uint64
}

// Container is
type Container struct {
	Total uint64
	Data  map[uint64]driverAttributes
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
		c.Data = make(map[uint64]driverAttributes)
	}
	return nil
}

// UpdateTotal is
func (c *Container) UpdateTotal() error {
	colorLog.Info("Function call to \"UpdateTotal\"")

	ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[0], true)
	colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[0], ret, ret)
	if err != nil {
		colorLog.Error("Invoke to %q failed, error is %q", methodName[0], err)
		return err
	}
	colorLog.Success("Invoke to %q success", methodName[0])
	c.Total = ret.(*big.Int).Uint64()

	return nil
}

// Insert is
func (c *Container) Insert(idx uint64) error {
	colorLog.Info("Start query data of %d", idx)
	defer colorLog.Info("End query data of %d", idx)
	var wg sync.WaitGroup
	wg.Add(1)
	var temp driverAttributes
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
		//temp.IsWorking = retArray[0].(*big.Int)
		temp.IsWorking = 1
		temp.CorX = float64(retArray[1].(*big.Int).Uint64()) * toDegree
		temp.CorY = float64(retArray[2].(*big.Int).Uint64()) * toDegree
		temp.ID = idx
		temp.State = retArray[3].(*big.Int).Uint64()
		temp.TotalJudge = retArray[4].(*big.Int).Int64()
		temp.AverageScore = retArray[5].(*big.Int).Int64()
		temp.Name = address(fmt.Sprintf("0x%x", retArray[6]))
		temp.Info0 = fmt.Sprintf("%s", retArray[7])
		temp.Info1 = fmt.Sprintf("%s", retArray[8])
		temp.Info2 = fmt.Sprintf("%s", retArray[9])
		return nil
	}()
	wg.Wait()
	c.Mu.Lock()
	c.Data[temp.ID] = temp
	c.Mu.Unlock()
	return nil
}

// Print is
func (c *Container) Print(idx uint64) {
	fmt.Println("--------------------------------------")
	fmt.Printf("编号: %d\n", c.Data[idx].ID)
	fmt.Printf("经度: %f\n", c.Data[idx].CorX)
	fmt.Printf("纬度: %f\n", c.Data[idx].CorY)
	if c.Data[idx].IsWorking != 0 {
		fmt.Printf("状态: %s\n", driverState[c.Data[idx].State])
	} else {
		fmt.Printf("状态: %s\n", driverState[c.Data[idx].State])
	}
	fmt.Printf("地址: %s\n", c.Data[idx].Name)
	fmt.Printf("姓名: %s\n", c.Data[idx].Info0)
	fmt.Printf("车牌: %s\n", c.Data[idx].Info1)
	fmt.Printf("车型: %s\n", c.Data[idx].Info2)
	fmt.Printf("总单数: %d\n", c.Data[idx].TotalJudge)
	fmt.Printf("平均分: %.2f\n", float32(c.Data[idx].AverageScore/1000))
	fmt.Println("--------------------------------------")
}

// Check is
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
