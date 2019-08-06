package myfdbstorage

func init() {
	IStorageGenerates = make(map[string]func(params map[string]interface{}) (IStorage, error))

	IStorageGenerates[IStorageSimpleName] = func(params map[string]interface{}) (IStorage, error) {
		s := &IStorageSimple{Data: make(map[int64]IItem)}

		return s, nil
	}

}
