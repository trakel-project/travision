# 一键部署说明   

#### 目录说明   
`contract`中为合约代码，carcoin.sol是taxi.sol中第一部分

#### 使用说明   
根据链IP修改main函数下newRPC的参数      
编译运行onekey.go，获得carcoin和Taxing两个合约的地址和ABI，并且会保存在`./const.txt`文件中，用户根据文件的内容去更新`gocode/contract/const.go`中的全局参数    
之后方可正常使用监控平台

**注意双引号"应该使用转义字符\"**



#### 一键部署功能拆解   
- 读取合约代码
- 编译carcoin合约代码返回ABI
- 部署carcoin合约代码返回合约地址
- 将carcoin合约的地址写入taxing合约中
- 编译carcoin合约代码返回ABI
- 部署taxing合约得到地址
- 调用carcoin合约中的exeonece函数，将taxing合约地址写入
- 用得到的合约地址和ABI去更新travision后台golang代码中的const.go中的参数