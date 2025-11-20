package authrepository

import (
	"database/sql"

	"github.com/Space-Software-LTDA/owncast/webserver/handlers/generated"
)

type AuthRepository interface {
	CreateBanIPTable(db *sql.DB)
	BanIPAddress(address, note string) error
	IsIPAddressBanned(address string) (bool, error)
	GetIPAddressBans() ([]generated.IPAddress, error)
	RemoveIPAddressBan(address string) error
}
