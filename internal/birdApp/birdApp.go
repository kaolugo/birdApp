package birdApp

// handle actual application logic here

import (
	"encoding/json"
	"fmt"
	"net/http"

	"io"
	"log"
	"sort"

	"github.com/HENNGE/kaoru-BirdApp/tweet"
	"github.com/HENNGE/kaoru-BirdApp/user"
	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

// add type interfaces to work with here (all of these will be separate packages)
// potential types: Amazon S3 SDK, tweets, users, images ?
// redis type interface
type RedisDB interface {
	// add functions that I want the database to be able to handle here

	// upload new content or any changes
	Upload(ID string, content []byte) error

	// get a tweet from database with tweetID
	Get(ID string) ([]byte, error)

	// deleting a tweet / user given tweetID + userID
	Delete(ID string) error

	// add new user to the list of all users in database
	UpdateUsers(ID string) error

	// get list of all users from the database
	GetAllUsers() ([]string, error)

	// add a tweet to the global timeline
	UpdateGlobal(ID string) error

	// delete a tweet from the global timeline
	RemoveGlobal(ID string) error

	// get the global timeline
	GetGlobal() ([]string, error)
}

// this struct is for receiving new tweets / edits from the front end
// and parsing the info with json decoder
type NewTweet struct {
	UserID  uuid.UUID `json:"userID"`
	TweetID uuid.UUID `json:"tweetID"`
	Content string    `json:"content"`
}

// this struct is for testing purposes
// for receiving new users from idk where and parsing the info with json decoder
type NewUser struct {
	Name     string    `json:"name"`
	UserID   uuid.UUID `json:"userID"`
	FollowID uuid.UUID `json:"followID"`
}

// for verifying follows
type verifyFollow struct {
	Followed string `json:"followed"`
}

// func New() for BirdApp
func New(tweetData RedisDB) *BirdApp {
	fmt.Println("Hi, I'm BirdApp !")
	return &BirdApp{
		// TwinA
		// TODO: add hatever assignment needed here
		// ex: attribute constructorParameters
		// TODO:
		redisDB: tweetData,
	}
}

// type definition of BirdApp
type BirdApp struct {
	// whatever the BirdApp needs to know / keep track of
	redisDB RedisDB
	// all the other users (array of users) -- for displaying timeline
}

// usrful function for removing an element from a slice
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// follows + user interactions
// I want to sign in to my account

// I want to sign out of my account

// I want to sign up for an account ? (birth a new user)

// for testing purposes to see if the app is functioning
func (birdApp *BirdApp) Test(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Pong!\n")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create a new User for testing purposes
// returns a json of the user that was just created
func (birdApp *BirdApp) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUserInfo NewUser

	// EXPECT: name from the http request json
	// get necessary info from the http request
	json.NewDecoder(r.Body).Decode(&newUserInfo)

	// create new user using the information gathered
	newUser := user.New(newUserInfo.Name)

	//userID, userInfo := newUser.PrepUser()
	userID := newUser.UserID.String()

	// turn user struct into json
	userJson, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println(err)
	}

	// upload new user to the redis database
	err2 := birdApp.redisDB.Upload(userID, userJson)

	if err2 != nil {
		fmt.Println("redisDB.Upload for users didn't work")
	}

	// add newUser to the list of all users in the database
	birdApp.redisDB.UpdateUsers(userID)

	// TODO: take out after testing
	_, err3 := io.WriteString(w, "something happened in createUser\n")
	if err3 != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// for testing purpose, send back user info that was just created
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)

}

// ability to fetch a user + userData from the database given a userID
// ability to view a user profile --> display my / a user's profile
func (birdApp *BirdApp) GetUser(w http.ResponseWriter, r *http.Request) {
	// json implementation
	// var wantedUser NewUser
	// EXPECT: userID (UUID) from the http request json
	// get necessary info from the http request
	//json.NewDecoder(r.Body).Decode(&wantedUser)
	//id := wantedUser.UserID.String()

	vars := mux.Vars(r)
	id := vars["id"]

	// call get user function from redis
	userJson, err := birdApp.redisDB.Get(id)

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}

