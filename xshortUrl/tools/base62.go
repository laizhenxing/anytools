package tools

import (
	"errors"
	"math"
	"strings"
)

var (
	Base = 62
	CharacterSet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func Base62Encode(num int) string {
	r := make([]rune, 0)

	for num > 0 {
		n := math.Mod(float64(num), float64(Base))
		num /= Base
		// 查找对应n位置的字符，添加到 r 的头部
		r = append([]rune{CharacterSet[int(n)]}, r...)
	}
	return string(r)
}

func Base62Decode(s string) (int, error) {
	var r, pow int

	for i, v := range s {
		// 对应每个字符所在位置的幂，如2，即为62^2
		pow = len(s) - (i + 1)
		pos := strings.IndexRune(string(CharacterSet), v)
		if pos == -1 {
			return 0, errors.New("invalid character: " + string(v))
		}
		r += pos * int(math.Pow(float64(Base), float64(pow)))
	}

	return r, nil
}