package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Space-Software-LTDA/owncast/core/chat"
	"github.com/Space-Software-LTDA/owncast/models"
	webutils "github.com/Space-Software-LTDA/owncast/webserver/utils"
)

// GetConnectedChatClients returns currently connected clients.
func GetConnectedChatClients(w http.ResponseWriter, r *http.Request) {
	clients := chat.GetClients()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(clients); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

// ExternalGetConnectedChatClients returns currently connected clients.
func ExternalGetConnectedChatClients(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	GetConnectedChatClients(w, r)
}
