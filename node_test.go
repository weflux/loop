package loopin

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/weflux/loopin/option"
	"log/slog"
	"testing"
)

func Test_NewHub(t *testing.T) {

	node := NewNode(option.Options{}, &slog.Logger{})
	err := node.Start(context.Background())
	require.NoError(t, err)

}
