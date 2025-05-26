package tsafe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	assert.NotPanics(t, func() {
		Go(func() {
			panic("test panic")
		})
	})
}

func TestGoWithRecover(t *testing.T) {
	var waitChan = make(chan struct{})
	assert.NotPanics(t, func() {
		GoWithRecover(func() {
			panic("test panic")
		}, func(err any) {
			assert.Equal(t, "test panic", err)
			waitChan <- struct{}{}
		})
	})
	<-waitChan
}
