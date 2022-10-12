# go-redis
implement redis by go

参考：

# 编码流程
## Part Ⅰ 项目的基本内容
配置文件：`config/config.go`，用于配置的管理，其中用到了反射的方式来实现文件配置和结构体的配置对应设置的功能
 
反射的设置的过程：
> 1. 首先使用 bufio.NewScanner 按照行读入，根据空格区分 key 和 value，存在到 rawMap 中
> 2. 创建反射对象：Type 和 Value，依次遍历所有的结构体成员，根据结构体成员 tag 中的 key 值从 rawMap 中取出值，根据反射 Type 的不同类型进行设置

文件处理函数：`lib/logger/files.go` 主要是封装了一些文件的判断函数，比如判断文件是否存在、是否有权限等

日志输出函数：`lib/logger/logger.go` 对日志输出函数进行封装，包括设置输出文件、日志前缀、日志等级等

bool原子操作：`lib/sync/atomic/bool.go` 因为在 go 标准库中没有提供 bool 的原子操作，所以这里做了一个封装

wait操作封装：`lib/sync/wait/wait.go` 这个和 go 标准库中的 wait 是差不多的，多提供了一个超时退出的等待机制

## Part Ⅱ TCP 服务器的实现

处理函数接口：`interface/tcp/handler.go`，主要是处理连接函数，以及关闭函数

TCP 服务器：`tcp/server.go` 实现了 TCP 的连接的监听和分发，这里的 `ListenAndServeWithSignal` 主要是注册了系统中几个信号的处理函数，到收到这些信号的时候，主要对一些资源进行关闭处理。
在 `ListenAndServe` 函数中则是死循环处理来自客户端的连接，然后分发给 handler 的 `Handle` 函数

Echo 处理Handle：`tcp/echo.go` 这个是实现一个回复Handler的示例函数，用于基本功能测试。比较重要的核心点就是关闭的处理，需要关闭所有的连接，所以在 Handler 中保存了所有的活跃连接。

主函数：`main.go` ：进行配置读入操作，调用 tcp 的监听和 handler 函数。

## Part Ⅲ Redis 的通信协议 RESP 实现

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

## Part Ⅳ 实现 Redis 内存数据库

`command.go`：定义了全局的命令处理函数集合 cmdTable，是一个 map 类型，同时注意，这个 map 并不需要保证并发写安全，所以只需要用一个普通的map就可以，因为所有的map中的内容都会在程序启动的初始化过程中，插入完成。普通的map并发读是没有问题的。

`dict.go`：定义 Dict 类型的接口，可以扩展底层的实现

`sync_dict.go`：封装 go 提供的 sync.Map 作为 Dict 的一个具体实现

`db.go`：redis中是支持多个 DB 的，所以这里是一个具体的 DB 实现，核心函数是 Exec 函数，该函数从 cmdTable 取出具体的执行函数。其他的实现则是存取、删除数据的操作函数，但是要注意的是，这里的数据进行了接口的封装是：DataEntity，而不是具体的数据类型，比如string，int之类的。

`keys.go`： 实现 Keys 操作的一些命令，比如 DEL、KEYS、FLUSHDB、TYPE、RENAME 等

`string.go`：实现字符串操作的命令，比如 SET、GET、SETNX 等

`wildcard.go`：实现一个简单的正则表示式匹配，核心步骤分为 Compile 和 Match，使用动态规划来实现

`database.go`：数据库中实现了 Database 的接口，需要特殊的处理的是选择哪一个 DB 进行操作，然后调用具体的 db 的 Exec 函数就可以了。

可以使用网络调试助手工具，也可以直接使用 redis-cli 工具

测试命令：

设置 key 和 value：`*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n`

获取 key 中设置的值：`*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n`

## Part Ⅴ 实现 AOF

AOF 的全称是 AppendOnlyFile，用于 Redis 的持久化。

核心是实现了 `aof.go` 这么一个文件，实现了打开一个外部文件，写入 AOF 内容，和加载 AOF 内容的功能。其中加载 `LoadAof` 函数复用了之前的解析函数 `ParseStream`。

打开该功能，后台会启动一个 `handleAof` 的程序监听 `aofChan` 管道得到需要写入的数据信息，其中需要对是否发生 DB 的切换进行特殊处理，如果发生了 DB 的切换，则需要多写入切换 DB 的命令到 AOF 文件中。

其中还做了一个小的转换，在 db.go 中的 `addAof` 函数是 `func(CmdLine)` ， 而在 `aof.go` 中的 `AddAof` 函数是 `func(dbIndex int, cmd CmdLine)`，所以在 `database.go` 文件中对数据库进行初始化的时候进行了转换操作。

这样做的好处是，将对象不关心的数据进行剥离，让每一个对象实体只看到自己需要的内容。比如对于 `db` 来说，我只需要知道怎么写入 AOF 就可以了，我当前的数据库 index 是多少我并不知关心，而且我也无从得知，所以交给上层的 database 来管理数据库的 index 就可以了。













