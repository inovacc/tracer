//go:build windows

package logger

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	cfgLogger := NewConfigModule().
		WithServiceName("testService").
		WithHostname("testHostname").
		WithLogFilePath("C:/arqprod_local/logs").
		WithLogMaxSize(5).
		WithLogMaxBackups(5).
		WithLogMaxAge(2).
		WithLogCompress(true).
		WithLevel("info").
		Build()

	if err := InitLogger(cfgLogger); err != nil {
		fmt.Printf("error initializing Logger: %v", err)
		return
	}

	m.Run()
}

func TestNewLoggerWithTraceID(t *testing.T) {
	ctx := context.Background()
	componentName := "test/componentName"
	spanName := "spanName"

	tracerLogger := NewLoggerWithTraceID(ctx, componentName, spanName)
	assert.NotNil(t, tracerLogger)

	defer tracerLogger.Close()

	tracerLogger.Info("sender", "format")
}
