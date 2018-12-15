package helpers

func IsIntInArray(value int, array []int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

type Tuple struct {
	Left, Right interface{}
}
