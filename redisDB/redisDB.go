package redisDB

import (
	"fmt"

	//"encoding/json"

	"errors"

	"context"

	"github.com/go-redis/redis"
)

// New returns a new redis service .....
/*
func New() *Redis {
	return &Redis{}
}
*/

// make a redis struct.....probably
type RedisDB struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("no matching records found in redis database :(")
	Ctx    = context.TODO()
)

// functions to write to and get from database here

// create a new Redis database !!
func New(address string) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	// check for a pong to test if the connection is working
	pong, err := client.Ping().Result()

	// check for a pong by printing the result
	fmt.Println(pong, err)

	if err != nil {
		return nil, err
	}

	return &RedisDB{
		Client: client,
	}, nil
}

// just for compilation purposes

// function definitions

// upload function - the verion where it takes json as input
func (db *RedisDB) Upload(ID string, content []byte) error {
	err := db.Client.Set(ID, content, 0).Err()

	if err != nil {
		return err
	}
	return nil
}

// IMPLEMENT NOW !!!!!
func (db *RedisDB) Get(ID string) ([]byte, error) {
	info, err := db.Client.Get(ID).Result()

	if err != nil {
		return nil, err
	}
	return []byte(info), err
}

// function to delete tweet or user from the database given ID
// refactor for users as needed (although CRUD operations are unnecessary for users)
func (db *RedisDB) Delete(ID string) error {
	err := db.Client.Del(ID).Err()

	if err != nil {
		return err
	}
	return nil
}

// update list of users in the database
func (db *RedisDB) UpdateUsers(ID string) error {
	key := "allUsers"
	newUser := []string{ID}

	err := db.Client.RPush(key, newUser).Err()

	if err != nil {
		return err
	}
	return nil
}

func (db *RedisDB) GetAllUsers() ([]string, error) {
	key := "allUsers"

	users, err := db.Client.LRange(key, 0, -1).Result()

	if err != nil {
		return nil, err
	}
	return users, err
}

// function to update the global timeline list
// which is just a list of the ids of all tweets ever made
func (db *RedisDB) UpdateGlobal(ID string) error {
	key := "globalTimeline"
	newTweet := []string{ID}

	err := db.Client.RPush(key, newTweet).Err()

	if err != nil {
		return err
	}
	return nil
}

func (db *RedisDB) RemoveGlobal(ID string) error {
	key := "globalTimeline"

	// remove the tweet id
	err := db.Client.LRem(key, 1, ID).Err()

	if err != nil {
		return err
	}
	return nil
}

func (db *RedisDB) GetGlobal() ([]string, error) {
	key := "globalTimeline"

	tweets, err := db.Client.LRange(key, 0, -1).Result()

	if err != nil {
		return nil, err
	}
	return tweets, err
}
