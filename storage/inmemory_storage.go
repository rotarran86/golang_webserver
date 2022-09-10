package storage

import (
	"math/rand"
	"webserver/model"
)

type Storage interface {
	Add(user *model.User) int
	Delete(userID int)
	FindByUserId(userID int) *model.User
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{make(map[int]*model.User)}
}

type InMemoryStorage struct {
	storage map[int]*model.User
}

func (ims *InMemoryStorage) Add(user *model.User) int {
	userID := rand.Int()
	ims.storage[userID] = user

	return userID
}

func (ims *InMemoryStorage) Delete(userID int) {
	_, ok := ims.storage[userID]
	if !ok {
		return
	}

	delete(ims.storage, userID)
}

func (ims *InMemoryStorage) FindByUserId(userID int) *model.User {
	userData, ok := ims.storage[userID]
	if !ok {
		return nil
	}

	return userData
}
