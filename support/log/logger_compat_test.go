package log

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// legacyLogger mirrors the Logger test doubles kept in other modules; it stops compiling if a method is added.
type legacyLogger struct{}

func (legacyLogger) DebugEnabled() bool                   { return false }
func (legacyLogger) TraceEnabled() bool                   { return false }
func (legacyLogger) Trace(...interface{})                 {}
func (legacyLogger) Debug(...interface{})                 {}
func (legacyLogger) Info(...interface{})                  {}
func (legacyLogger) Warn(...interface{})                  {}
func (legacyLogger) Error(...interface{})                 {}
func (legacyLogger) Tracef(string, ...interface{})        {}
func (legacyLogger) Debugf(string, ...interface{})        {}
func (legacyLogger) Infof(string, ...interface{})         {}
func (legacyLogger) Warnf(string, ...interface{})         {}
func (legacyLogger) Errorf(string, ...interface{})        {}
func (legacyLogger) Structured() StructuredLogger         { return nil }
func (legacyLogger) GetTracingContext() map[string]string { return nil }
func (legacyLogger) SetTracingContext(map[string]string)  {}

var _ Logger = legacyLogger{}

func TestLoggerInterfaceUnchanged(t *testing.T) {
	want := "Debug DebugEnabled Debugf Error Errorf GetTracingContext Info Infof SetTracingContext Structured Trace " +
		"TraceEnabled Tracef Warn Warnf"
	it := reflect.TypeOf((*Logger)(nil)).Elem()
	var got []string
	for i := 0; i < it.NumMethod(); i++ {
		got = append(got, it.Method(i).Name)
	}
	sort.Strings(got)
	if strings.Join(got, " ") != want {
		t.Errorf("Logger methods = %v, want %s", got, want)
	}
}
