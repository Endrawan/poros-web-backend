package seeder

import (
	"github.com/divisi-developer-poros/poros-web-backend/config"
)

var (
	dbModel    config.DBModel
	connection = dbModel.PostgreConn()
)

// Execute menjalankan seeder yang telah dibuat
func Execute() {
	TagSeeder()
	PostTypeSeeder()
	UserTypeSeeder()
	UserSeeder()
}
