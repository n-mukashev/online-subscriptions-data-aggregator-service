package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/helpers"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/logger"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/models"
)

// @Summary Get subscription by ID
// @Description Get a single subscription by its ID
// @Tags Subscription
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/{id} [get]
func getById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You need to provide an ID"})
		return
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error("Could not parse id", slog.Any("err", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse id"})
		return
	}
	subscription, err := models.GetById(id)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Subscription not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch subscription"})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

// @Summary Get all subscriptions
// @Description Get a list of all subscriptions
// @Tags Subscription
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscription/all [get]
func getAll(ctx *gin.Context) {
	subscriptions, err := models.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch all subscriptions"})
		return
	}
	ctx.JSON(http.StatusOK, subscriptions)
}

// @Summary Create new subscription
// @Description Create a new subscription entry
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription [post]
func create(ctx *gin.Context) {
	var subscription models.Subscription
	if !helpers.BindJSONWithValidation(ctx, &subscription) {
		return
	}
	err := subscription.Create()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create the new subscription"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": subscription.Id, "message": "New subscription was created"})
}

// @Summary Update subscription
// @Description Update an existing subscription by ID
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription body models.UpdateSubscription true "Subscription data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription [put]
func update(ctx *gin.Context) {
	var subscription models.UpdateSubscription
	if !helpers.BindJSONWithValidation(ctx, &subscription) {
		return
	}
	err := subscription.Update()
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Subscription was not found with given ID"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update the subscription"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "The subscription was successfully updated"})
}

// @Summary Delete subscription
// @Description Delete a subscription by its ID
// @Tags Subscription
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/{id} [delete]
func delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You need to provide an ID"})
		return
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error("Could not parse id", slog.Any("err", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse int"})
		return
	}
	err = models.Delete(id)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Subscription not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the subscription"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "The subscription was successfully deleted"})
}
