package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// saveTweet saves the image attached to a tweet.
func saveTweet(tweet *twitter.Tweet) error {
	createdAt, err := time.Parse(time.RubyDate, tweet.CreatedAt)
	if err != nil {
		return fmt.Errorf("parse created_at: %w", err)
	}

	saveDirPath, err := newSaveDirectory(createdAt)
	if err != nil {
		return err
	}

	if tweet.ExtendedEntities == nil {
		fmt.Println("[INFO] no ExtendedEntities tweet")
		fmt.Printf("\thttps://twitter.com/%s/status/%d\n", tweet.User.ScreenName, tweet.ID)
		return nil
	}

	for _, medium := range tweet.ExtendedEntities.Media {
		path := medium.MediaURLHttps

		if !isWantExtension(filepath.Ext(path)) {
			fmt.Printf("[INFO] following media extension is not supported: %s", path)
			break
		}

		mediumID := filepath.Base(path[:len(path)-len(filepath.Ext(path))])

		saveName := createdAt.Format("20060102150405") +
			"_" + tweet.IDStr +
			"_" + tweet.User.IDStr +
			"_" + mediumID +
			".png"

		savePath := filepath.Join(saveDirPath, saveName)
		_, err := os.Stat(savePath)
		if err == nil {
			break // Don't execute because the file already exists.
		}
		if !os.IsNotExist(err) {
			return fmt.Errorf("exists file: %w", err)
		}

		imageData, err := fetchImage(path)
		if err != nil {
			return fmt.Errorf("fetch image: \"%s\": %w", path, err)
		}

		if err := saveImage(savePath, imageData); err != nil {
			return fmt.Errorf("save image: \"%s\": %w", savePath, err)
		}

		fmt.Println("done: " + saveName)

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// newSaveDirectory returns string of the save destination directory.
// If the directory does not exist, it will be created.
func newSaveDirectory(createdAt time.Time) (string, error) {
	dir := path.Join("twitterFavoritePictures", createdAt.Format("200601"))

	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("exists save dir: %w", err)
		}

		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return "", fmt.Errorf("mkdir save dir: %w", err)
		}
	}

	return dir, nil
}

// isWantExtension determines if extension is expected one.
func isWantExtension(ext string) bool {
	wantExts := []string{".png", ".jpg", ".jpeg"}

	in := strings.ToLower(ext)

	for _, v := range wantExts {
		if in == v {
			return true
		}
	}
	return false
}

// fetchImage fetches image data.
func fetchImage(imageURL string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, imageURL+"?format=png&name=large", nil)
	if err != nil {
		return nil, fmt.Errorf("new get request: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get request status: %v", response.Status)
	}

	return io.ReadAll(response.Body)
}

// saveImage saves image data.
func saveImage(savePath string, data []byte) error {
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		if err2 := os.Remove(savePath); err2 != nil {
			err = fmt.Errorf("%w: (fail remove file)", err)
		}
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
