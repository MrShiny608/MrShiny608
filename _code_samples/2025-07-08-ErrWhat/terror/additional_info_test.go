package terror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct{}
	type result struct {
		flattened AdditionalInfos
	}
	type testConfig struct {
		name          string
		instance      *AdditionalInfos
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns an empty map when no additional info is provided",
			instance: &AdditionalInfos{},
			args:     &args{},
			result: &result{
				flattened: AdditionalInfos{},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "returns keys from multiple additional info of different types",
			instance: &AdditionalInfos{
				WithStringInfo("key1", "value1"),
				WithIntInfo("key2", 1),
			},
			args: &args{},
			result: &result{
				flattened: AdditionalInfos{
					WithStringInfo("key1", "value1"),
					WithIntInfo("key2", 1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "returns the last key's value when multiple keys have the same value",
			instance: &AdditionalInfos{
				WithStringInfo("key1", "value1"),
				WithIntInfo("key1", 1),
			},
			args: &args{},
			result: &result{
				flattened: AdditionalInfos{
					WithIntInfo("key1", 1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "flattens taking the last instance of any given key, but retaining the order of the first instance",
			instance: &AdditionalInfos{
				WithIntInfo("one", 3),
				WithIntInfo("two", 2),
				WithIntInfo("one", 1),
			},
			args: &args{},
			result: &result{
				flattened: AdditionalInfos{
					WithIntInfo("one", 1),
					WithIntInfo("two", 2),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			// Arrange
			instance := config.instance
			args := config.args
			result := config.result
			arrangeFunc, actFunc, assertFunc := config.generateHooks(t, args, result)

			arrangeFunc(t)

			// Act
			actFunc(t)
			flattened := (AdditionalInfos).Flatten(*instance)

			// Assert
			assert.Equal(t, result.flattened, flattened)

			assertFunc(t)
		})
	}
}

func TestToJSON(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct{}
	type result struct {
		flattened map[string]any
	}
	type testConfig struct {
		name          string
		instance      *AdditionalInfos
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns an empty map when no additional info is provided",
			instance: &AdditionalInfos{},
			args:     &args{},
			result: &result{
				flattened: map[string]any{},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "returns keys from multiple additional info of different types",
			instance: &AdditionalInfos{
				WithStringInfo("key1", "value1"),
				WithIntInfo("key2", 1),
			},
			args: &args{},
			result: &result{
				flattened: map[string]any{
					"key1": "value1",
					"key2": int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "returns the last key's value when multiple keys have the same value",
			instance: &AdditionalInfos{
				WithStringInfo("key1", "value1"),
				WithIntInfo("key1", 1),
			},
			args: &args{},
			result: &result{
				flattened: map[string]any{
					"key1": int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "flattens into a map",
			instance: &AdditionalInfos{
				WithIntInfo("one", 1),
				WithIntInfo("two", 2),
			},
			args: &args{},
			result: &result{
				flattened: map[string]any{
					"one": int64(1),
					"two": int64(2),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "flattens taking the last instance of any given key",
			instance: &AdditionalInfos{
				WithIntInfo("one", 3),
				WithIntInfo("two", 2),
				WithIntInfo("one", 1),
			},
			args: &args{},
			result: &result{
				flattened: map[string]any{
					"one": int64(1),
					"two": int64(2),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			// Arrange
			instance := config.instance
			args := config.args
			result := config.result
			arrangeFunc, actFunc, assertFunc := config.generateHooks(t, args, result)

			arrangeFunc(t)

			// Act
			actFunc(t)
			flattened := (AdditionalInfos).ToJSON(*instance)

			// Assert
			assert.Equal(t, result.flattened, flattened)

			assertFunc(t)
		})
	}
}

func TestTypedInfo(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T) (instance AdditionalInfo)
	type assertFunc func(t *testing.T)
	type args struct{}
	type result struct {
		info AdditionalInfo
	}
	type testConfig struct {
		name          string
		instance      *interface{}
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "with bool info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[bool]{
					key:   "one",
					value: true,
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithBoolInfo("one", true)
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with bool slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[bool]{
					key:   "one",
					value: []bool{true},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithBoolSliceInfo("one", []bool{true})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[int64]{
					key:   "one",
					value: 1,
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntInfo("one", int(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[int64]{
					key:   "one",
					value: []int64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntSliceInfo("one", []int{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int8 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[int64]{
					key:   "one",
					value: int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntInfo("one", int8(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int8 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[int64]{
					key:   "one",
					value: []int64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntSliceInfo("one", []int8{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int16 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[int64]{
					key:   "one",
					value: int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntInfo("one", int16(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int16 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[int64]{
					key:   "one",
					value: []int64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntSliceInfo("one", []int16{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int16 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[int64]{
					key:   "one",
					value: int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntInfo("one", int16(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int32 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[int64]{
					key:   "one",
					value: []int64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntSliceInfo("one", []int32{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int64 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[int64]{
					key:   "one",
					value: int64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntInfo("one", int64(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with int64 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[int64]{
					key:   "one",
					value: []int64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithIntSliceInfo("one", []int64{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[uint64]{
					key:   "one",
					value: 1,
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintInfo("one", uint(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[uint64]{
					key:   "one",
					value: []uint64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintSliceInfo("one", []uint{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint8 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[uint64]{
					key:   "one",
					value: uint64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintInfo("one", uint8(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint8 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[uint64]{
					key:   "one",
					value: []uint64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintSliceInfo("one", []uint8{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint16 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[uint64]{
					key:   "one",
					value: uint64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintInfo("one", uint16(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint16 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[uint64]{
					key:   "one",
					value: []uint64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintSliceInfo("one", []uint16{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint16 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[uint64]{
					key:   "one",
					value: uint64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintInfo("one", uint16(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint32 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[uint64]{
					key:   "one",
					value: []uint64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintSliceInfo("one", []uint32{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint64 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[uint64]{
					key:   "one",
					value: uint64(1),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintInfo("one", uint64(1))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with uint64 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[uint64]{
					key:   "one",
					value: []uint64{1},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithUintSliceInfo("one", []uint64{1})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with float32 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[float64]{
					key:   "one",
					value: float64(1.0),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithFloatInfo("one", float32(1.0))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with float32 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[float64]{
					key:   "one",
					value: []float64{1.0},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithFloatSliceInfo("one", []float32{1.0})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with float64 info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[float64]{
					key:   "one",
					value: 1.0,
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithFloatInfo("one", float64(1.0))
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with float64 slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[float64]{
					key:   "one",
					value: []float64{1.0},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithFloatSliceInfo("one", []float64{1.0})
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with string info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfo[string]{
					key:   "one",
					value: "string",
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithStringInfo("one", "string")
				}, func(t *testing.T) {}
			},
		},
		{
			name:     "with string slice info",
			instance: nil,
			args:     &args{},
			result: &result{
				info: typedInfoSlice[string]{
					key:   "one",
					value: []string{"string"},
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) (instance AdditionalInfo) {
					return WithStringSliceInfo("one", []string{"string"})
				}, func(t *testing.T) {}
			},
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			// Arrange
			args := config.args
			result := config.result
			arrangeFunc, actFunc, assertFunc := config.generateHooks(t, args, result)

			arrangeFunc(t)

			// Act
			instance := actFunc(t)

			// Assert
			assert.Equal(t, result.info, instance)

			assertFunc(t)
		})
	}
}
