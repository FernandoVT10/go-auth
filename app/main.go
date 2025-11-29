package main

import (
	"fmt"
	"net/http"

	"github.com/FernandoVT10/go-auth/app/constants"
	"github.com/FernandoVT10/go-auth/app/db"
	"github.com/FernandoVT10/go-auth/app/utils"
)

func main() {
    constants.Initialize()

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

    router := GetRoutes()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if router.Serve(w, r) {
            return
        }

        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Not Found")
    })

    utils.LogInfo("Server listening on port 3000")
    http.ListenAndServe(":3000", nil);
}
