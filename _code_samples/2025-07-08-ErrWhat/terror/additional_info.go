package terror

import (
	"slices"
)

type AdditionalInfo interface {
	GetKey() (key string)
	GetValue() (value any)
}

type AdditionalInfos []AdditionalInfo

// Flatten removes duplicate keys from the AdditionalInfos slice.
// It keeps the last occurrence of each key based on the order in the slice.
func (instance AdditionalInfos) Flatten() (flattened AdditionalInfos) {
	firstIndexOfKey := make(map[string]int)
	lastIndexOfKey := make(map[string]int)
	for i, info := range instance {
		key := info.GetKey()
		lastIndexOfKey[key] = i
		_, exists := firstIndexOfKey[key]
		if !exists {
			firstIndexOfKey[key] = i
		}
	}

	type tuple struct {
		index int
		value AdditionalInfo
	}
	tuples := make([]tuple, 0, len(firstIndexOfKey))
	for key, firstIndex := range firstIndexOfKey {
		tuples = append(tuples, tuple{
			index: firstIndex,
			value: instance[lastIndexOfKey[key]],
		})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return a.index - b.index
	})

	flattened = make(AdditionalInfos, 0, len(instance))
	for _, tuple := range tuples {
		flattened = append(flattened, tuple.value)
	}

	return flattened
}

// ToJSON converts a slice of AdditionalInfo into a map where each info's key maps to its value.
// If multiple info entries have the same key, the last one's value will be used in the resulting map.
// Returns a map[string]any containing all key-value pairs from the AdditionalInfos slice.
func (instance AdditionalInfos) ToJSON() (flattened map[string]any) {
	flattened = make(map[string]any)
	for _, info := range instance {
		flattened[info.GetKey()] = info.GetValue()
	}

	return flattened
}

type jsonValue interface {
	~bool | ~int64 | ~uint64 | ~float64 | ~string
}

type typedInfo[T jsonValue] struct {
	key   string
	value T
}

func (instance typedInfo[T]) GetKey() (key string) {
	return instance.key
}

func (instance typedInfo[T]) GetValue() (value any) {
	return instance.value
}

type typedInfoSlice[T jsonValue] struct {
	key   string
	value []T
}

func (instance typedInfoSlice[T]) GetKey() (key string) {
	return instance.key
}

func (instance typedInfoSlice[T]) GetValue() (value any) {
	return instance.value
}

func WithBoolInfo[T ~bool](key string, value T) (info typedInfo[T]) {
	return typedInfo[T]{key: key, value: value}
}

func WithBoolSliceInfo[T ~bool](key string, value []T) (info typedInfoSlice[T]) {
	return typedInfoSlice[T]{key: key, value: value}
}

type intTypes interface {
	~int64 | ~int32 | ~int16 | ~int8 | ~int
}

func WithIntInfo[T intTypes](key string, value T) (info typedInfo[int64]) {
	return typedInfo[int64]{key: key, value: int64(value)}
}

func WithIntSliceInfo[T intTypes](key string, value []T) (info typedInfoSlice[int64]) {
	castValues := make([]int64, len(value))
	for i := range value {
		castValues[i] = int64(value[i])
	}

	return typedInfoSlice[int64]{key: key, value: castValues}
}

type uintTypes interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint
}

func WithUintInfo[T uintTypes](key string, value T) (info typedInfo[uint64]) {
	return typedInfo[uint64]{key: key, value: uint64(value)}
}

func WithUintSliceInfo[T uintTypes](key string, value []T) (info typedInfoSlice[uint64]) {
	castValues := make([]uint64, len(value))
	for i := range value {
		castValues[i] = uint64(value[i])
	}

	return typedInfoSlice[uint64]{key: key, value: castValues}
}

type floatTypes interface {
	~float64 | ~float32
}

func WithFloatInfo[T floatTypes](key string, value T) (info typedInfo[float64]) {
	return typedInfo[float64]{key: key, value: float64(value)}
}

func WithFloatSliceInfo[T floatTypes](key string, value []T) (info typedInfoSlice[float64]) {
	castValues := make([]float64, len(value))
	for i := range value {
		castValues[i] = float64(value[i])
	}

	return typedInfoSlice[float64]{key: key, value: castValues}
}

func WithStringInfo[T ~string](key string, value T) (info typedInfo[T]) {
	return typedInfo[T]{key: key, value: value}
}

func WithStringSliceInfo[T ~string](key string, value []T) (info typedInfoSlice[T]) {
	return typedInfoSlice[T]{key: key, value: value}
}
