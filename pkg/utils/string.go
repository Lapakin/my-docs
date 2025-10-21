package utils

import (
	"strconv"
)

func UInt64ToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
