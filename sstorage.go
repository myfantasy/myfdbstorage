package myfdbstorage

import "errors"

// SStorageGenerates funcs for generate storage
var SStorageGenerates map[string]func(params map[string]interface{}) (SStorage, error)

// SStorageCreate create instance
func SStorageCreate(t string, params map[string]interface{}) (SStorage, error) {
	f, ok := SStorageGenerates[t]
	if ok {
		s, err := f(params)
		return s, err
	}

	return nil, errors.New("storage type not exists")

}

// SStorage - string kv storage
type SStorage interface {
	// Exists item
	Exists(key string) (ok bool, err error)
	// Get item
	Get(key string) (itm SItem, ok bool, err error)
	// Set - Insert or update item
	Set(itm SItem) (err error)
	// Del - delete item
	Del(key string) (exists bool, err error)

	// Type Storage for restore
	Type() string
	// Params Storage for restore
	Params() map[string]interface{}
	// Flush changes into storages
	Flush() error
	// ClearAndDeleteStorage table clear and then delete from storage
	ClearAndDeleteStorage() error
}

// SItem - string kv item
type SItem interface {
	// ID - uniqe object identity
	ID() string
	// Data - object data
	Data() []byte

	// FullData - all data without ID neaded for restore object
	FullData() []byte
}
