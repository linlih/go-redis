package reply

import (
	"bytes"
	"go-redis/interface/resp"
	"strconv"
)

/*
	构建 redis 常用的回复信息函数
    其中 redis 定义了一套序列化协议，包含五种：
    1. 正常回复 :以 + 号开头，\r\n 结尾
	2. 错误回复 :以 - 号开头，\r\n 结尾
    3. 整数    :以 : 开头，\r\n 结尾
    4. 字符串:以 $ 开头，\r\n 结尾
	5. 数组    :以 * 开头，\r\n 结尾
	其中注意到：5 是多行的内容
*/

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type BulkReply struct {
	Arg []byte
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

type MultiBulkReply struct {
	Args [][]byte
}

func (m *MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	buf := bytes.Buffer{}
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// MakeMultiBulkReply 生成的是一个数组
func MakeMultiBulkReply(arg [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: arg,
	}
}

type StatusReply struct {
	status string
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		status: status,
	}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.status + CRLF)
}

type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

type StandErrReply struct {
	Status string
}

func MakeStandErrReply(status string) *StandErrReply {
	return &StandErrReply{Status: status}
}

func (s *StandErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func (s *StandErrReply) Error() string {
	return s.Status
}

func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}

type ErrReply struct {
	Msg string
}

func MakeErrReply(msg string) *ErrReply {
	return &ErrReply{
		Msg: msg,
	}
}

func (e *ErrReply) ToBytes() []byte {
	return []byte("-" + e.Msg + CRLF)
}
