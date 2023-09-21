package loop

import (
	"context"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func Test_NewHub(t *testing.T) {

	hub := NewHub(Options{}, &slog.Logger{})
	err := hub.Start(context.Background())
	require.NoError(t, err)

}
