package main

import "github.com/CollActionteam/collaction_backend/pkg/handler/http"

func main() {
	router := http.NewRouter()
	http.NewContactHandler().Register(router)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
