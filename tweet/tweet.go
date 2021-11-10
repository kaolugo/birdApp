package tweet

import (
	"time"

	"github.com/google/uuid"
)

// what should a tweet know ??
type Tweet struct {
	// which database the tweet is writing to (type interface whatever the fuck)
	//redisDB RedisDB

	// the userID of the user that this tweet is associated with
	UserID uuid.UUID `json:"userID"`
	// the unique tweetID of this tweet (for databse purposes I think lol)
	//tweetID string
	TweetID uuid.UUID `json:"tweetID"`

	// time that this tweet was created
	Time int64 `json:"time"`
	// actual content of this tweet
	Content string `json:"content"`

	// include additional data about images that may be associated with this tweet
	// but implement later bc idk wtf I'm doing
}

// constructor ....
// returns a brand spanking new empty tweet with a redis database associated with it
func New(user uuid.UUID, inputText string) *Tweet {
	// get current time using time go package
	now := time.Now()
	// convert to seconds so it is easier to reorder
	timeInSecs := now.Unix()

	return &Tweet{
		//redisDB: tweetData,
		UserID:  user,
		TweetID: uuid.New(),
		Time:    timeInSecs,
		Content: inputText,
	}
}

// functionality that a tweet should have !!!

// prep tweet to be uploaded to redis database
// string is the key for redis database, which is the string of tweetID
// the map is the tweet struct converted to a hash
func (t *Tweet) PrepTweet() (string, map[string]interface{}) {
	// convert tweet object to map
	// do this programatically
	tweetMap := map[string]interface{}{
		"UserID":  t.UserID.String(),
		"TweetID": t.TweetID.String(),
		"Time":    t.Time,
		"Content": t.Content,
	}

	// convert uuid type to string and attach a tag to it
	id := t.TweetID.String()
	tweetTag := "t:"
	id = tweetTag + id

	return id, tweetMap
}

// REMOVED: edit the content of a tweet

// REMOVED: delete a wholeass tweet

// return tweetID

// potential image things to consider
// add an image associated with a tweet
// replace an already existing image
// delete an image associated with a tweet
