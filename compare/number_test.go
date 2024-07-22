package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxInt(t *testing.T) {
	assert.Equal(t, Max(10, 20), 20)
	assert.Equal(t, Max(-10, 20), 20)
	assert.Equal(t, Max(-10, -20), -10)
}
func TestMaxInt8(t *testing.T) {
	assert.Equal(t, Max(int8(10), int8(20)), int8(20))
	assert.Equal(t, Max(int8(-10), int8(20)), int8(20))
	assert.Equal(t, Max(int8(-10), int8(-20)), int8(-10))
}
func TestMaxInt16(t *testing.T) {
	assert.Equal(t, Max(int16(10), int16(20)), int16(20))
	assert.Equal(t, Max(int16(-10), int16(20)), int16(20))
	assert.Equal(t, Max(int16(-10), int16(-20)), int16(-10))
}
func TestMaxInt32(t *testing.T) {
	assert.Equal(t, Max(int32(10), int32(20)), int32(20))
	assert.Equal(t, Max(int32(-10), int32(20)), int32(20))
	assert.Equal(t, Max(int32(-10), int32(-20)), int32(-10))
}
func TestMaxInt64(t *testing.T) {
	assert.Equal(t, Max(int64(10), int64(20)), int64(20))
	assert.Equal(t, Max(int64(-10), int64(20)), int64(20))
	assert.Equal(t, Max(int64(-10), int64(-20)), int64(-10))
}
func TestMaxUint(t *testing.T) {
	assert.Equal(t, Max(uint(10), uint(20)), uint(20))
	assert.Equal(t, Max(uint(10), uint(0)), uint(10))
}
func TestMaxUint8(t *testing.T) {
	assert.Equal(t, Max(uint8(10), uint8(20)), uint8(20))
	assert.Equal(t, Max(uint8(10), uint8(0)), uint8(10))
}
func TestMaxUint16(t *testing.T) {
	assert.Equal(t, Max(uint16(10), uint16(20)), uint16(20))
	assert.Equal(t, Max(uint16(10), uint16(0)), uint16(10))
}
func TestMaxUint32(t *testing.T) {
	assert.Equal(t, Max(uint32(10), uint32(20)), uint32(20))
	assert.Equal(t, Max(uint32(10), uint32(0)), uint32(10))
}
func TestMaxUint64(t *testing.T) {
	assert.Equal(t, Max(uint64(10), uint64(20)), uint64(20))
	assert.Equal(t, Max(uint64(10), uint64(0)), uint64(10))
}

func TestMin(t *testing.T) {
	for _, tc := range []struct {
		a, b int
		want int
	}{
		{5, 3, 3},
		{1, 5, 1},
		{5, 5, 5},
		{-5, 3, -5},
		{-5, -3, -5},
		{0, -5, -5},
	} {
		t.Run("", func(t *testing.T) {
			got := Min(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Min(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMaxSlice(t *testing.T) {
	slice := []int{1, 5, 3, 7, 4}
	maxSlice, err := MaxSlice(slice...)
	assert.Nil(t, err)
	assert.Equal(t, 7, maxSlice)

	// 测试空切片的情况
	var emptySlice []int
	_, err = MaxSlice(emptySlice...)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, EmptySliceError)
}

func TestMinSlice(t *testing.T) {
	t.Run("non-empty slice with min value", func(t *testing.T) {
		min, err := MinSlice[int](-1, 0, 1)
		assert.Nil(t, err)
		assert.Equal(t, -1, min)
	})
	t.Run("non-empty slice with different types", func(t *testing.T) {
		min, err := MinSlice[int32](int32(100), int32(3), int32(10), int32(0), int32(-10))
		assert.Nil(t, err)
		assert.Equal(t, int32(-10), min)
	})
	t.Run("empty slice", func(t *testing.T) {
		_, err := MinSlice[int]()
		assert.ErrorIs(t, err, EmptySliceError)
	})
	t.Run("custom error check", func(t *testing.T) {
		_, err := MinSlice[int]()
		var emptySliceErr error
		if len(err.Error()) > 0 {
			emptySliceErr = EmptySliceError
		}
		assert.Equal(t, emptySliceErr, err)
	})
}
