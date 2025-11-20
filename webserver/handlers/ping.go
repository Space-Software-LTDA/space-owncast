package handlers

import (
	"net/http"

	"github.com/Space-Software-LTDA/owncast/core"
	"github.com/Space-Software-LTDA/owncast/models"
)

// Ping is fired by a client to show they are still an active viewer.
func Ping(w http.ResponseWriter, r *http.Request) {
	viewer := models.GenerateViewerFromRequest(r)
	core.SetViewerActive(&viewer)
	w.WriteHeader(http.StatusOK)
}
