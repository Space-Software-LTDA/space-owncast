package admin

import (
	"net/http"

	"github.com/Space-Software-LTDA/owncast/persistence/configrepository"
	webutils "github.com/Space-Software-LTDA/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// ResetYPRegistration will clear the YP protocol registration key.
func ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	configRepository := configrepository.Get()
	if err := configRepository.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "reset")
}
