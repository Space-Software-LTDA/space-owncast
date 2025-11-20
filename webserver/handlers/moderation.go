package handlers

import (
	"net/http"

	"github.com/Space-Software-LTDA/owncast/webserver/handlers/generated"
	"github.com/Space-Software-LTDA/owncast/webserver/handlers/moderation"
	"github.com/Space-Software-LTDA/owncast/webserver/router/middleware"
)

func (*ServerInterfaceImpl) GetUserDetails(w http.ResponseWriter, r *http.Request, userId string, params generated.GetUserDetailsParams) {
	middleware.RequireUserModerationScopeAccesstoken(moderation.GetUserDetails)(w, r)
}