// ability to get the list of all users existing in the database
func (birdApp *BirdApp) AllUsers(w http.ResponseWriter, r *http.Request) {
	// get list of users from the database
	userIDs, err := birdApp.redisDB.GetAllUsers()

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allUsers []user.User

	// for each userID string, make a userID struct
	for id := range userIDs {
		// fetch userJson from database
		userJson, _ := birdApp.redisDB.Get(userIDs[id])

		// turn said json into user struct
		var aUser user.User
		json.Unmarshal(userJson, &aUser)

		// put said struct into list
		allUsers = append(allUsers, aUser)
	}

	var response user.Roster
	response.AllUsers = allUsers

	// turn user roster struct into json
	responseJson, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

// check if the users in question are already following each other
func (birdApp *BirdApp) VerifyFollow(w http.ResponseWriter, r *http.Request) {
	//var theUser NewUser

	//json.NewDecoder(r.Body).Decode(&theUser)

	//friendID := theUser.FollowID.String()

	vars := mux.Vars(r)
	id := vars["id"]
	friendID := vars["friendID"]

	// call get user function from redis
	userJson, err := birdApp.redisDB.Get(id)
	friendJson, err2 := birdApp.redisDB.Get(friendID)

	// error handling
	if err != nil || err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// turn json data into structs
	var you user.User
	json.Unmarshal(userJson, &you)

	var friend user.User
	json.Unmarshal(friendJson, &friend)

	var followed = false

	// check if the user is already following their friend
	for i := range you.Following {
		if you.Following[i] == friendID {
			followed = true
		}
	}

	var result verifyFollow

	if followed == true {
		result = verifyFollow{"true"}
	} else {
		result = verifyFollow{"false"}
	}

	//var result = verifyFollow{followed}

	responseJson, _ := json.Marshal(result)

	// return to front end
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)

}

// ability to follow a user
func (birdApp *BirdApp) FollowUser(w http.ResponseWriter, r *http.Request) {

	var followUser NewUser
	json.NewDecoder(r.Body).Decode(&followUser)

	//id := followUser.UserID.String()
	friendID := followUser.FollowID.String()

	vars := mux.Vars(r)
	id := vars["id"]

	// call get user function from redis
	userJson, err := birdApp.redisDB.Get(id)
	friendJson, err2 := birdApp.redisDB.Get(friendID)

	fmt.Println(friendID)

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// turn json data into structs
	var you user.User
	json.Unmarshal(userJson, &you)

	var friend user.User
	json.Unmarshal(friendJson, &friend)

	// add friendID to the slice of following in the user
	you.Following = append(you.Following, friendID)
	// add your ID to the slice of followers in friend
	friend.Followers = append(you.Followers, id)

	// package these changes into json again
	userJson, err = json.Marshal(you)
	friendJson, err2 = json.Marshal(friend)

	// error checking
	if err != nil || err2 != nil {
		fmt.Println(err)
	}

	// upload these edits to the database for both users
	birdApp.redisDB.Upload(id, userJson)
	birdApp.redisDB.Upload(friendID, friendJson)

	// just to confirm something happened
	_, err = io.WriteString(w, "This user followed a friend\n")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ability to unfollow a user
func (birdApp *BirdApp) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	// EXPECT: main userID, userID of the person they want to unfollow
	var unfollowUser NewUser
	json.NewDecoder(r.Body).Decode(&unfollowUser)

	//id := unfollowUser.UserID.String()
	friendID := unfollowUser.FollowID.String()

	vars := mux.Vars(r)
	id := vars["id"]

	// call get user function from redis
	userJson, err := birdApp.redisDB.Get(id)
	friendJson, err2 := birdApp.redisDB.Get(friendID)

	// error handling
	if err != nil || err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// turn json data into structs
	var you user.User
	json.Unmarshal(userJson, &you)

	var friend user.User
	json.Unmarshal(friendJson, &friend)

	// delete friendID from the following slice in user
	//index := sort.SearchStrings(you.Following, friendID)

	var index int
	for i := range you.Following {
		if friendID == you.Following[i] {
			index = i
		}
	}

	you.Following = RemoveIndex(you.Following, index)

	// delete your id from your friend's follower slice
	for i := range friend.Followers {
		if id == friend.Followers[i] {
			index = i
		}
	}
	friend.Followers = RemoveIndex(friend.Followers, index)

	// package these changes into json again
	userJson, err = json.Marshal(you)
	friendJson, err2 = json.Marshal(friend)

	// error checking
	if err != nil || err2 != nil {
		fmt.Println(err)
	}

	// upload these edits to the database for both users
	birdApp.redisDB.Upload(id, userJson)
	birdApp.redisDB.Upload(friendID, friendJson)

	// just to confirm something happened
	_, err = io.WriteString(w, "This user unfollowed a friend\n")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// timeline related things
// list all tweets by the users I am following + myself (personal timeline)
// return a list of json
func (birdApp *BirdApp) ShowPersonal(w http.ResponseWriter, r *http.Request) {
	// json implementation
	// EXPECT: userID
	//var personal NewUser
	//json.NewDecoder(r.Body).Decode(&personal)
	// get the userID whose personal timeline you want to show from frontend
	//id := personal.UserID.String()

	vars := mux.Vars(r)
	id := vars["id"]

	userJson, err := birdApp.redisDB.Get(id)

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// turn userJson into a usable struct
	var you user.User
	json.Unmarshal(userJson, &you)

	// get the "following" list from their user data
	following := you.Following

	var timeline []tweet.Tweet

	// iterate over the "following list"
	for i := range following {
		// for each friend id string
		friendJson, err2 := birdApp.redisDB.Get(following[i])

		// error handling
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// turn Json into user struct
		var friend user.User
		json.Unmarshal(friendJson, &friend)

		// iterate over the list of tweets they have
		// get the tweet from database for each
		// list into a list of tweet structs
		for j := range friend.Tweets {
			tweetJson, err3 := birdApp.redisDB.Get(friend.Tweets[j])

			// error handling
			if err3 != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var aTweet tweet.Tweet
			json.Unmarshal(tweetJson, &aTweet)

			timeline = append(timeline, aTweet)
		}
	}

	// sort in chronological order
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].Time > timeline[j].Time })

	// turn this list into json
	var response user.Timeline
	response.AllTweets = timeline

	// turn into json
	responseJson, _ := json.Marshal(response)

	// return to frontend
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

