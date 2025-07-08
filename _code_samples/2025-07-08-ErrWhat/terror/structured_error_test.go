package terror

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testContext struct {
	context.Context
	additionalInfo AdditionalInfos
}

func (instance *testContext) GetAdditionalInfo() (additionalInfos AdditionalInfos) {
	return instance.additionalInfo
}

func TestError(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
		message string
	}
	type testConfig struct {
		name          string
		instance      *StructuredError
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns the root error message",
			instance: New(nil, fmt.Errorf("root error")),
			args:     &args{},
			result: &result{
				message: "root error",
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "recursively calls until it returns the root error message",
			instance: New(nil, New(nil, fmt.Errorf("root error"))),
			args:     &args{},
			result: &result{
				message: "root error",
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
			message := (*StructuredError).Error(instance)

			// Assert
			assert.Equal(t, result.message, message)

			assertFunc(t)
		})
	}
}

func TestGetCallstack(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
		callstack string
	}
	type testConfig struct {
		name          string
		instance      *StructuredError
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns the callstack",
			instance: New(nil, fmt.Errorf("root error")),
			args:     &args{},
			result: &result{
				callstack: "",
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {
					// Find the current file path
					_, currentFilePath, _, ok := runtime.Caller(1)
					assert.True(t, ok)

					result.callstack = currentFilePath
				}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "doesn't generate the callstack multiple times",
			instance: New(nil, New(nil, fmt.Errorf("root error"))),
			args:     &args{},
			result: &result{
				callstack: "",
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
			callstack := (*StructuredError).getCallstack(instance)

			// Assert
			assert.Contains(t, callstack, result.callstack)

			assertFunc(t)
		})
	}
}

func TestGetAdditionalInfo(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
		additionalInfo AdditionalInfos
	}
	type testConfig struct {
		name          string
		instance      *StructuredError
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns the additional info",
			instance: New(nil, fmt.Errorf("root error"), WithStringInfo("key1", "value1")),
			args:     &args{},
			result: &result{
				additionalInfo: AdditionalInfos{
					WithStringInfo("key1", "value1"),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "recursively calls until it returns the additional info, appending the nested error's additional info",
			instance: New(nil, New(nil, fmt.Errorf("root error"), WithStringInfo("key1", "value1")), WithStringInfo("key2", "value2")),
			args:     &args{},
			result: &result{
				additionalInfo: AdditionalInfos{
					WithStringInfo("key2", "value2"),
					WithStringInfo("key1", "value1"),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "returns additional info from the context, prepending the context's additional info",
			instance: New(&testContext{
				additionalInfo: AdditionalInfos{WithStringInfo("key2", "value2")},
			}, fmt.Errorf("root error"), WithStringInfo("key1", "value1")),
			args: &args{},
			result: &result{
				additionalInfo: AdditionalInfos{
					WithStringInfo("key2", "value2"),
					WithStringInfo("key1", "value1"),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name: "recursively calls until it returns the additional info, appending the nested error's context's additional info if they have different roots",
			instance: New(
				&testContext{
					additionalInfo: AdditionalInfos{WithStringInfo("key4", "value4")},
				},
				New(&testContext{
					additionalInfo: AdditionalInfos{WithStringInfo("key2", "value2")},
				}, fmt.Errorf("root error"), WithStringInfo("key1", "value1")),
				WithStringInfo("key3", "value3"),
			),
			args: &args{},
			result: &result{
				additionalInfo: AdditionalInfos{
					WithStringInfo("key4", "value4"),
					WithStringInfo("key3", "value3"),
					WithStringInfo("key2", "value2"),
					WithStringInfo("key1", "value1"),
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
			additionalInfo := (*StructuredError).getAdditionalInfo(instance, nil)

			// Assert
			assert.Equal(t, result.additionalInfo, additionalInfo)

			assertFunc(t)
		})
	}
}

func TestGetLoggingInfo(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
		cause          string
		callstack      string
		additionalInfo AdditionalInfos
	}
	type testConfig struct {
		name          string
		instance      error
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "returns the info from a structured error",
			instance: New(nil, fmt.Errorf("root error"), WithStringInfo("key1", "value1")),
			args:     &args{},
			result: &result{
				cause:     "root error",
				callstack: "",
				additionalInfo: AdditionalInfos{
					WithStringInfo("key1", "value1"),
				},
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {
					// Find the current file path
					_, currentFilePath, _, ok := runtime.Caller(1)
					assert.True(t, ok)

					result.callstack = currentFilePath
				}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "returns the info from a standard go error",
			instance: fmt.Errorf("root error"),
			args:     &args{},
			result: &result{
				cause:          "root error",
				callstack:      "unwrapped error - no callstack",
				additionalInfo: AdditionalInfos{},
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
			cause, callstack, additionalInfo := GetLoggingInfo(instance)

			// Assert
			assert.Equal(t, result.cause, cause)
			assert.Contains(t, callstack, result.callstack)
			assert.Equal(t, result.additionalInfo, additionalInfo)

			assertFunc(t)
		})
	}
}

func TestPrintError(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
		errString string
	}
	type testConfig struct {
		name          string
		instance      error
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "formats a structured error",
			instance: New(nil, fmt.Errorf("root error"), WithStringInfo("key1", "value1")),
			args:     &args{},
			result: &result{
				errString: "Cause: root error\nAdditional Info:\n\tkey1: value1\nCallstack:\n\t",
			},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {
					// Find the current file path
					_, currentFilePath, _, ok := runtime.Caller(1)
					assert.True(t, ok)

					result.errString += currentFilePath
				}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "formats a standard go error",
			instance: fmt.Errorf("root error"),
			args:     &args{},
			result: &result{
				errString: "Unknown error: root error\n",
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
			errString := PrintError(instance)

			// Assert
			assert.Contains(t, errString, result.errString)

			assertFunc(t)
		})
	}
}

func TestSupportsErrorsPackage(t *testing.T) {
	t.Parallel()

	type arrangeFunc func(t *testing.T)
	type actFunc func(t *testing.T)
	type assertFunc func(t *testing.T)
	type args struct {
	}
	type result struct {
	}
	type testConfig struct {
		name          string
		instance      error
		args          *args
		result        *result
		generateHooks func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc)
	}
	configs := []testConfig{
		{
			name:     "works with a structured error",
			instance: New(nil, errors.New("root error"), WithStringInfo("key1", "value1")),
			args:     &args{},
			result:   &result{},
			generateHooks: func(t *testing.T, args *args, result *result) (arrangeFunc arrangeFunc, actFunc actFunc, assertFunc assertFunc) {
				return func(t *testing.T) {}, func(t *testing.T) {}, func(t *testing.T) {}
			},
		},
		{
			name:     "works when wrapped with errors.Join",
			instance: errors.Join(New(nil, errors.New("root error"), WithStringInfo("key1", "value1"))),
			args:     &args{},
			result:   &result{},
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
			var e *StructuredError
			success := errors.As(instance, &e)
			is := errors.Is(instance, e)
			unwrapped := errors.Unwrap(instance)

			// Assert
			assert.True(t, success)
			assert.True(t, is)
			assert.Nil(t, unwrapped) // We don't allow unwrapping

			assertFunc(t)
		})
	}
}
