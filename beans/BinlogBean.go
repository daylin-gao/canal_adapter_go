package beans

import(
	protocol "github.com/gao111/canal-adapter-go/protocol"
)

type BinlogBean struct {
	Columns     []*protocol.Column // 更新内容
	EventType   protocol.EventType // 事件类型
	SchemaName  string  // 库名
	TableName   string  // 表明
}