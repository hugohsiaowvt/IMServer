package redis

import (
	"rz/util"
	"rz/config"
	rd "github.com/gomodule/redigo/redis"
)

type RedisConn struct {
	client rd.Conn
}

var instance *RedisConn

func InitRedis() {
	var err error
	if (config.IS_TEST) {
		instance.client, err = rd.Dial("tcp", config.TEST_REDIS_HOST, rd.DialPassword(config.TEST_REDIS_PASSWORD))
	} else {
		instance.client, err = rd.Dial("tcp", config.REDIS_HOST, rd.DialPassword(config.REDIS_PASSWORD))
	}

	util.CheckErr(err)
}

func Instance() *RedisConn {
	if instance == nil {
		instance = &RedisConn{}
		InitRedis()
	}
	return instance
}

func (rc *RedisConn) Keys(key string) (interface{}, error) {
	return rc.client.Do("KEYS", key)
}

func (rc *RedisConn) Set(key string, value interface{}) (interface{}, error) {
	return rc.client.Do("SET", key, value)
}

func (rc *RedisConn) SetAndExpire(key string, value interface{}, time int) (interface{}, error) {
	return rc.client.Do("SET", key, value, "ex", time)
}

func (rc *RedisConn) Get(key string) (string, error) {
	return rd.String(rc.client.Do("GET", key))
}

func (rc *RedisConn) Del(key string) (interface{}, error) {
	return rc.client.Do("DEL", key)
}

func (rc *RedisConn) Expire(key string, time int) (interface{}, error) {
	return rc.client.Do("EXPIRE", key, time)
}

func (rc *RedisConn) HSet(key, field string, value interface{}) (interface{}, error) {
	return rc.client.Do("HSET", key, field, value)
}

func (rc *RedisConn) HGet(key, field interface{}) (string, error) {
	return rd.String(rc.client.Do("HGET", key, field))
}

func (rc *RedisConn) LPush(key string, value string) (interface{}, error) {
	return rc.client.Do("LPUSH" ,key, value)
}

func (rc *RedisConn) LRange(key string, min, max interface{}) ([]string, error) {
	return rd.Strings(rc.client.Do("LRange" ,key, min, max))
}

func (rc *RedisConn) SAdd(key string, value string) (interface{}, error) {
	return rc.client.Do("SADD" ,key, value)
}

func (rc *RedisConn) SMembers(key string) ([]string, error) {
	return rd.Strings(rc.client.Do("SMEMBERS" ,key))
}