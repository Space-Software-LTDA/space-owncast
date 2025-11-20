package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Space-Software-LTDA/owncast/core"
	"github.com/Space-Software-LTDA/owncast/metrics"
	"github.com/Space-Software-LTDA/owncast/models"
	"github.com/Space-Software-LTDA/owncast/persistence/configrepository"
	"github.com/Space-Software-LTDA/owncast/webserver/router/middleware"
	log "github.com/sirupsen/logrus"
)

// Status gets the details of the inbound broadcaster.
func Status(w http.ResponseWriter, r *http.Request) {
	configRepository := configrepository.Get()

	broadcaster := core.GetBroadcaster()
	status := core.GetStatus()
	currentBroadcast := core.GetCurrentBroadcast()
	health := metrics.GetStreamHealthOverview()
	response := adminStatusResponse{
		Broadcaster:            broadcaster,
		CurrentBroadcast:       currentBroadcast,
		Online:                 status.Online,
		Health:                 health,
		ViewerCount:            status.ViewerCount,
		OverallPeakViewerCount: status.OverallMaxViewerCount,
		SessionPeakViewerCount: status.SessionMaxViewerCount,
		VersionNumber:          status.VersionNumber,
		StreamTitle:            configRepository.GetStreamTitle(),
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type adminStatusResponse struct {
	Broadcaster            *models.Broadcaster          `json:"broadcaster"`
	CurrentBroadcast       *models.CurrentBroadcast     `json:"currentBroadcast"`
	Health                 *models.StreamHealthOverview `json:"health"`
	StreamTitle            string                       `json:"streamTitle"`
	VersionNumber          string                       `json:"versionNumber"`
	ViewerCount            int                          `json:"viewerCount"`
	OverallPeakViewerCount int                          `json:"overallPeakViewerCount"`
	SessionPeakViewerCount int                          `json:"sessionPeakViewerCount"`
	Online                 bool                         `json:"online"`
}
