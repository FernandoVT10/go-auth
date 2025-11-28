package main

import (
	"net/http"

	"github.com/FernandoVT10/go-auth/app/controllers"
	httpUtils "github.com/FernandoVT10/go-auth/app/utils/http"
)

func GetRoutes() Router {
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

    return router
}
