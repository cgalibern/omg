package statusbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"opensvc.com/opensvc/core/kind"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/core/status"
)

func TestRefuseStartTwice(t *testing.T) {
	bus := T{}
	defer bus.Stop()
	bus.Start()
	assert.PanicsWithError(t, ErrorStarted.Error(), bus.Start)
}

func TestPanicIfNotStarted(t *testing.T) {
	bus := T{}
	defer bus.Stop()
	t.Run("Post", func(t *testing.T) {
		assert.PanicsWithError(
			t,
			ErrorNeedStart.Error(),
			func() {
				bus.Post(path.T{}, "app#1", status.Warn, false)
			})
	})
	t.Run("Post", func(t *testing.T) {
		assert.PanicsWithError(
			t,
			ErrorNeedStart.Error(),
			func() {
				bus.Get(path.T{}, "app#1")
			})
	})
}

func TestPost(t *testing.T) {
	bus := T{}
	bus.Start()
	defer bus.Stop()
	p := path.T{
		Name:      "foo",
		Namespace: "root",
		Kind:      kind.Svc,
	}
	bus.Post(p, "app#1", status.Down, false)
	bus.Post(p, "app#2", status.Up, false)

	cases := []struct {
		rid   string
		state status.T
	}{
		{"app#1", status.Down},
		{"app#2", status.Up},
		{"app#3", status.Undef},
	}
	for _, ridState := range cases {
		t.Logf("ensure rid %s status is %v", ridState.rid, ridState.state)
		found := bus.Get(p, ridState.rid)
		assert.Equal(t, ridState.state, found)
	}
	t.Run("status is undef when service is not found", func(t *testing.T) {
		assert.Equal(t, status.Undef, bus.Get(path.T{}, ""))
		assert.Equal(t, status.Undef, bus.Get(path.T{}, "app#1"))
	})
}
