package terror

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ///////////////////////////////////////////////////////////////
// Otel context is a span managing context, in reality you'd have any additional information required
// for your application, e.g. user information, current request scope, etc.
type OtelContext struct {
	context.Context

	otelSpan       trace.Span
	traceID        trace.TraceID
	spanID         trace.SpanID
	parentID       trace.SpanID
	additionalInfo AdditionalInfos
	dispatched     bool
}

// Creates a new root context. If this is part of a distributed operation, i.e. another application began
// the trace, then the traceID and parentID should be set to the values from that operation, otherwise
// they can be nil and we generate IDs where appropriate.
func NewRootContext(name string, traceID *trace.TraceID, parentID *trace.SpanID, additionalInfo ...AdditionalInfo) (instance *OtelContext) {
	// If there is no traceID provided, we generate a new trace ID
	if traceID == nil {
		traceID = &trace.TraceID{}
		// copied from otel/trace.randomIDGenerator.NewIDs
		for {
			_, _ = rand.Read(traceID[:]) // #nosec G104 docs say this never returns an error
			if traceID.IsValid() {
				break
			}
		}
	}

	// If there is no parentID provided, we create an invalid one - otel doesn't like nil
	if parentID == nil {
		parentID = &trace.SpanID{}
	}

	// Set up the otel span
	tracer := otel.Tracer("default")

	ctx := context.Background()
	if !parentID.IsValid() {
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: *traceID,
		})
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
	} else {
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    *traceID,
			SpanID:     *parentID,          // For a remote span this is the parent ID
			Remote:     true,               // Recorded as a remote span
			TraceFlags: trace.FlagsSampled, // Sampled or it wont be recorded
		})
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
	}

	spanContext, otelSpan := tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithTimestamp(time.Now().UTC()),
	)

	spanID := otelSpan.SpanContext().SpanID()

	return &OtelContext{
		Context: spanContext,

		otelSpan:       otelSpan,
		spanID:         spanID,
		traceID:        *traceID,
		parentID:       *parentID,
		additionalInfo: additionalInfo,
		dispatched:     false,
	}
}

// A child context is a span within the current root, it will inherit the traceID and parentID from the current context.
// Really it's just a helper so you don't have to worrk about trace ID and parentID, but it also helps to document
// where a span enters the system (root) and where it is expected to be nested (child).
func (instance *OtelContext) NewChildContext(name string, additionalInfo ...AdditionalInfo) (childContext *OtelContext) {
	// Set up the otel span
	tracer := otel.Tracer("default")

	ctx := context.Background()
	if !instance.parentID.IsValid() {
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: instance.traceID,
		})
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
	} else {
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    instance.traceID,
			SpanID:     instance.parentID,  // For a remote span this is the parent ID
			Remote:     true,               // Recorded as a remote span
			TraceFlags: trace.FlagsSampled, // Sampled or it wont be recorded
		})
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
	}

	spanContext, otelSpan := tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithTimestamp(time.Now().UTC()),
	)

	spanID := otelSpan.SpanContext().SpanID()

	return &OtelContext{
		Context: spanContext,

		otelSpan:       otelSpan,
		spanID:         spanID,
		traceID:        instance.traceID,
		parentID:       instance.parentID,
		additionalInfo: additionalInfo,
		dispatched:     false,
	}
}

func (instance *OtelContext) Close(err *error) {
	if instance.dispatched {
		return
	}

	// Add the additional info to the span
	attributes := make([]attribute.KeyValue, len(instance.additionalInfo))
	for i, info := range instance.additionalInfo {
		switch value := info.GetValue().(type) {
		case bool:
			attributes[i] = attribute.Bool(info.GetKey(), value)
		case int64:
			attributes[i] = attribute.Int64(info.GetKey(), value)
		case uint64:
			// otel/attribute doesn't support uint64
			attributes[i] = attribute.Int64(info.GetKey(), int64(value))
		case float64:
			attributes[i] = attribute.Float64(info.GetKey(), value)
		case string:
			attributes[i] = attribute.String(info.GetKey(), value)
		case []bool:
			attributes[i] = attribute.BoolSlice(info.GetKey(), value)
		case []int64:
			attributes[i] = attribute.Int64Slice(info.GetKey(), value)
		case []uint64:
			// otel/attribute doesn't support uint64
			values := make([]int64, len(value))
			for j, v := range value {
				values[j] = int64(v)
			}
			attributes[i] = attribute.Int64Slice(info.GetKey(), values)
		case []float64:
			attributes[i] = attribute.Float64Slice(info.GetKey(), value)
		case []string:
			attributes[i] = attribute.StringSlice(info.GetKey(), value)
		}
	}

	instance.otelSpan.SetAttributes(attributes...)

	// Set the status
	var status codes.Code
	var cause string
	if *err != nil {
		status = codes.Error
		cause = (*err).Error()
	} else {
		status = codes.Ok
	}

	instance.otelSpan.SetStatus(status, cause)

	instance.otelSpan.End(
		trace.WithTimestamp(time.Now().UTC()),
	)

	instance.dispatched = true
}

func (instance *OtelContext) GetAdditionalInfo() (additionalInfo AdditionalInfos) {
	return instance.additionalInfo
}

// ///////////////////////////////////////////////////////////////
// Error types for testing purposees
var ErrSomethingBad = errors.New("something bad happened")

type BadError struct{}

func (e *BadError) Error() string {
	return ErrSomethingBad.Error()
}

func (e *BadError) Is(target error) bool {
	return target == ErrSomethingBad
}

// ///////////////////////////////////////////////////////////////
// F1 is the start of our chain of execution, it creates a root context
func F1() (err error) {
	// If we were continuing a distributed trace, we would pass the traceID and parentID from the incoming
	// request, instead of nils
	ctx := NewRootContext("F1", nil, nil, WithStringInfo("key1", "value1"))
	defer ctx.Close(&err)

	// Emulate work
	time.Sleep(100 * time.Millisecond)

	err = F2(ctx)
	if err != nil {
		return New(ctx, err, WithStringInfo("error1", "value1"))
	}

	return nil
}

// F2 is still within our system, so it creates a child context to denote the nested span
func F2(ctx *OtelContext) (err error) {
	ctx = ctx.NewChildContext("F2", WithStringInfo("key2", "value2"))
	defer ctx.Close(&err)

	// Emulate work
	time.Sleep(100 * time.Millisecond)

	err = F3(ctx)
	if err != nil {
		return New(ctx, err, WithStringInfo("error2", "value2"))
	}

	return nil
}

// F3 is an external call, it doesn't understand OtelContext, and will return a normal error
func F3(ctx context.Context) (err error) {
	return &BadError{}
}

func TestExample(t *testing.T) {
	t.Parallel()

	err := F1()

	var badError *BadError
	assert.ErrorAs(t, err, &badError)
	assert.ErrorIs(t, err, ErrSomethingBad)
	assert.Equal(t, err.Error(), ErrSomethingBad.Error())

	cause, callstack, additionalInfo := GetLoggingInfo(err)
	assert.Equal(t, cause, ErrSomethingBad.Error())
	assert.Contains(t, callstack, "example_test.go:")
	assert.Equal(t, additionalInfo, AdditionalInfos{
		WithStringInfo("key1", "value1"),
		WithStringInfo("error1", "value1"),
		WithStringInfo("key2", "value2"),
		WithStringInfo("error2", "value2"),
	})
}
