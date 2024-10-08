package handlers

import (
	"strconv"
	fn "validator-api/pkgs/functions"
	"validator-api/responses"
	s "validator-api/server"
	"validator-api/services"

	"github.com/gofiber/fiber/v2"
)

type HandlerSyncDuties struct {
	Server *s.Server
}

func NewHandlerSyncDuties(server *s.Server) *HandlerSyncDuties {
	return &HandlerSyncDuties{
		Server: server,
	}
}

// @Summary Get sync duties
// @Description Get sync duties by slot
// @Tags syncduties
// @Accept json
// @Produce json
// @Param slot path integer true "Slot"
// @Success 200 {object} responses.SyncDuties
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /syncduties/{slot} [get]
func (h *HandlerSyncDuties) GetSyncDuties(ctx *fiber.Ctx) error {

	slot, _ := strconv.Atoi(ctx.Params("slot"))

	// Check if slot is in future
	if fn.IsSlotInFuture(slot) {
		return responses.ErrorResponse(ctx,
			fiber.StatusBadRequest,
			"Requested slot is in the future.",
		)
	}

	keys := []string{}

	err := services.SyncCommitteeIndices(slot, h.Server.Config.QuckNode.Http, &keys)

	if err != nil {
		return responses.ErrorResponse(ctx,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return responses.ResponseValidators(ctx, fiber.StatusOK, keys)
}
