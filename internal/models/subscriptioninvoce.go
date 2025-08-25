package models

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/logger"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/storage"
)

// SubscriptionInvoiceRequest represents a request to calculate the total cost of subscriptions.
// swagger:model SubscriptionInvoiceRequest
type SubscriptionInvoiceRequest struct {
	// Service name
	// example: "Netflix"
	ServiceName string `json:"service_name" binding:"required"`

	// Unique user identifier
	// example: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	UserId uuid.UUID `json:"user_id" binding:"required"`

	// Start date of the period (month/year)
	// example: "06-2006"
	FromDate MonthYear `json:"from_date" binding:"required"`

	// End date of the period (month/year). Can be null, in which case the current date is used
	// example: "06-2006"
	ToDate *MonthYear `json:"to_date" binding:"required"`
}

func (f *SubscriptionInvoiceRequest) GetSubscriptionsInvoice() (int32, error) {
	query := `
	SELECT 
		COALESCE(
			SUM(
				(DATE_PART('year', age(
					LEAST(end_date, $2), 
					GREATEST(start_date, $1)
				)) * 12 
				+ DATE_PART('month', age(
					LEAST(end_date, $2), 
					GREATEST(start_date, $1)
				))
				) * monthly_price
			), 0
		) AS total_cost
	FROM subscription
	WHERE service_name = $3
	AND user_id = $4
	AND start_date <= $2
	AND end_date >= $1;
	`
	var invoice int32
	err := storage.DB.QueryRow(query, f.FromDate.ToTime(), f.ToDate.ToTime(), f.ServiceName, f.UserId).Scan(&invoice)
	if err != nil {
		logger.Log.Error("failed to fetch subscriptions invoice", slog.Any("err", err))
		return 0, err
	}
	return invoice, nil
}
