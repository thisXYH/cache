package cache_test

import (
	"encoding/json"
	"github.com/thisXYH/cache"
	"strconv"
	"testing"
	"time"
)

func TestUnixTime(t *testing.T) {
	now := time.Now()
	millisecond := now.UnixNano() / int64(time.Millisecond)
	unix := cache.UnixTime(now)
	t.Log(millisecond)

	b, err := json.Marshal(unix)
	if err != nil {
		t.Error(err)
	}

	var temp cache.UnixTime
	json.Unmarshal(b, &temp)

	if temp.String() != unix.String() || temp.String() != strconv.Itoa(int(millisecond)) {
		t.Fail()
	}
}
