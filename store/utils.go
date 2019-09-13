package store


func Get(key string)(interface{},error){
	return rdb.Get(key).Result()
}

func Set(key string, value interface{}) error{
	return rdb.Set(key,value,2000).Err()
}

func Del(key string) error{
	return  rdb.Del(key).Err()
}