# 一个简单的例子   

本文展现如何编译、部署智能合约，并且使用travision监控平台来监控数据。   
使用的平台是Ubuntu 14.04LTS 32位系统。

1. 配置Golang运行环境   
由于Ubuntu的apt-get中Golang版本过低，我们需要从官网下载最新的版本进行安装   
`cd`   
`wget https://storage.googleapis.com/golang/go1.8.linux-386.tar.gz`   
`sudo tar -xzf go1.8.linux-386.tar.gz -C /usr/local`      
设定环境变量   
`export GOPATH="/root/go"`   
`export GOROOT="/usr/local/go"`   
`export PATH=$GOROOT/bin:$PATH`   
之后我们在命令行中敲入  
`go version`   
`echo $GOPATH`   
`echo $GOROOT`   
`echo $PATH`   
检查是否配置完成   

2. Clone源代码   
`cd $GOPATH`   
`mkdir src`   
`cd src`   
`git clone https://github.com/trakel-project/travision`   
`cd travision`   

3. 配置依赖包
项目所需要的依赖包都包含在`vendor`目录中，我们所需要做的就是将依赖包复制到`GOPATH`中   
`cd vendor`   
`cp -rf . $GOPATH/src`   

4. 部署合约进行编译   
首先部署我们的合约。项目提供了一键部署合约的功能。   
`cd $GOPATH/src/travision/onekey`   
`contract`目录中存放了我们的合约代码   
修改`oneky.go`中`main`函数第四行代码`rpc.NewRpc()`函数的第一个参数为运行链的服务器IP   
`go run onekey.go`   
根据程序的输出来判断运行是否成功  
运行成功之后会在当前目录下生成`const.txt`文件，里面记录了`carcoin`和`taxi`合约的地址和ABI
利用新生成的地址和ABI还有链的IP替换`travision/gocode/contract/const.go`中的参数

5. 运行监控平台   
`cd $GOPATH/src/travision`   
`go build`   
`./travision`   
即可打开服务，开始监听8080端口。浏览器[打开](http://localhost:8080)即可使用监控平台。   
