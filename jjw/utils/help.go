package utils

import (
	"bytes"
	"encoding/gob"
	"time"
	"strconv"
	"math/rand"
	"encoding/json"
)

const fmt = "2006-01-02 15:04:05"

// 公有成员 深拷贝
func DeepCopy(src, target interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(target)
}

// 秒时间转换
func SecFormat(timeUnix int64) string {
	return time.Unix(timeUnix,0).Format(fmt)
}
// 毫秒时间转换
func MilliSecFormat(timeUnix int64) string {
	return time.Unix(timeUnix/1e3,0).Format(fmt)
}

// now 毫秒
func Now() int64 {
	return time.Now().UnixNano()/1e6
}

//带时序随机数
func RandomId() string {
	return strconv.FormatInt( time.Now().UnixNano()/1e6 , 10) + strconv.Itoa(rand.Intn(10000))
}

//
func Sleep(milliSec int) {
	time.Sleep(time.Duration(milliSec) * time.Millisecond)
}

// json编码
func JsonMarshal(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		return "---------解析json出错---------------"
	}
	return string(j)
}
