package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type ageUpdate struct {
	Age int `json:"age"`
}

type userDelete struct {
	TargetId int `json:"target_id"`
}

type FriendshipRequest struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type userData struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type userService struct {
	storage map[int]*userData
}

func main() {
	us := userService{make(map[int]*userData)}
	mux := http.NewServeMux()
	mux.HandleFunc("/create", us.Create)
	mux.HandleFunc("/make_friends", us.MakeFriends)
	mux.HandleFunc("/delete", us.Delete)
	mux.HandleFunc("/friends/", us.GetFriendsByUserId)
	mux.HandleFunc("/", us.UpdateAgeByUserId)

	http.ListenAndServe("localhost:8081", mux)
}

func (us *userService) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userData := &userData{}
	err = json.Unmarshal(jsonBody, &userData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	userID := rand.Int()
	us.storage[userID] = userData
	rw.Write([]byte(fmt.Sprintf(`{"userID":%d}`, userID)))
	rw.WriteHeader(http.StatusCreated)
}

func (us *userService) MakeFriends(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	friendshipRequest := &FriendshipRequest{}
	err = json.Unmarshal(jsonBody, &friendshipRequest)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	_, ok := us.storage[friendshipRequest.SourceId]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	us.storage[friendshipRequest.SourceId].Friends = append(us.storage[friendshipRequest.SourceId].Friends, friendshipRequest.TargetId)
	firstUser := us.storage[friendshipRequest.SourceId].Name
	secondUser := us.storage[friendshipRequest.TargetId].Name
	rw.Write([]byte(firstUser + " and " + secondUser + " are friends now"))
	rw.WriteHeader(http.StatusOK)
}

func (us *userService) Delete(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	userDelete := &userDelete{}
	err = json.Unmarshal(jsonBody, &userDelete)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	_, ok := us.storage[userDelete.TargetId]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	userName := us.storage[userDelete.TargetId].Name
	delete(us.storage, userDelete.TargetId)
	rw.Write([]byte(userName + ` was deleted`))
	rw.WriteHeader(http.StatusOK)
}

func (us *userService) GetFriendsByUserId(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/friends/"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error parsing user data`))
		return
	}

	user, ok := us.storage[userID]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	for _, val := range user.Friends {
		friend := us.storage[val]
		rw.Write([]byte(friend.Name))
	}

	rw.WriteHeader(http.StatusOK)
}

func (us *userService) UpdateAgeByUserId(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status Bad Request`))
		return
	}

	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error parsing user data`))
		return
	}

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`Status internal server error due to reading body`))
		return
	}

	ageUpdate := &ageUpdate{}
	err = json.Unmarshal(jsonBody, &ageUpdate)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user, ok := us.storage[userID]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	user.Age = ageUpdate.Age
	us.storage[userID] = user
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(fmt.Sprintf("возраст пользователя %s успешно обновлён. Теперь ему %d", user.Name, user.Age)))
}
