package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"pos-saas/internal/domain"
)

// Monetary limits match NUMERIC(10,2). Preserve the existing 10% tax and zero fee.
const maximumOrderCents int64 = 9_999_999_999

func pricedOrderItem(product domain.Product, request domain.CreateOrderItemRequest) (domain.OrderItem, error) {
	if product.Status != "active" {
		return domain.OrderItem{}, errors.New("product is not available")
	}
	if request.Quantity <= 0 {
		return domain.OrderItem{}, errors.New("quantity must be positive")
	}
	if product.TrackInventory && product.QuantityInStock < request.Quantity {
		return domain.OrderItem{}, errors.New("insufficient inventory")
	}
	price := product.Price
	if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
		price = *product.DiscountPrice
	}
	if math.IsNaN(price) || math.IsInf(price, 0) || price < 0 || price > float64(maximumOrderCents)/100 {
		return domain.OrderItem{}, errors.New("invalid product price")
	}
	cents := int64(math.Round(price * 100))
	if cents > 0 && int64(request.Quantity) > maximumOrderCents/cents {
		return domain.OrderItem{}, errors.New("order amount exceeds supported limit")
	}
	return domain.OrderItem{ProductID: request.ProductID, VariantID: request.VariantID, ProductName: product.NameEn, Quantity: request.Quantity, UnitPrice: float64(cents) / 100, TotalPrice: float64(cents*int64(request.Quantity)) / 100, SpecialInstructions: request.SpecialInstructions}, nil
}

func applyOrderPricing(order *domain.Order) error {
	var subtotal int64
	for _, item := range order.Items {
		subtotal += int64(math.Round(item.TotalPrice * 100))
		if subtotal > maximumOrderCents {
			return errors.New("order amount exceeds supported limit")
		}
	}
	tax := (subtotal + 5) / 10 // Round the established 10% tax to cents, as PostgreSQL does.
	if subtotal+tax > maximumOrderCents {
		return errors.New("order amount exceeds supported limit")
	}
	order.Subtotal = float64(subtotal) / 100
	order.TaxAmount = float64(tax) / 100
	order.DeliveryFee = 0
	order.TotalAmount = float64(subtotal+tax)/100 - order.DiscountAmount
	if math.IsNaN(order.TotalAmount) || math.IsInf(order.TotalAmount, 0) || order.TotalAmount < 0 {
		return errors.New("order total cannot be negative or invalid")
	}
	return nil
}

func (uc *OrderUseCase) QuoteOrder(ctx context.Context, tenant, restaurant int64, request *domain.OrderQuoteRequest) (*domain.OrderQuote, error) {
	if request == nil || len(request.Items) == 0 || len(request.Items) > 100 {
		return nil, fmt.Errorf("%w: quote requires 1 to 100 items", domain.ErrInvalidOrderQuote)
	}
	ids := make([]int64, 0, len(request.Items))
	quantities := make(map[int64]int)
	for _, item := range request.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 || item.Quantity > 10000 {
			return nil, fmt.Errorf("%w: invalid quote item", domain.ErrInvalidOrderQuote)
		}
		ids = append(ids, item.ProductID)
		quantities[item.ProductID] += item.Quantity
	}
	products, err := uc.productRepo.GetOrderQuoteProducts(ctx, tenant, restaurant, ids)
	if err != nil {
		return nil, err
	}
	order := &domain.Order{}
	for _, item := range request.Items {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("%w: product is unavailable", domain.ErrInvalidOrderQuote)
		}
		if !product.IsAvailable || product.TrackInventory && product.QuantityInStock < quantities[item.ProductID] {
			return nil, fmt.Errorf("%w: product is unavailable or insufficient inventory", domain.ErrInvalidOrderQuote)
		}
		priced, err := pricedOrderItem(product, item)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInvalidOrderQuote, err)
		}
		order.Items = append(order.Items, priced)
	}
	if err := applyOrderPricing(order); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidOrderQuote, err)
	}
	lines := make([]domain.OrderQuoteItem, 0, len(order.Items))
	for _, item := range order.Items {
		lines = append(lines, domain.OrderQuoteItem{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: item.UnitPrice, TotalPrice: item.TotalPrice})
	}
	return &domain.OrderQuote{Items: lines, Subtotal: order.Subtotal, TaxAmount: order.TaxAmount, DeliveryFee: order.DeliveryFee, TotalAmount: order.TotalAmount, Currency: "EGP"}, nil
}
