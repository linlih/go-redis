package dict

type Consumer func(key string, val interface{}) bool

// Dict 定义了字典类型的接口，可以扩展字典的底层实现
type Dict interface {
	Get(key string) (val interface{}, exits bool)
	Len() int
	Put(key string, val interface{}) (result int)
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
