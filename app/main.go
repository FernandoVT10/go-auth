package main

import (
	"fmt"
	"net/http"

	"github.com/FernandoVT10/go-auth/app/controllers"
	"github.com/FernandoVT10/go-auth/app/db"
	"github.com/FernandoVT10/go-auth/app/utils"
	"github.com/FernandoVT10/go-auth/app/utils/http"
)

func main() {
    utils.LogInfo("Connecting to db...")
    err := db.Connect()
    if err != nil {
        utils.LogFatal(err.Error())
    }
    defer db.Close()

    utils.LogInfo("Initializing db...")
    err = db.Initialize()
    if err != nil {
        utils.LogFatal(err.Error())
    }

    router := NewRouter()

    router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
        var data controllers.RegisterUserData

        err := httpUtils.ParseJSON(r, &data)
        if err != nil {
            httpUtils.SendErrorMsg(w, http.StatusBadRequest, err.Error())
            return
        }

        err = data.Validate()
        if err != nil {
            httpUtils.SendErrorMsg(w, http.StatusBadRequest, err.Error())
            return
        }

        err = controllers.RegisterUser(data)
        if err != nil {
            httpUtils.HandleError(w, err)
            return
        }

        w.WriteHeader(http.StatusOK)
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if router.Serve(w, r) {
            return
        }

        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "404")
    })

    utils.LogInfo("Server listening on port 3000")
    http.ListenAndServe(":3000", nil);
}
