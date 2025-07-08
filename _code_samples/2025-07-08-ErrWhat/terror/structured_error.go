package terror

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
)

type StructuredContext interface {
	GetAdditionalInfo() (additionalInfos AdditionalInfos)
}

type StructuredError struct {
	context        StructuredContext
	cause          error
	callstack      []string
	additionalInfo AdditionalInfos
}

func New(ctx StructuredContext, cause error, additionalInfo ...AdditionalInfo) (err *StructuredError) {
	var callstack []string
	switch cause.(type) {
	case *StructuredError:
		// Don't generate the callstack multiple times
	default:
		callstack = make([]string, 0)
		stackDepth := 1
		for stackDepth < 1000 {
			_, file, line, success := runtime.Caller(stackDepth)
			if !success {
				break
			}
			stackDepth++
			callstack = append(callstack, fmt.Sprintf("%s: %d", file, line))
		}
	}

	return &StructuredError{
		context:        ctx,
		cause:          cause,
		callstack:      callstack,
		additionalInfo: additionalInfo,
	}
}

func (instance *StructuredError) Error() (message string) {
	return instance.cause.Error()
}

func (instance *StructuredError) Is(other error) (is bool) {
	switch e := instance.cause.(type) {
	case *StructuredError:
		return e.Is(other)
	default:
		return errors.Is(instance.cause, other)
	}
}

func (instance *StructuredError) As(other any) (ok bool) {
	switch e := instance.cause.(type) {
	case *StructuredError:
		return e.As(other)
	default:
		return errors.As(instance.cause, other)
	}
}

func (instance *StructuredError) getCallstack() (callstack string) {
	switch e := instance.cause.(type) {
	case *StructuredError:
		return e.getCallstack()
	default:
		return strings.Join(instance.callstack, "\n")
	}
}

func (instance *StructuredError) getAdditionalInfo(visitedContexts map[StructuredContext]bool) (additionalInfo AdditionalInfos) {
	if visitedContexts == nil {
		visitedContexts = make(map[StructuredContext]bool)
	}

	additionalInfo = AdditionalInfos{}

	// If the context is not nil, we can extract the additional info from it
	if instance.context != nil {
		_, found := visitedContexts[instance.context]
		if !found {
			additionalInfo = append(instance.context.GetAdditionalInfo(), additionalInfo...)
			visitedContexts[instance.context] = true
		}
	}

	// This error itself will have additional info, this should overwrite anything
	// from the context, i.e. come later in the array
	additionalInfo = append(additionalInfo, instance.additionalInfo...)

	// Finally recurse the errors, as each one may have more context and additional info
	// and deeper in the chain is more specific, so they come after the current info
	var childAdditionalInfo AdditionalInfos
	switch e := instance.cause.(type) {
	case *StructuredError:
		childAdditionalInfo = e.getAdditionalInfo(visitedContexts)
	default:
		// Wrapped another type of error, don't traverse further
	}
	additionalInfo = append(additionalInfo, childAdditionalInfo...)

	return additionalInfo
}

// When logging to otel we will want each part separately, so we can transform to their types etc
func GetLoggingInfo(err error) (cause string, callstack string, additionalInfo AdditionalInfos) {
	switch e := err.(type) {
	case *StructuredError:
		cause = e.Error()
		callstack = e.getCallstack()
		additionalInfo = e.getAdditionalInfo(nil).Flatten()
	default:
		cause = e.Error()
		callstack = "unwrapped error - no callstack"
		additionalInfo = make(AdditionalInfos, 0)
	}

	return cause, callstack, additionalInfo
}

func PrintError(err error) (errString string) {
	switch e := err.(type) {
	case *StructuredError:
		callstack := e.getCallstack()
		additionalInfo := e.getAdditionalInfo(nil).ToJSON()

		errString = fmt.Sprintf("Cause: %s\n", e.Error())

		// Sort additional info keys
		sortedKeys := make([]string, 0, len(additionalInfo))
		for key := range additionalInfo {
			sortedKeys = append(sortedKeys, key)
		}

		if len(sortedKeys) > 0 {
			slices.Sort(sortedKeys)
			errString += "Additional Info:\n"
		}

		for _, key := range sortedKeys {
			value := additionalInfo[key]
			errString += fmt.Sprintf("\t%s: %v\n", key, value)
		}

		errString += fmt.Sprintf("Callstack:\n\t%s\n", strings.ReplaceAll(callstack, "\n", "\n\t"))

	default:
		errString = fmt.Sprintf("Unknown error: %s\n", e.Error())
	}

	return errString
}
