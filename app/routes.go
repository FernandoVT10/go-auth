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
            httpUtils.HandleError(w, err)
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

    router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
        var data controllers.LoginUserData

        err := httpUtils.ParseJSON(r, &data)
        if err != nil {
            httpUtils.HandleError(w, err)
            return
        }

        err = data.Validate()
        if err != nil {
            httpUtils.SendErrorMsg(w, http.StatusBadRequest, err.Error())
            return
        }

        token, err := controllers.LoginUser(data)
        if err != nil {
            httpUtils.HandleError(w, err)
            return
        }

        res := map[string]string{"token": token}
        httpUtils.SendJSONResponse(w, http.StatusOK, res)
    })

    router.Post("/isAuthenticated", func(w http.ResponseWriter, r *http.Request) {
        token, err := httpUtils.GetAuthToken(r)
        if err != nil {
            httpUtils.HandleError(w, err)
            return
        }

        res := map[string]bool{
            "isAuthenticated": controllers.IsAuthenticated(token),
        }
        httpUtils.SendJSONResponse(w, http.StatusOK, res)
    })

    router.Get("/protectedRoute", func(w http.ResponseWriter, r *http.Request) {
        err := controllers.AuthenticateUser(r)
        if err != nil {
            httpUtils.HandleError(w, err)
            return
        }

        res := map[string]string{"message": "This is a protected route, if you're here, it means you're authenticated"}
        httpUtils.SendJSONResponse(w, http.StatusOK, res)
    })

    return router
}
