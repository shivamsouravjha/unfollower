package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseFollowingURL = "https://api.github.com/user/following?per_page=100"
	baseFollowersURL = "https://api.github.com/user/followers?per_page=100"
	token            = ""
)

type User struct {
	Login string `json:"login"`
}

func fetchPaginatedData(baseURL, token string) ([]User, error) {
	var allUsers []User
	page := 1

	for {
		url := fmt.Sprintf("%s&page=%d", baseURL, page)
		users, err := fetchGitHubData(url, token)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, users...)

		// If less than 100 users were returned, we are on the last page
		if len(users) < 100 {
			break
		}

		page++
	}

	return allUsers, nil
}

func fetchGitHubData(url, token string) ([]User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add authorization header with the token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return users, nil
}

func main() {
	err := godotenv.Load()
	token := os.Getenv("TOKEN")

	// Fetch all following users with pagination
	following, err := fetchPaginatedData(baseFollowingURL, token)
	if err != nil {
		fmt.Printf("Error fetching following data: %v\n", err)
		return
	}

	// Store following users in a map and initialize them to false
	followingMap := make(map[string]bool)
	for _, user := range following {
		followingMap[user.Login] = false
	}

	// Fetch all followers with pagination
	followers, err := fetchPaginatedData(baseFollowersURL, token)
	if err != nil {
		fmt.Printf("Error fetching followers data: %v\n", err)
		return
	}

	// Mark followers in the map as true
	for _, user := range followers {
		if _, found := followingMap[user.Login]; found {
			followingMap[user.Login] = true
		}
	}

	notFollowingBack := []string{}
	for user, followsBack := range followingMap {
		if !followsBack {
			notFollowingBack = append(notFollowingBack, user)
		}
	}

	// Print the users who don't follow you back
	if len(notFollowingBack) > 0 {
		fmt.Println("Users you follow who don't follow you back:")
		for _, user := range notFollowingBack {
			fmt.Println(user)
		}
	} else {
		fmt.Println("Everyone you follow is following you back!")
	}
}
