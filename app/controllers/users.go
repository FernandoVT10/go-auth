package controllers

import (
	"fmt"
	"net/http"

	"github.com/FernandoVT10/go-auth/app/db"
	"github.com/FernandoVT10/go-auth/app/utils"
	"github.com/FernandoVT10/go-auth/app/utils/http"
	"github.com/FernandoVT10/go-auth/app/validator"

	"golang.org/x/crypto/bcrypt"
)

const SALT_ROUNDS = 10

func usernameExists(username string) bool {
    row := db.QueryRow("SELECT COUNT(*) FROM Users WHERE Username = ?", username)

    var count int
    err := row.Scan(&count)
    if err != nil {
        utils.LogError(err.Error())
        return false
    }

    return count > 0
}

type RegisterUserData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (data RegisterUserData) Validate() error {
    return validator.ValidateStruct(
        validator.Field("username", data.Username, validator.Required, validator.Length(6, 20)),
        validator.Field("password", data.Password, validator.Required, validator.Length(8, 1000)),
    )
}

func RegisterUser(data RegisterUserData) error {
    if usernameExists(data.Username) {
        msg := fmt.Sprintf("Username \"%s\" is already taken", data.Username)
        return &httpUtils.HttpError{
            Code: http.StatusBadRequest,
            Message: msg,
        }
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), SALT_ROUNDS)
    if err != nil {
        utils.LogError(err.Error())
        return httpUtils.InternalServerError()
    }

    _, err = db.Exec("INSERT INTO Users (Username, Password) VALUES (?, ?)", data.Username, hashedPassword)
    if err != nil {
        utils.LogError(err.Error())
        return httpUtils.InternalServerError()
    }

    return nil
}

