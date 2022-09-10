package main

import (
	"net/http"
	"webserver/controller"
	"webserver/storage"
)

func main() {
	memoryStorage := storage.NewStorage()
	ctrl := controller.NewController(memoryStorage)
	mux := http.NewServeMux()
	mux.HandleFunc("/create", ctrl.Create)
	mux.HandleFunc("/make_friends", ctrl.MakeFriends)
	mux.HandleFunc("/delete", ctrl.Delete)
	mux.HandleFunc("/friends/", ctrl.GetFriendsByUserId)
	mux.HandleFunc("/", ctrl.UpdateAgeByUserId)

	http.ListenAndServe("localhost:8081", mux)
}
