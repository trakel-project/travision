package order

import (
	"encoding/gob"
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
var methodName = []string{"tsTotalNumOfOrder", "getOrderInfo0", "getOrderInfo1", "getOrderPlaceName", "getOrderDrivInfo"}
var testAccount = "0x38ba6f023d4a600f8a1cc566bba49190f3d8467e"
var hrpc *rpc.Rpc
var orderState = []string{"非法", "待抢单", "待支付", "完成", "终止"}

const (
	toDegree      = 0.000001
	degreeToMeter = 0.1
	toCoin        = 0.01
)

type address string
type orderAttributes struct {
	ID             uint64
	Passenger      address
	Driver         address
	StartPointX    float64
	StartPointY    float64
	DestPointX     float64
	DestPointY     float64
	StartName      string
	DestName       string
	Distance       float64
	PreFee         float64
	ActFeeDistance float64
	ActFeeTime     float64
	StartTime      time.Time
	PickTime       time.Time
	EndTime        time.Time
	State          int64
	PassInfo       string
	DriverInfo0    string
	DriverInfo1    string
	DriverInfo2    string
}

// Container is
type Container struct {
	Total uint64
	Data  map[uint64]orderAttributes
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
	// if c.Load() {
	// 	return nil
	// }
	if c.Data == nil {
		c.Data = make(map[uint64]orderAttributes)
	}
	return nil
}

// UpdateTotal returns the number of orders in blockchain.
func (c *Container) UpdateTotal() error {
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
	//更新第一部分
	wg.Add(4)
	var temp orderAttributes
	go func() error {
		defer wg.Done()
		ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[1], true, fmt.Sprintf("%d", idx))
		//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[1], ret, ret)
		if err != nil {
			colorLog.Error("Invoke to %q failed, error is %q", methodName[1], err)
			return err
		}
		//colorLog.Success("Invoke to %q success", methodName[1])
		retArray := ret.([]interface{})
		temp.ID = retArray[0].(*big.Int).Uint64()
		temp.Passenger = address(fmt.Sprintf("0x%x", retArray[1]))
		temp.StartPointX = float64(retArray[2].(*big.Int).Uint64()) * toDegree
		temp.StartPointY = float64(retArray[3].(*big.Int).Uint64()) * toDegree
		temp.DestPointX = float64(retArray[4].(*big.Int).Uint64()) * toDegree
		temp.DestPointY = float64(retArray[5].(*big.Int).Uint64()) * toDegree
		temp.Distance = float64(retArray[6].(*big.Int).Uint64()) * degreeToMeter
		temp.PreFee = float64(retArray[7].(*big.Int).Uint64()) * toCoin
		temp.StartTime = time.Unix(retArray[8].(*big.Int).Int64(), 0)
		temp.PassInfo = fmt.Sprintf("%s", retArray[9])
		return nil
	}()
	//更新第二部分
	go func() error {
		defer wg.Done()
		ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[2], true, fmt.Sprintf("%d", idx))
		//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[2], ret, ret)
		if err != nil {
			colorLog.Error("Invoke to %q failed, error is %q", methodName[2], err)
			return err
		}
		//colorLog.Success("Invoke to %q success", methodName[2])
		retArray := ret.([]interface{})
		temp.Driver = address(fmt.Sprintf("0x%x", retArray[0]))
		temp.ActFeeDistance = float64(retArray[1].(*big.Int).Uint64()) * toCoin
		temp.ActFeeTime = float64(retArray[2].(*big.Int).Uint64()) * toCoin
		temp.PickTime = time.Unix(retArray[3].(*big.Int).Int64(), 0)
		temp.EndTime = time.Unix(retArray[4].(*big.Int).Int64(), 0)
		temp.State = retArray[5].(*big.Int).Int64()
		return nil
	}()

	//更新第三部分
	go func() error {
		defer wg.Done()
		ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[3], true, fmt.Sprintf("%d", idx))
		//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[3], ret, ret)
		if err != nil {
			colorLog.Error("Invoke to %q failed, error is %q", methodName[3], err)
			return err
		}
		//colorLog.Success("Invoke to %q success", methodName[3])
		retArray := ret.([]interface{})
		temp.StartName = fmt.Sprintf("%s", retArray[0])
		temp.DestName = fmt.Sprintf("%s", retArray[1])
		return nil
	}()

	//更新第四部分
	go func() error {
		defer wg.Done()
		ret, err := hrpc.Invoke(ownerAddr, contractAddr, ownerPriKey, contractABI, methodName[4], true, fmt.Sprintf("%d", idx))
		//colorLog.Info("Invoke to %q, Return Type: %T, Return Value: %v", methodName[4], ret, ret)
		if err != nil {
			colorLog.Error("Invoke to %q failed, error is %q", methodName[4], err)
			return err
		}
		//colorLog.Success("Invoke to %q success", methodName[4])
		retArray := ret.([]interface{})
		temp.DriverInfo0 = fmt.Sprintf("%s", retArray[0])
		temp.DriverInfo1 = fmt.Sprintf("%s", retArray[1])
		temp.DriverInfo2 = fmt.Sprintf("%s", retArray[2])
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
	timeFormat := "2006年01月02日 15:04:05"
	fmt.Println("--------------------------------------")
	fmt.Printf("订单编号: %d\n", c.Data[idx].ID)
	fmt.Printf("乘客地址: %s\n", c.Data[idx].Passenger)
	fmt.Printf("司机地址: %s\n", c.Data[idx].Driver)
	fmt.Printf("起点经度: %f°\n", c.Data[idx].StartPointX)
	fmt.Printf("起点纬度: %f°\n", c.Data[idx].StartPointY)
	fmt.Printf("终点经度: %f°\n", c.Data[idx].DestPointX)
	fmt.Printf("终点纬度: %f°\n", c.Data[idx].DestPointY)
	fmt.Printf("起点地名: %s\n", c.Data[idx].StartName)
	fmt.Printf("终点地名: %s\n", c.Data[idx].DestName)
	fmt.Printf("距离: %.1f米\n", c.Data[idx].Distance)
	fmt.Printf("预付款: %.2f趣币\n", c.Data[idx].PreFee)
	fmt.Printf("距离费: %.2f趣币\n", c.Data[idx].ActFeeDistance)
	fmt.Printf("时长费: %.2f趣币\n", c.Data[idx].ActFeeTime)
	fmt.Println("订单提交时间: ", c.Data[idx].StartTime.Format(timeFormat))
	fmt.Println("行程开始时间: ", c.Data[idx].PickTime.Format(timeFormat))
	fmt.Println("行程结束时间: ", c.Data[idx].EndTime.Format(timeFormat))
	fmt.Printf("订单状态: %d %s\n", c.Data[idx].State, orderState[c.Data[idx].State])
	fmt.Printf("乘客信息: %s\n", c.Data[idx].PassInfo)
	fmt.Printf("司机姓名: %s\n", c.Data[idx].DriverInfo0)
	fmt.Printf("司机车牌: %s\n", c.Data[idx].DriverInfo1)
	fmt.Printf("司机车型: %s\n", c.Data[idx].DriverInfo2)
	fmt.Println("--------------------------------------")
}

// Save is
func (c *Container) Save() error {
	filePath := "./order.dat"
	f, err := os.Create(filePath)
	if err != nil {
		colorLog.Warning("Can't create Data file")
		fmt.Println(err)
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(*c); err != nil {
		colorLog.Warning("Can't encode Data")
		fmt.Println(err)
		return err
	}
	colorLog.Success("Save order data success, path is %q", filePath)
	return nil
}

// Load is
func (c *Container) Load() bool {
	filePath := "./order.dat"
	f, err := os.Open(filePath)
	if err != nil {
		colorLog.Warning("Can't open Data file, new Data file will be created in path \"./order.dat\" ")
		fmt.Println(err)
		return false
	}
	defer f.Close()

	enc := gob.NewDecoder(f)
	if err := enc.Decode(&c); err != nil {
		colorLog.Error("Can't decode Data file")
		fmt.Println(err)
		return false
	}
	colorLog.Success("Load order data success, path is %q", filePath)
	return true
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
