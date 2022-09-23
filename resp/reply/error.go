package reply

type UnKnowErrReply struct {
}

var unkownErrBytes = []byte("-Err unknown\r\n")

func (u UnKnowErrReply) Error() string {
	return "Err unknown"
}

func (u UnKnowErrReply) ToBytes() []byte {
	return unkownErrBytes
}

type ArgNumErrReply struct {
	Cmd string
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{cmd}
}

func (a ArgNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + a.Cmd + "' command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

type SyntaxErrReply struct {
}

func (s SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (w WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (w WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

type ProtocolErrReply struct {
	Msg string
}

func (p ProtocolErrReply) Error() string {
	return "-ERR Protocol error: '" + p.Msg
}

func (p ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + p.Msg + "'\r\n")
}
