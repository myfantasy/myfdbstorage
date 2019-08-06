package myfdbstorage

import "errors"

// IStorageGenerates funcs for generate storage
var IStorageGenerates map[string]func(params map[string]interface{}) (IStorage, error)

// IStorageCreate create instance
func IStorageCreate(t string, params map[string]interface{}) (IStorage, error) {
	f, ok := IStorageGenerates[t]
	if ok {
		s, err := f(params)
		return s, err
	}

	return nil, errors.New("storage type not exists")

}

// IStorage - int64 kv storage
type IStorage interface {
	// Exists item
	Exists(key int64) (ok bool, err error)
	// Get item
	Get(key int64) (itm IItem, ok bool, err error)
	// Set - Insert or update item
	Set(itm IItem) (err error)
	// Del - delete item
	Del(key int64) (exists bool, err error)

	// Type Storage for restore
	Type() string
	// Params Storage for restore
	Params() map[string]interface{}
	// ClearAndDeleteStorage table clear and then delete from storage
	ClearAndDeleteStorage() error
}

// IItem - int64 kv item
type IItem interface {
	// ID - uniqe object identity
	ID() int64
	// Data - object data
	Data() []byte
}
