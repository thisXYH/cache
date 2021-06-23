package cache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestUnixTime(t *testing.T) {
	now := time.Now()
	millisecond := now.UnixNano() / int64(time.Millisecond)
	unix := UnixTime(now)
	t.Log(millisecond)

	b, err := json.Marshal(unix)
	if err != nil {
		t.Error(err)
	}

	var temp UnixTime
	json.Unmarshal(b, &temp)

	if temp.String() != unix.String() || temp.String() != strconv.Itoa(int(millisecond)) {
		t.Fail()
	}
}

func TestXX(t *testing.T) {
	data := int64(1)
	var i interface{} = data

	var data2 int64
	var i2 interface{} = &data2

	reflect.ValueOf(i2).Elem().Set(reflect.ValueOf(i))

	fmt.Println(data2)
}
