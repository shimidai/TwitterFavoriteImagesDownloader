![GitHub](https://img.shields.io/github/license/shimidai/TwitterFavoriteImagesDownloader)

# TwitterFavoriteImagesDownloader
A program to save the images of your favorite tweets locally.

## Usage 
It is assumed that the twitter app has already been created.  
If you haven't created it yet, please do so from [here](https://developer.twitter.com/en/apps/).  

1. Copy `.env.example` and create `.env`.
2. Set value of each item in `.env` appropriately.
   - `SCREEN_NAME` is your twitter ID.
   - Fill the following items with the values obtained from "Keys and tokens" of the twitter app you created.
     - `CONSUMER_KEY`
     - `CONSUMER_SECRET`
     - `ACCESS_TOKEN`
     - `ACCESS_TOKEN_SECRET`
3. Execute `go run ./...` to start fetching and saving images. 

## About saved images
The file name will be created with format `{DateTime}_{TweetID}_{UserID}_{ImageID}.png`.  
Such as `20170823190152_900432923400310784_783214_DH7xKArXkAENtyr.png`.

If the image fails to be saved, a log will be emitted.  

## Open the target page for each ID
- By TweetID
  - https://twitter.com/i/web/status/900432923400310784
- By UserID
  - https://twitter.com/intent/user?user_id=783214
- By ImageID
  - https://pbs.twimg.com/media/DH7xKArXkAENtyr?format=png&name=medium

# License
The source code is licensed MIT.
