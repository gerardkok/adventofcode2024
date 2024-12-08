package conv

import "strconv"

func MustAtoi(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}
