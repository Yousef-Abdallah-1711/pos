package usecase

import (
	"context"
	"encoding/json"

	"pos-saas/internal/domain"
)

func (uc *OrderUseCase) CancelOrderItem(ctx context.Context, request domain.OrderItemCancellation) (json.RawMessage, error) {
	return uc.orderRepo.CancelOrderItem(ctx, request)
}
