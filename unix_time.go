package caching

import (
	"encoding/json"
	"strconv"
	"time"
)

// UnixTime 毫秒级时间戳。
type UnixTime time.Time

var (
	_ json.Marshaler   = UnixTime{}
	_ json.Unmarshaler = (*UnixTime)(nil)
)

// MarshalJSON implements json.Marshaler.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	millisecond := time.Time(t).UnixNano() / int64(time.Millisecond)
	return []byte(strconv.FormatInt(millisecond, 10)), nil
}

// UnmarshalJSON implements json.UnmarshalJSON.
func (t *UnixTime) UnmarshalJSON(bytes []byte) error {
	v, err := strconv.ParseInt(string(bytes), 0, 64)
	if err != nil {
		return err
	}
	temp := time.Unix(0, v*int64(time.Millisecond))
	*t = (UnixTime)(temp)
	return nil
}

// String implements fmt.Stringer.
func (t UnixTime) String() string {
	millisecond := time.Time(t).UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(millisecond, 10)
}
