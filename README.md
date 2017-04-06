# 趣快出行监控平台

### 简介
采用Beego作为web框架，访问[Beego官网](https://beego.me/)来获得更多相关的信息   
利用该平台可以查询链上的所有订单、司机还有转账记录信息   

### 部署说明
1. 配置Golang环境   
访问[Golang官网](https://golang.org/) 安装对应的Goalng运行环境，并且设置好$GOROOT和$GOPATH   
**Golang版本应该高于1.8**   
项目文件夹应该放置在src目录下
一个设定好的$GOPATH应该有如下的目录组织   
$GOPATH    
&emsp;&emsp;src   
&emsp;&emsp;&emsp;&emsp;travision   
`cd $GOPATH/src/travision`

2. 配置Beego包(可选)   
参见[Beego官网](https://beego.me/)   
若不需要进一步对web框架进行开发，也可以直接使用项目内置的包，参见3

3. 配置相应的包   
编译所需要的依赖库都包含在`vendor`文件夹下  
`cd vendor`   
`cp -rf . $GOPATH/src`       
注意：如果已经配置了Beego环境，则忽略`vendor/github.com/astaxie/beego`

4. 修改IP   
根据部署链的服务器的IP修改`travision/gocode/const/const.go`中IP参数   

5. 部署合约   
合约代码存放在`onekey/contract/`中   
第一次使用时我们需要重新部署合约，并更新
`gocode/contract/const.go`
中的参数，具体参看`onekey`目录下README  

6. 启动服务   
`cd $GOPATH/src/travision`   
`go build`   
`./travision`   
即可打开服务，开始监听8080端口。浏览器[打开](http://localhost:8080)即可使用监控平台。   

### 特别说明
1. 平台每次刷新都会从合约中获取最新的数据，所以刷新时间较长是正常的，可以从运行`./travision`的终端看到数据读取的进度   

2. `example`中为一个完整流程的例子，可以参考   
