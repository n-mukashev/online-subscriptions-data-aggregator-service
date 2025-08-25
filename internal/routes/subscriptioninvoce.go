package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/helpers"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/models"
)

// getSubscriptionsInvoice calculates the invoice for subscriptions
// @Summary Get subscriptions invoice
// @Description Calculate the total cost of subscriptions for a given user and period
// @Tags Subscription
// @Accept json
// @Produce json
// @Param request body models.SubscriptionInvoiceRequest true "Invoice Request"
// @Success 200 {object} map[string]int32
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/invoice [post]
func getSubscriptionsInvoice(ctx *gin.Context) {
	var request models.SubscriptionInvoiceRequest
	if !helpers.BindJSONWithValidation(ctx, &request) {
		return
	}
	invoice, err := request.GetSubscriptionsInvoice()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch subscriptions invoice"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sum": invoice})
}
