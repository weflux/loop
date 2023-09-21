package loop

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/weflux/loop/option"
	"log/slog"
	"testing"
)

func Test_NewHub(t *testing.T) {

	hub := NewHub(option.Options{}, &slog.Logger{})
	err := hub.Start(context.Background())
	require.NoError(t, err)

}