// list all tweets by all users + myself (global timeline)
// return a list of json
// {"global": [{}{}{}]}
func (birdApp *BirdApp) ShowGlobal(w http.ResponseWriter, r *http.Request) {

	// get global timeline from the database
	tweetIDs, err := birdApp.redisDB.GetGlobal()

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var timeline []tweet.Tweet

	// for each tweetID string, make a tweetID struct
	for id := range tweetIDs {
		// fetch tweetJson from database
		tweetJson, _ := birdApp.redisDB.Get(tweetIDs[id])

		// turn said json into tweet struct
		var aTweet tweet.Tweet
		json.Unmarshal(tweetJson, &aTweet)

		// put said struct into list
		timeline = append(timeline, aTweet)
	}

	// sort the timeline according to the time
	// greatest time -> least time
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].Time > timeline[j].Time })

	// declare a new global timeline struct from user package
	var response user.Timeline
	response.AllTweets = timeline

	// turn global timeline struct into Json
	responseJson, _ := json.Marshal(response)

	// return to frontend
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

// tweet + content related things
// I want to create a new tweet
func (birdApp *BirdApp) NewTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	fmt.Println(userID)

	userUUID, _ := uuid.Parse(userID)

	var newTweetInfo NewTweet

	// EXPECT: tweetContent from the http request json
	// get necessary info from the http request
	json.NewDecoder(r.Body).Decode(&newTweetInfo)
	newTweet := tweet.New(userUUID, newTweetInfo.Content)

	// convert tweet object to redis-compatible format abd store it
	// key string, and value hash
	//twitID, twitContent := newTweet.PrepTweet()
	twitID := newTweet.TweetID.String()

	// turn tweet struct into json
	tweetJson, err := json.Marshal(newTweet)
	if err != nil {
		fmt.Println(err)
	}

	// upload said tweet to redis by calling uploadNewTweet
	// call function from redisDB type interface
	err2 := birdApp.redisDB.Upload(twitID, tweetJson)

	if err2 != nil {
		fmt.Println("redisDB.UploadNewTweet didn't work lol")
	}

	// add this tweet to the global timeline
	err2 = birdApp.redisDB.UpdateGlobal(twitID)
	if err2 != nil {
		fmt.Println("updating the global timeline didn't work...")
	}

	// add this tweetID to user's tweet slice
	youJson, err3 := birdApp.redisDB.Get(userID)
	if err3 != nil {
		fmt.Println("could not add this tweet to the user's tweet list")
	}

	// turn json into user struct
	var you user.User
	json.Unmarshal(youJson, &you)

	you.Tweets = append(you.Tweets, twitID)

	youJson, err3 = json.Marshal(you)
	if err3 != nil {
		fmt.Println(err)
	}

	// upload to Redis
	birdApp.redisDB.Upload(userID, youJson)

	// return tweet data for testing purposes
	w.Header().Set("Content-Type", "application/json")
	w.Write(tweetJson)
}

