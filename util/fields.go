package util

import (
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
)

type ExtraFields map[string]string

func MergeExtraFields(extraFields map[string]string, newFields map[string]string) ExtraFields {
	if extraFields == nil {
		return newFields
	}
	maps.Copy(extraFields, newFields)

	return extraFields
}

func (fields ExtraFields) ToAttrs() []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, attribute.Key(key).String(value))
	}
	return attrs
}

func (extra ExtraFields) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for key, value := range extra {
		encoder.AddString(key, value)
	}
	return nil
}
