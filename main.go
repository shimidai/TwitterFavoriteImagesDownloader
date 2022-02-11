package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/joho/godotenv"
)

const (
	count = 200 // How many tweets to fetch per request. Maximum is 200.
	page  = 20  // How many requests to execute. count * page is the total number of tweets to fetch.
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading .env file")
	}

	var ok bool

	screenName, ok = os.LookupEnv("SCREEN_NAME")
	if !ok {
		log.Fatal("not set SCREEN_NAME")
	}
	consumerKey, ok = os.LookupEnv("CONSUMER_KEY")
	if !ok {
		log.Fatal("not set CONSUMER_KEY")
	}
	consumerSecret, ok = os.LookupEnv("CONSUMER_SECRET")
	if !ok {
		log.Fatal("not set CONSUMER_SECRET")
	}
	accessToken, ok = os.LookupEnv("ACCESS_TOKEN")
	if !ok {
		log.Fatal("not set ACCESS_TOKEN")
	}
	accessTokenSecret, ok = os.LookupEnv("ACCESS_TOKEN_SECRET")
	if !ok {
		log.Fatal("not set ACCESS_TOKEN_SECRET")
	}
}

func main() {
	client := newTwitterClient()

	var cursorID int64 // Fetch older than this ID.

	for i := 0; i < page; i++ {
		fmt.Printf("=== request %d ===\n", i+1)

		lastCursorID, isFinish, err := run(client, cursorID)
		if err != nil {
			panic(err)
		}
		if isFinish {
			break
		}
		cursorID = lastCursorID
	}

	fmt.Println("=== Finish ===")
}

func run(client *twitter.Client, cursorID int64) (lastCursorID int64, fin bool, err error) {
	favoListParams := newFavoriteListParams(cursorID)

	tweets, resp, err := client.Favorites.List(favoListParams)
	if err != nil {
		return 0, false, fmt.Errorf("fetch tweets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, false, fmt.Errorf("fetch tweets: status %v", resp.Status)
	}

	if len(tweets) == 0 {
		return 0, true, nil
	}

	for _, tweet := range tweets {
		if err := saveTweet(&tweet); err != nil {
			fmt.Printf("[ERROR] %v\n", err)
		}
	}

	lastCursorID = int64(tweets[len(tweets)-1].ID)

	// NOTE
	//   When cursor ID is specified, tweets before the specified ID will be fetched.
	//   This will cause the last tweet to be duplicated, so the value of ID is set to -1.
	lastCursorID -= 1

	return lastCursorID, false, nil
}
