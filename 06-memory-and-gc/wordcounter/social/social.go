package social

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	gc "github.com/patrickmn/go-cache"
)

// use a cache for downloaded data
const (
	expiration = time.Minute * 10
	cleanup    = time.Hour
)

var cache = gc.New(expiration, cleanup)

// GetUserName returns the name of the user with id userID
func GetUserName(userID int) (string, error) {
	users, err := getUsers()
	if err != nil {
		return "", err
	}

	// search for id userID to get the name
	for _, user := range users {
		if user.ID == userID {
			return user.Name, nil
		}
	}

	return "", errors.New("user id not found")
}

// GetUserText returns all posts body and related comments
// of the user with id userID
func GetUserText(userID int) ([]string, error) {
	var text []string

	posts, err := getPosts()
	if err != nil {
		return nil, err
	}

	comments, err := getComments()
	if err != nil {
		return nil, err
	}

	// get all user's posts and save their ids
	var postIDs []int
	for _, post := range posts {
		if post.UserID == userID {
			postIDs = append(postIDs, post.ID)
			text = append(text, post.Body)
		}
	}

	// get all comments related to the user's posts
	for _, comment := range comments {
		for _, postID := range postIDs {
			if comment.PostID == postID {
				text = append(text, comment.Body)
			}
		}
	}

	return text, nil
}

func getUsers() ([]User, error) {
	// Look in the cache
	v, found := cache.Get(usersURI)
	if found {
		// cache hit
		return v.([]User), nil
	}

	// cache miss: download data
	resp, err := http.Get(usersURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode downloaded JSON data
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	// save data into the cache
	cache.Set(usersURI, users, expiration)

	return users, nil
}

func getPosts() ([]Post, error) {
	var posts []Post

	// for each posts uri
	for _, uri := range postsURIs {
		// Look in the cache
		v, found := cache.Get(uri)
		if found {
			// cache hit
			posts = append(posts, v.([]Post)...)
			continue
		}

		// cache miss
		err := func() error {
			// download data
			resp, err := http.Get(uri)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// decode downloaded JSON data
			var newPosts []Post
			if err := json.NewDecoder(resp.Body).Decode(&newPosts); err != nil {
				return err
			}

			// save data into the cache
			cache.Set(uri, newPosts, expiration)

			// add new posts to result
			posts = append(posts, newPosts...)

			return nil
		}()
		if err != nil {
			return nil, err
		}

	}

	return posts, nil
}

func getComments() ([]Comment, error) {
	var comments []Comment

	// for each posts uri
	for _, uri := range commentsURIs {
		// Look in the cache
		v, found := cache.Get(uri)
		if found {
			// cache hit
			comments = append(comments, v.([]Comment)...)
			continue
		}

		// cache miss
		err := func() error {
			// download data
			resp, err := http.Get(uri)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// decode downloaded JSON data
			var newComments []Comment
			if err := json.NewDecoder(resp.Body).Decode(&newComments); err != nil {
				return err
			}

			// save data into the cache
			cache.Set(uri, newComments, expiration)

			// add new posts to result
			comments = append(comments, newComments...)

			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	return comments, nil
}