// I want to edit a tweet
func (birdApp *BirdApp) EditTweet(w http.ResponseWriter, r *http.Request) {
	// receive edits necessary
	var tweetEdits NewTweet

	// EXPECT: tweetContent from the http request json
	// get necessary info from the http request
	json.NewDecoder(r.Body).Decode(&tweetEdits)

	// twitID is a string with the tweetID
	//twitID := tweetEdits.TweetID.String()

	vars := mux.Vars(r)
	twitID := vars["id"]

	originalTweet, err := birdApp.redisDB.Get(twitID)

	if err != nil {
		fmt.Println("could not retrieve the tweet for editTweet")
	}

	// turn json into tweet struct
	var tweet tweet.Tweet
	json.Unmarshal(originalTweet, &tweet)

	// replace the content with the new edits
	tweet.Content = tweetEdits.Content

	// repackage into json
	tweetJson, err2 := json.Marshal(tweet)
	if err2 != nil {
		fmt.Println(err)
	}

	// upload changes
	birdApp.redisDB.Upload(twitID, tweetJson)

}

// I want to delete a tweet
func (birdApp *BirdApp) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	// json implementation
	// EXPECT: tweetID + userID from http request json
	// get necessary info from the http request
	//var deletedTweet NewTweet
	//json.NewDecoder(r.Body).Decode(&deletedTweet)
	// convert id to string
	//id := deletedTweet.TweetID.String()

	fmt.Println("currently deleting tweet")

	vars := mux.Vars(r)
	id := vars["id"]

	// get userID of tweet
	tweetJson, err := birdApp.redisDB.Get(id)
	if err != nil {
		fmt.Println("failed to find the tweet")
	}
	// turn json into tweet struct
	var aTweet tweet.Tweet
	json.Unmarshal(tweetJson, &aTweet)

	userID := aTweet.UserID.String()

	// delete the tweet using tweetID from the redis database
	birdApp.redisDB.Delete(id)

	// delete this tweet from the global timeline
	birdApp.redisDB.RemoveGlobal(id)

	// TODO: take out the tweet from the user's list of tweets
	youJson, err3 := birdApp.redisDB.Get(userID)
	if err3 != nil {
		fmt.Println("could not delete this tweet from the user's tweet list")
	}

	// turn json into user struct
	var you user.User
	json.Unmarshal(youJson, &you)

	// find the tweet in the user's list of tweets and remove
	//index := sort.SearchStrings(you.Tweets, id)
	//RemoveIndex(you.Tweets, index)
	var index int
	for i := range you.Tweets {
		if id == you.Tweets[i] {
			index = i
		}
	}

	if index != -1 {
		you.Tweets = RemoveIndex(you.Tweets, index)
	}

	youJson, err3 = json.Marshal(you)
	if err3 != nil {
		fmt.Println(err)
	}

	fmt.Println(you.Tweets)

	// upload to Redis
	birdApp.redisDB.Upload(userID, youJson)

}

// tweet testing purposes
func (birdApp *BirdApp) GetTweet(w http.ResponseWriter, r *http.Request) {
	// json implementation
	// EXPECT: tweetID wanted
	// get necessary info from the http request
	//var wantedTweet NewTweet
	//json.NewDecoder(r.Body).Decode(&wantedTweet)
	// convert id to string
	//id := wantedTweet.TweetID.String()

	vars := mux.Vars(r)
	id := vars["id"]

	// get the wanted tweet
	tweetJson, err := birdApp.redisDB.Get(id)

	// error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tweetJson)
}
