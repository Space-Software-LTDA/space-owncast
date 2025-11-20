package models

import (
	"time"

	"github.com/Space-Software-LTDA/owncast/utils"
)

const (
	moderatorScopeKey = "MODERATOR"
)

type User struct {
	CreatedAt       time.Time  `json:"createdAt"`
	DisabledAt      *time.Time `json:"disabledAt,omitempty"`
	NameChangedAt   *time.Time `json:"nameChangedAt,omitempty"`
	AuthenticatedAt *time.Time `json:"-"`
	ID              string     `json:"id"`
	DisplayName     string     `json:"displayName"`
	Email           string     `json:"email"`
	PreviousNames   []string   `json:"previousNames"`
	Scopes          []string   `json:"scopes,omitempty"`
	DisplayColor    int        `json:"displayColor"`
	IsBot           bool       `json:"isBot"`
	Authenticated   bool       `json:"authenticated"`
}

// IsEnabled will return if this single user is enabled.
func (u *User) IsEnabled() bool {
	return u.DisabledAt == nil
}

// IsModerator will return if the user has moderation privileges.
func (u *User) IsModerator() bool {
	_, hasModerationScope := utils.FindInSlice(u.Scopes, moderatorScopeKey)
	return hasModerationScope
}
