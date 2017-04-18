package main

import (
	"log"
	"sync"
	"travision/gocode/carcoin"
	"travision/gocode/driver"
	"travision/gocode/order"
)

func main() {
	updataOrders()
	// updataRecord()
	// updataDriver()
}

// Example for update order data from contract to local file
func updataOrders() {
	var orderData order.Container
	orderData.Initial()
	//defer orderData.Save()

	undoList := orderData.Check(1, orderData.Total)
	//fmt.Println(undoList)

	// use goroutine to accelerates
	// for large amount, YOU should limit the goroutine's quantity
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

	log.Println(orderData.Total)

	var i uint64
	for i = 1; i < orderData.Total; i++ {
		orderData.Print(i)
	}
}

func updataDriver() {
	var driverDate driver.Container
	driverDate.Initial()
	//defer orderData.Save()

	undoList := driverDate.Check(1, driverDate.Total)
	//fmt.Println(undoList)

	// use goroutine to accelerates
	// for large amount, YOU should limit the goroutine's quantity
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

	log.Println(driverDate.Total)

	var i uint64
	for i = 1; i < driverDate.Total; i++ {
		driverDate.Print(i)
	}
}

func updataRecord() {
	var recordDate carcoin.Container
	recordDate.Initial()
	//defer orderData.Save()

	undoList := recordDate.Check(1, recordDate.Total)
	//fmt.Println(undoList)

	// use goroutine to accelerates
	// for large amount, YOU should limit the goroutine's quantity
	var wg sync.WaitGroup
	wg.Add(len(undoList))
	for _, idx := range undoList {
		idx := idx
		go func() {
			defer wg.Done()
			recordDate.Insert(idx)
		}()
	}
	wg.Wait()

	log.Println(recordDate.Total)

	var i uint64
	for i = 1; i < recordDate.Total; i++ {
		recordDate.Print(i)
	}
}
