package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"volume-ts/pkg/errors"
	"volume-ts/pkg/middleware"
)

func (h Handler) SortFlights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := middleware.GetUserId()
	if err != nil {
		errors.JSONError(w, errors.Error{
			Status:  false,
			Message: "check auth details",
			Data:    nil,
		}, http.StatusUnauthorized)
	}

	sortedPath, err := h.flightService.SortFlightPaths(ctx, userId)
	log.WithError(err).Errorf("failed to sort flight paths")
	if err != nil {
		errors.JSONError(w, errors.Error{
			Status:  false,
			Message: "failed to sort flight paths",
			Data:    nil,
		}, http.StatusBadRequest)
	}

	_ = json.NewEncoder(w).Encode(sortedPath)
}
