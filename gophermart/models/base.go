package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Money struct {
	Value float32
}

func (m *Money) Scan(src interface{}) (err error) {
    var value float64
    switch src := src.(type) {
	case string: value, err = strconv.ParseFloat(strings.TrimLeft(src, "$"), 32)
	case float64: value = src
	case nil: value = 0
	default: return fmt.Errorf("Money not allow Scan type=%T", src)
    }
    *m = Money{float32(value)}
    return err
}
