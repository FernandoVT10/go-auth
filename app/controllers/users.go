package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/FernandoVT10/go-auth/app/constants"
	"github.com/FernandoVT10/go-auth/app/db"
	"github.com/FernandoVT10/go-auth/app/utils"
	"github.com/FernandoVT10/go-auth/app/utils/http"
	"github.com/FernandoVT10/go-auth/app/validator"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const USERNAME_MIN_LEN = 6
const USERNAME_MAX_LEN = 20

const PASSWORD_MIN_LEN = 8
const PASSWORD_MAX_LEN = 1000

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

// returns false and "" when the user wasn't found
func getUserPasswordByUsername(username string) (bool, string) {
    row := db.QueryRow("SELECT Password FROM Users WHERE Username = ?", username)
    var password string
    err := row.Scan(&password)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // No user was found
            return false, ""
        }

        // Unexpected error was found
        utils.LogError(err.Error())
        return false, ""
    }

    return true, password
}

type RegisterUserData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (data RegisterUserData) Validate() error {
    return validator.ValidateStruct(
        validator.Field(
            "username",
            data.Username,
            validator.Required,
            validator.Length(USERNAME_MIN_LEN, USERNAME_MAX_LEN),
        ),
        validator.Field(
            "password",
            data.Password,
            validator.Required,
            validator.Length(PASSWORD_MIN_LEN, PASSWORD_MAX_LEN),
        ),
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

type LoginUserData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (data LoginUserData) Validate() error {
    return validator.ValidateStruct(
        validator.Field(
            "username",
            data.Username,
            validator.Required,
            validator.Length(USERNAME_MIN_LEN, USERNAME_MAX_LEN),
        ),
        validator.Field("password", data.Password, validator.Required),
    )
}

func LoginUser(data LoginUserData) (string, error) {
    wasFound, hashedPassword := getUserPasswordByUsername(data.Username)
    if !wasFound {
        return "", &httpUtils.HttpError{
            Code: http.StatusBadRequest,
            Message: "Username doesn't exist",
        }
    }

    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(data.Password))
    if err != nil {
        return "", &httpUtils.HttpError{
            Code: http.StatusBadRequest,
            Message: "Incorrect password",
        }
    }

    t := jwt.New(jwt.SigningMethodHS256)
    token, err := t.SignedString([]byte(constants.SECRET_JWT_KEY))
    if err != nil {
        utils.LogError(err.Error())
        return "", httpUtils.InternalServerError()
    }

    return token, nil
}

func IsAuthenticated(token string) bool {
    tkn, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
        return []byte(constants.SECRET_JWT_KEY), nil
    })

    return err == nil && tkn.Valid
}

func AuthenticateUser(r *http.Request) error {
    tkn, err := httpUtils.GetAuthToken(r)
    if err != nil {
        return err
    }

    if !IsAuthenticated(tkn) {
        return &httpUtils.HttpError{
            Code: http.StatusUnauthorized,
            Message: "You need to be logged in",
        }
    }

    return nil
}
