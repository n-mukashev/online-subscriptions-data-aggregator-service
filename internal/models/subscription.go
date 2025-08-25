package models

import (
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/logger"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/storage"
)

// Subscription represents a user's subscription
// @Description A subscription that a user has to a service
type Subscription struct {
	// ID of the subscription
	// example: 1
	Id int64 `json:"id"`
	// Name of the service
	// example: "Netflix"
	ServiceName string `json:"service_name" binding:"required"`
	// Monthly price of the subscription
	// example: 100
	MonthlyPrice int32 `json:"monthly_price" binding:"required"`
	// ID of the user who owns the subscription
	// example: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	UserId uuid.UUID `json:"user_id" binding:"required"`
	// Start date of the subscription
	// example: "06-2006"
	StartDate MonthYear `json:"start_date" binding:"required"`
	// End date of the subscription
	// example: "06-2006"
	EndDate *MonthYear `json:"end_date"`
}

// UpdateSubscription represents data for updating a subscription
// @Description Fields for updating an existing subscription
type UpdateSubscription struct {
	// ID of the subscription to update
	// example: 1
	Id int64 `json:"id" binding:"required"`
	// Name of the service
	// example: "Netflix"
	ServiceName string `json:"service_name"`
	// Monthly price of the subscription
	// example: 100
	MonthlyPrice int32 `json:"monthly_price"`
	// ID of the user
	// example: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	UserId uuid.UUID `json:"user_id"`
	// Start date of the subscription
	// example: "06-2006"
	StartDate *MonthYear `json:"start_date"`
	// End date of the subscription
	// example: "06-2006"
	EndDate *MonthYear `json:"end_date"`
}

func GetById(id int64) (*Subscription, error) {
	row := storage.DB.QueryRow(
		`SELECT id, service_name, monthly_price, user_id, start_date, end_date
		 FROM subscription WHERE id = $1`, id)

	var s Subscription
	err := row.Scan(&s.Id, &s.ServiceName, &s.MonthlyPrice, &s.UserId, &s.StartDate, &s.EndDate)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		logger.Log.Error("failed to get subscription by id", slog.Any("err", err))
		return nil, err
	}
	return &s, nil
}

func GetAll() ([]Subscription, error) {
	rows, err := storage.DB.Query(
		`SELECT id, service_name, monthly_price, user_id, start_date, end_date FROM subscription`)
	if err != nil {
		logger.Log.Error("failed to get subscriptions", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var s Subscription
		err := rows.Scan(&s.Id, &s.ServiceName, &s.MonthlyPrice, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			logger.Log.Error("failed to scan subscription row", slog.Any("err", err))
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	return subscriptions, nil
}

func (s *Subscription) Create() error {
	query := `
		INSERT INTO subscription (id, service_name, monthly_price, user_id, start_date, end_date)
		VALUES (nextval('subscription_seq'), $1, $2, $3, $4, $5) RETURNING id`
	err := storage.DB.QueryRow(query,
		s.ServiceName, s.MonthlyPrice, s.UserId, s.StartDate.ToTime(), s.EndDate.ToTime()).
		Scan(&s.Id)
	if err != nil {
		logger.Log.Error("failed to create subscription", slog.Any("err", err))
		return err
	}
	return nil
}

func (req *UpdateSubscription) Update() error {
	s, err := GetById(req.Id)
	if err != nil {
		return err
	}
	s.compareAndUpdate(req)
	query := `
	UPDATE subscription 
	SET service_name = $1, monthly_price = $2, user_id = $3, start_date = $4, end_date = $5 
	WHERE id = $6`
	stmt, err := storage.DB.Prepare(query)
	if err != nil {
		logger.Log.Error("failed to prepare update statement", slog.Any("err", err))
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(s.ServiceName, s.MonthlyPrice, s.UserId, s.StartDate.ToTime(), s.EndDate.ToTime(), s.Id)
	if err != nil {
		logger.Log.Error("failed to execute update", slog.Any("err", err))
		return err
	}
	logger.Log.Info("updated subscription", slog.Any("id", s.Id))
	return nil
}

func (to *Subscription) compareAndUpdate(from *UpdateSubscription) {
	if from.ServiceName != "" {
		to.ServiceName = from.ServiceName
	}
	if from.MonthlyPrice != 0 {
		to.MonthlyPrice = from.MonthlyPrice
	}
	if from.UserId != uuid.Nil {
		to.UserId = from.UserId
	}
	if from.StartDate != nil {
		to.StartDate = *from.StartDate
	}
	if from.EndDate != nil {
		to.EndDate = from.EndDate
	}
}

func Delete(id int64) error {
	query := `DELETE FROM subscription WHERE id = $1`
	stmt, err := storage.DB.Prepare(query)
	if err != nil {
		logger.Log.Error("failed to prepare delete statement", slog.Any("id", id), slog.Any("err", err))
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		logger.Log.Error("failed to execute delete", slog.Any("id", id), slog.Any("err", err))
		return err
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error("failed to get rows affected for delete", slog.Any("id", id), slog.Any("err", err))
		return err
	}
	if deleted == 0 {
		logger.Log.Warn("no record deleted", slog.Any("id", id))
		return sql.ErrNoRows
	}
	logger.Log.Info("deleted subscription", slog.Any("id", id))
	return nil
}
