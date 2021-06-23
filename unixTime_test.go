package cache

import (
	"encoding/json"
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
	m := make(map[int]int, 10)

	for i := 0; i < 100; i++ {
		m[i] = i
	}
}
