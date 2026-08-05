package notification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationDispatchDisabled(t *testing.T) {
	t.Parallel()
	require.NoError(t, NewNotificationUseCase(false).Dispatch(context.Background(), Notification{To: "a@example.com"}))
}
