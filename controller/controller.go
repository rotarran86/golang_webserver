package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"webserver/model"
	"webserver/request"
	"webserver/storage"
)

type Controller struct {
	storage storage.Storage
}

func NewController(storage storage.Storage) *Controller {
	return &Controller{storage: storage}
}

func (c *Controller) Create(rw http.ResponseWriter, r *http.Request) {
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

	userData := &request.UserData{}
	err = json.Unmarshal(jsonBody, &userData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("Status internal server error due to unmarshal body:%s", err)))
		return
	}

	user := &model.User{Name: userData.Name, Age: userData.Age, Friends: userData.Friends}
	userID := c.storage.Add(user)

	rw.Write([]byte(fmt.Sprintf(`{"userID":%d}`, userID)))
	rw.WriteHeader(http.StatusCreated)
}

func (c *Controller) MakeFriends(rw http.ResponseWriter, r *http.Request) {
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

	friendshipRequest := &request.FriendshipRequest{}
	err = json.Unmarshal(jsonBody, &friendshipRequest)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user := c.storage.FindByUserId(friendshipRequest.SourceId)
	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	user.Friends = append(user.Friends, friendshipRequest.TargetId)
	secondUser := c.storage.FindByUserId(friendshipRequest.TargetId)
	rw.Write([]byte(user.Name + " and " + secondUser.Name + " are friends now"))
	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) Delete(rw http.ResponseWriter, r *http.Request) {
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

	userDelete := &request.UserDelete{}
	err = json.Unmarshal(jsonBody, &userDelete)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user := c.storage.FindByUserId(userDelete.TargetId)
	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`user with requested ID not found`))
		return
	}

	rw.Write([]byte(user.Name + ` was deleted`))
	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) GetFriendsByUserId(rw http.ResponseWriter, r *http.Request) {
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

	user := c.storage.FindByUserId(userID)
	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	for _, val := range user.Friends {
		friend := c.storage.FindByUserId(val)
		rw.Write([]byte(friend.Name))
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateAgeByUserId(rw http.ResponseWriter, r *http.Request) {
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

	ageUpdate := &request.AgeUpdate{}
	err = json.Unmarshal(jsonBody, &ageUpdate)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Status internal server error due to unmarshal body`))
		return
	}

	user := c.storage.FindByUserId(userID)
	if user == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("user with requested id %d not found", userID)))
		return
	}

	user.Age = ageUpdate.Age
	c.storage.Add(user)
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(fmt.Sprintf("возраст пользователя %s успешно обновлён. Теперь ему %d", user.Name, user.Age)))
}
