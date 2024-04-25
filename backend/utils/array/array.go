package array_operation

func InArray[T comparable](element T, array []T) bool {
	for _, value := range array {
		if value == element {
			return true
		}
	}
	return false
}

func Remove[T comparable](element T, array []T) []T {
	for i, v := range array {
		if v == element {
			// 通过切片操作删除元素
			array = append(array[:i], array[i+1:]...)
			break
		}
	}
	return array
}
