package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tmdgo/dependencies"
	"github.com/tmdgo/reflection/functions"
)

func routerHandler(router *Router, route Route, handleFunc interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		manager := dependencies.Manager{}
		manager.InitWithOtherManager(router.manager)

		if route.UseVars {
			manager.Add(Vars{Value: mux.Vars(r)})
		}
		if route.UseOptionalVars {
			manager.Add(OptionalVars{Value: r.URL.Query()})
		}
		if route.UseRequestModel {
			requestModel := functions.CallFunc(route.RequestModel, nil)[0].Interface()
			_ = json.NewDecoder(r.Body).Decode(requestModel)
			manager.AddModel(requestModel)
		}

		handlerResult := manager.CallFunc(handleFunc)

		result := handlerResult[0].Interface().(Result)
		routerErr := handlerResult[1].Interface().(Error)

		if routerErr.Err != nil {
			jsonModel, err := json.Marshal(jsonError{
				Status:  "Error",
				Message: routerErr.Message,
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(routerErr.StatusCode)
			w.Write(jsonModel)

			return
		}

		jsonModel, err := json.Marshal(result.Model)

		if err != nil {
			return
		}

		w.WriteHeader(result.StatusCode)
		w.Write(jsonModel)
	}
}
