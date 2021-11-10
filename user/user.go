package user

import (
	"github.com/HENNGE/kaoru-BirdApp/tweet"
	"github.com/google/uuid"
)

// what should a user know ?
type User struct {
	// the name of this user
	Name string `json:"name"`
	// the userID of this user
	UserID uuid.UUID `json:"userID"`

	// convert any uuid to string before adding to any of these arrays
	// because apparently redis doesn't like uuid types

	// userIDs
	// array/slice of follower userIDs
	Followers []string `json:"followers"`
	// array/slice of following userIDs
	Following []string `json:"following"`

	// tweetIDs
	// slice of tweets that they have posted
	Tweets []string `json:"tweets"`
}

type Timeline struct {
	AllTweets []tweet.Tweet
}

type Roster struct {
	AllUsers []User
}

// constructor babyyyy
func New(name string) *User {
	// do I have to initialize all the slices in New() as well ?
	id := uuid.New()

	return &User{
		Name:   name,
		UserID: id,
		// a user is already following themselves by default
		Following: []string{id.String()},
	}
}

// prep user to be uploaded to redis database
// string is the key for redis database
// the map is the user struct converted to a hash
func (u *User) PrepUser() (string, map[string]interface{}) {
	// convert user object to a map
	userMap := map[string]interface{}{
		"Name":      u.Name,
		"UserID":    u.UserID.String(),
		"Followers": u.Followers,
		"Following": u.Following,
		"Tweets":    u.Tweets,
	}

	// convert uuid type to string and attach a tag
	id := "u:" + u.UserID.String()

	return id, userMap
}

// also identified in implementation of tweet things
// function to add a tweet to a user's tweet slice
// function to remove a tweet from a user's tweet slice

// what should a user be able to do ??

// display profile
// follow other users
// unfollow other users
// display personal timeline page
// display global timeline page

// TODO: Think about google OAuth integration before implementing this change
// sign in using google
// sign up using google
