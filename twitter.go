package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	screenName        string
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
)

// newTwitterClient generates twitter-client.
func newTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// newFavoriteListParams generates new twitter.FavoriteListParams.
func newFavoriteListParams(maxID int64) *twitter.FavoriteListParams {
	params := twitter.FavoriteListParams{
		ScreenName:      screenName,
		Count:           count, // Maximum is 200
		IncludeEntities: twitter.Bool(true),
		TweetMode:       "extended",
	}

	if maxID > 0 {
		params.MaxID = maxID
	}

	return &params
}
