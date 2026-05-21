package content

import "context"

func (uc *UseCase) CatchUpScheduled(ctx context.Context) (int, error) {
	return uc.PublishDue(ctx)
}
