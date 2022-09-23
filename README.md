# go-redis
implement redis by go

# 编码流程
### Part Ⅰ 项目的基本内容
配置文件：`config/config.go`，用于配置的管理，其中用到了反射的方式来实现文件配置和结构体的配置对应设置的功能
 
反射的设置的过程：
> 1. 首先使用 bufio.NewScanner 按照行读入，根据空格区分 key 和 value，存在到 rawMap 中
> 2. 创建反射对象：Type 和 Value，依次遍历所有的结构体成员，根据结构体成员 tag 中的 key 值从 rawMap 中取出值，根据反射 Type 的不同类型进行设置

文件处理函数：`lib/logger/files.go` 主要是封装了一些文件的判断函数，比如判断文件是否存在、是否有权限等

日志输出函数：`lib/logger/logger.go` 对日志输出函数进行封装，包括设置输出文件、日志前缀、日志等级等

bool原子操作：`lib/sync/atomic/bool.go` 因为在 go 标准库中没有提供 bool 的原子操作，所以这里做了一个封装

wait操作封装：`lib/sync/wait/wait.go` 这个和 go 标准库中的 wait 是差不多的，多提供了一个超时退出的等待机制

### Part Ⅱ TCP 服务器的实现

处理函数接口：`interface/tcp/handler.go`，主要是处理连接函数，以及关闭函数

TCP 服务器：`tcp/server.go` 实现了 TCP 的连接的监听和分发，这里的 `ListenAndServeWithSignal` 主要是注册了系统中几个信号的处理函数，到收到这些信号的时候，主要对一些资源进行关闭处理。
在 `ListenAndServe` 函数中则是死循环处理来自客户端的连接，然后分发给 handler 的 `Handle` 函数

Echo 处理Handle：`tcp/echo.go` 这个是实现一个回复Handler的示例函数，用于基本功能测试。比较重要的核心点就是关闭的处理，需要关闭所有的连接，所以在 Handler 中保存了所有的活跃连接。

主函数：`main.go` ：进行配置读入操作，调用 tcp 的监听和 handler 函数。

### Part Ⅲ Redis 的通信协议 RESP 实现

首先要明确 RESP 的协议格式，分为如下五种：
- 正常回复：以 `+` 号开头，`\r\n` 结尾的字符串形式
- 错误回复：以 `-` 号开头，`\r\n` 结尾的字符串形式
- 整数： 以 `:` 开头，以 `\r\n` 结尾的字符串形式
- 字符串：以 `$` 开头后跟实际发送字节数，以 `\r\n`结尾
- 数组：以 `*` 开头，后面跟成员个数

示例：
`*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n`

实现的要点是定义了 readState 这么一个状态：
```go
type readState struct {
	readingMultiline  bool // 是否多行读取模式
	expectedArgsCount int  // 需要多少个参数
	msgType           byte // 
	args              [][]byte
	bulkLen           int64
}
```
整个的调用流程如下：

`ListenAndServeWithSignal -> ListenAndServe -> RespHandler.Handle -> ParseStream -> parse0`

其中 `parse0` 死循环从 socket io 中将数据读出，进行解析。

我们来看一个例子的解析过程：
`*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n`

开始的时候会初始化一个 readState 变量，所有数值都设置为默认值。

首先调用的是 readLine 得到第一个 msg，得到的结果是：`*3\r\n`。

然后因为一开始的时候并不是多行读取模式，所以对 readLine 结果的 msg 进行判断。

此时的 msg 开头是 *，调用 parseMultiBulkHeader 来解析，设置为多行读取模式，同时将想要读取 Args 数量设置为 3。

进入下一个循环的 readLine 得到 msg 为 `$3\r\n`，走到 readBody 函数，此时函数的开头是 $ ，所以读取得到参数是 bulkLen 为 3。

继续进入下一个循环走到 readLine ，得到 msg 为 `SET\r\n`，这个时候将数据放到 readState 的 args 参数中。

然后依次循环读取，最终的结束条件是看当前读取的参数数量是否满足。

其他的流程可以按照根据代码进行分析。










