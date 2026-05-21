package notification

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestNotificationDispatchDisabled(t *testing.T) {
	t.Parallel()
	require.NoError(t, New(false).Dispatch(context.Background(), entity.Notification{To: "a@example.com"}))
}
