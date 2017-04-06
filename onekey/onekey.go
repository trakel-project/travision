package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"travision/gocode/colorLog"

	"strings"

	"git.hyperchain.cn/yeyc/hyperkit/rpc"
)

var ownerAddr = "0x2cd84f9e3c182c5c543571ea00611c41009c7024"
var ownerPriKey = "0x437cace9ccb62f0e3e5bd71d2793aa8ac4a0e9d42262028e4a4dc7797d060dff"

func main() {
	// 初始化rpc服务
	var hrpc *rpc.Rpc
	colorLog.Info("初始化中....")
	hrpc, err := rpc.NewRpc("http://114.55.64.145:8081", time.Second*10)
	if err != nil {
		colorLog.Error("无法连接服务器，错误代码为 %q", err)
		os.Exit(1)
	}
	colorLog.Success("连接服务器成功")

	// 编译Carcoin合约
	colorLog.Info("编译Carcoin合约")
	codeCarcoin, err := ioutil.ReadFile("./contract/carcoin.sol")
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	compiledCarcoin, err := hrpc.CompileContract(string(codeCarcoin))
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	//fmt.Println(compiledCarcoin.Abi[0])

	// 部署Carcoin合约
	colorLog.Info("部署Carcoin合约")
	addressCarcoin, err := hrpc.Deploy(ownerAddr, compiledCarcoin.Bin[0], ownerPriKey)
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	fmt.Println("Carcoin合约地址为：", addressCarcoin)

	// 编译Taxi合约
	colorLog.Info("编译Taxi合约")
	codeTaxing, err := ioutil.ReadFile("./contract/taxi.sol")
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	compiledTaxing, err := hrpc.CompileContract(string(codeTaxing))
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}

	// 部署Taxi合约
	colorLog.Info("部署Taxi合约")
	addressTaxing, err := hrpc.DeployWithArgs(ownerAddr, compiledTaxing.Bin[1], ownerPriKey, compiledTaxing.Abi[1], addressCarcoin)
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	fmt.Println("Taxi合约地址为：", addressTaxing)

	// 调用carcoin合约中的exeonece函数，将taxi合约地址写入
	colorLog.Info("调用carcoin合约中的exeonece函数，将taxi合约地址写入")
	_, err = hrpc.Invoke(ownerAddr, addressCarcoin, ownerPriKey, compiledCarcoin.Abi[0], "exeOnce", false, addressTaxing)
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}

	// 调用Taxi合约中verify函数，验证写入结果
	colorLog.Info("调用Taxi合约中verify函数，验证写入结果")
	ret, err := hrpc.Invoke(ownerAddr, addressTaxing, ownerPriKey, compiledTaxing.Abi[1], "verify", true)
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	str := fmt.Sprintf("0x%x", ret)
	fmt.Println("验证结果为：", str)
	if strings.Compare(str, addressTaxing) == 0 {
		colorLog.Success("验证成功，一键部署完毕")
	} else {
		colorLog.Error("验证失败，一键部署失败")
	}

	// 将合约地址和ABI写入文件
	result := fmt.Sprintf("Carcoin地址:\n%s\nCarcoin ABI:\n%s\nTaxing地址:\n%s\nTaxing ABI:\n%s\n", addressCarcoin, compiledCarcoin.Abi[0], addressTaxing, compiledTaxing.Abi[1])
	err = ioutil.WriteFile("./const.txt", []byte(result), 0666)
	if err != nil {
		colorLog.Error("%v", err)
		os.Exit(1)
	}
	colorLog.Info("合约地址和ABI已经保存至./const.txt，请手动至travision/gocode/contract/const.go中更新合约地址和ABI")

	//TODO 用得到的合约地址和ABI去更新travision后台golang代码中的const.go中的参数
}

// func insertAddress(code string, addr string) string {
// 	idx := strings.Index(code, "@@@@@@")
// 	bytes := []byte(code)
// 	newBytes := make([]byte, idx)
// 	copy(newBytes, bytes[0:idx])
// 	newBytes = append(newBytes, []byte(addr)...)
// 	newBytes = append(newBytes, bytes[idx+6:]...)
// 	return string(newBytes)
// }
