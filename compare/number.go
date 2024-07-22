package compare

import "errors"

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

var EmptySliceError = errors.New("empty slice")

// Max 是一个泛型函数，用于比较两个整数或浮点数类型参数 a 和 b 的大小，并返回较大的值。
//
// 参数：
// a, b: 类型为 T 的参数，可以是 int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32 或 float64 中的任意一种。
//
// 返回值：
// 类型为 T 的值，表示 a 和 b 中较大的一个。
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 是一个泛型函数，用于比较两个整数或浮点数类型值的大小，并返回较小的值。
//
// 参数：
// a, b：要比较的两个整数或浮点数类型值，类型必须相同且为 int、int8、int16、int32、int64、uint、uint8、uint16、uint32、uint64、float32 或 float64 之一。
//
// 返回值：
// 返回 a 和 b 中较小的值，类型为 T。
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// MaxSlice 是一个泛型函数，用于查找整数或浮点数切片中的最大值
// T 可以是 int、int8、int16、int32、int64、uint、uint8、uint16、uint32、uint64、float32 或 float64 类型的整数或浮点数
// slice 是要查找最大值的切片，类型为 T 的可变参数
// 返回值：
// - 第一个返回值是切片中的最大值，类型为 T
// - 第二个返回值是一个 error，如果切片为空则返回 EmptySliceError，否则为 nil
func MaxSlice[T Number](slice ...T) (T, error) {
	if len(slice) == 0 {
		return 0, EmptySliceError
	}
	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if res < slice[i] {
			res = slice[i]
		}
	}
	return res, nil
}

// MinSlice 是一个泛型函数，用于找到整数切片中的最小值
// T 可以是 int, int8, int16, int32, int64, uint, uint8, uint16, uint32 或 uint64 类型的整数
// slice 是要查找最小值的切片，可变参数
// 函数返回切片中的最小值以及一个 error，如果切片为空，则返回 EmptySliceError
func MinSlice[T Number](slice ...T) (T, error) {
	if len(slice) == 0 {
		return 0, EmptySliceError
	}
	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if res > slice[i] {
			res = slice[i]
		}
	}

	return res, nil
}
