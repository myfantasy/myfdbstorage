package myfdbstorage

func init() {
	IStorageGenerates = make(map[string]func(params map[string]interface{}) (IStorage, error))

	IStorageGenerates[IStorageSimpleName] = func(params map[string]interface{}) (IStorage, error) {
		s, e := IStorageSimpleLoad(params)
		return s, e
	}

	SStorageGenerates = make(map[string]func(params map[string]interface{}) (SStorage, error))

	SStorageGenerates[SStorageSimpleName] = func(params map[string]interface{}) (SStorage, error) {
		s, e := SStorageSimpleLoad(params)
		return s, e
	}

}

// LogFunc - if func is not null Write Log
var LogFunc func(error)

// Exception write exception to log
func Exception(err error) {
	if LogFunc != nil {
		LogFunc(err)
	}
}
