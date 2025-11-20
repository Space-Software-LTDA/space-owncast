package events

import "github.com/Space-Software-LTDA/owncast/models"

// ConnectedClientInfo represents the information about a connected client.
type ConnectedClientInfo struct {
	User *models.User `json:"user"`
	Event
}
