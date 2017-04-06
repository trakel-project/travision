package controllers

import (
	"encoding/json"
	"fmt"
	"sync"
	"travision/gocode/carcoin"
	"travision/gocode/driver"
	"travision/gocode/order"

	"github.com/astaxie/beego"
)

//----MainController----//
type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

//----OrderController----//
type OrderController struct {
	beego.Controller
}

func (c *OrderController) Get() {
	c.TplName = "order.html"
	str := quearyOrder()
	fmt.Println(str)
	c.Data["orders"] = str
}

type orderJSON struct {
	ID        uint64 `json:"id"`
	StartName string `json:"sname"`
	DestName  string `json:"dname"`
	Distance  string `json:"distance"`
	PreFee    string `json:"prefee"`
	ActFee    string `json:"actfee"`
	Duration  string `json:"duration"`
	State     string `json:"state"`
}

func quearyOrder() string {
	var orderData order.Container
	orderData.Initial()
	undoList := orderData.Check(1, orderData.Total)
	var wg sync.WaitGroup
	wg.Add(len(undoList))
	for _, idx := range undoList {
		idx := idx
		go func() {
			defer wg.Done()
			orderData.Insert(idx)
		}()
	}
	wg.Wait()

	var orderState = []string{"非法", "待抢单", "待支付", "完成", "终止"}
	var result []orderJSON
	var o orderJSON
	var i uint64
	for i = 1; i < orderData.Total; i++ {
		o.ID = orderData.Data[i].ID
		o.StartName = orderData.Data[i].StartName
		o.DestName = orderData.Data[i].DestName
		o.PreFee = fmt.Sprintf("%.2f趣币", orderData.Data[i].PreFee)
		o.ActFee = fmt.Sprintf("%.2f趣币", orderData.Data[i].ActFeeDistance+orderData.Data[i].ActFeeTime)
		duration := (orderData.Data[i].EndTime.Unix() - orderData.Data[i].PickTime.Unix()) / 60
		if duration < 0 {
			duration = 0
		}
		o.Duration = fmt.Sprintf("%d分钟", duration)
		o.Distance = fmt.Sprintf("%.2f千米", orderData.Data[i].Distance/1000)
		o.State = fmt.Sprintf("%s\n", orderState[orderData.Data[i].State])
		result = append(result, o)
	}
	body, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(body)
}

//----DriverController----//
type DriverController struct {
	beego.Controller
}

func (c *DriverController) Get() {
	c.TplName = "driver.html"
	str := quearyDriver()
	fmt.Println(str)
	c.Data["drivers"] = str
}

type driverJSON struct {
	ID           uint64 `json:"id"`
	CorX         string `json:"x"`
	CorY         string `json:"y"`
	Name         string `json:"address"`
	Info         string `json:"info"`
	TotalJudge   int64  `json:"toljudge"`
	AverageScore string `json:"avgscore"`
	State        string `json:"state"`
}

func quearyDriver() string {
	var driverDate driver.Container
	driverDate.Initial()
	undoList := driverDate.Check(1, driverDate.Total)
	var wg sync.WaitGroup
	wg.Add(len(undoList))
	for _, idx := range undoList {
		idx := idx
		go func() {
			defer wg.Done()
			driverDate.Insert(idx)
		}()
	}
	wg.Wait()

	var driverState = []string{"挂起", "已抢单", "待接客", "行程中"}
	var result []driverJSON
	var o driverJSON
	var i uint64
	for i = 1; i < driverDate.Total; i++ {
		o.ID = driverDate.Data[i].ID
		o.CorX = fmt.Sprintf("%.6f", driverDate.Data[i].CorX)
		o.CorY = fmt.Sprintf("%.6f", driverDate.Data[i].CorY)
		o.Name = string(driverDate.Data[i].Name)
		o.Info = driverDate.Data[i].Info2
		o.TotalJudge = driverDate.Data[i].TotalJudge
		o.AverageScore = fmt.Sprintf("%.1f", float64(driverDate.Data[i].AverageScore)/1000.0)
		o.State = driverState[driverDate.Data[i].State]
		result = append(result, o)
	}
	body, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(body)
}

//----CarcoinController----//
type CarcoinController struct {
	beego.Controller
}

func (c *CarcoinController) Get() {
	c.TplName = "carcoin.html"
	str := quearyCarcoin()
	fmt.Println(str)
	c.Data["carcoins"] = str
}

type carcoinJSON struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
	Comment string `json:"comment"`
}

func quearyCarcoin() string {
	var carcoinData carcoin.Container
	carcoinData.Initial()
	undoList := carcoinData.Check(1, carcoinData.Total)
	var wg sync.WaitGroup
	wg.Add(len(undoList))
	for _, idx := range undoList {
		idx := idx
		go func() {
			defer wg.Done()
			carcoinData.Insert(idx)
		}()
	}
	wg.Wait()

	var result []carcoinJSON
	var o carcoinJSON
	var i uint64
	for i = 1; i < carcoinData.Total; i++ {
		o.From = string(carcoinData.Data[i].From)
		o.To = string(carcoinData.Data[i].To)
		o.Amount = fmt.Sprintf("%.2f趣币", float64(carcoinData.Data[i].Amount)/100.0)
		o.Comment = carcoinData.Data[i].Comment
		result = append(result, o)
	}
	body, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(body)
}
