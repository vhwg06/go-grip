package webapi

import (
	"context"
	"fmt"
)

type AdminNotifier interface {
	SendBroadcast(ctx context.Context, title, body string) error
	SendTargeted(ctx context.Context, userID, title, body string) error
}

type noopAdminNotifier struct{}

func NewNoopAdminNotifier() AdminNotifier {
	return &noopAdminNotifier{}
}

func (n *noopAdminNotifier) SendBroadcast(_ context.Context, _ string, _ string) error { return nil }

func (n *noopAdminNotifier) SendTargeted(_ context.Context, _ string, _ string, _ string) error {
	return nil
}

type prefixedAdminNotifier struct {
	prefix string
}

func NewTelegramAdminNotifier() AdminNotifier { return &prefixedAdminNotifier{prefix: "telegram"} }
func NewBarkAdminNotifier() AdminNotifier     { return &prefixedAdminNotifier{prefix: "bark"} }
func NewEmailAdminNotifier() AdminNotifier    { return &prefixedAdminNotifier{prefix: "email"} }

func (n *prefixedAdminNotifier) SendBroadcast(ctx context.Context, title, body string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%s notifier broadcast: %w", n.prefix, err)
	}
	_ = title
	_ = body
	return nil
}

func (n *prefixedAdminNotifier) SendTargeted(ctx context.Context, userID, title, body string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%s notifier targeted: %w", n.prefix, err)
	}
	_ = userID
	_ = title
	_ = body
	return nil
}
