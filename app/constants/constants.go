package constants

import (
	"os"

	"github.com/FernandoVT10/go-auth/app/utils"
)

var SECRET_JWT_KEY string

func Initialize() {
    SECRET_JWT_KEY = os.Getenv("SECRET_JWT_KEY")
    if SECRET_JWT_KEY == "" {
        utils.LogFatal("SECRET_JWT_KEY env variable is required")
    }
}
