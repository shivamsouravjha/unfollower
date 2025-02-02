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
)

type User struct {
    Login string `json:"login"`
}

var fetchGitHubData = func(url, token string) ([]User, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

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

        if len(users) < 100 {
            break
        }

        page++
    }

    return allUsers, nil
}

func main() {
    err := godotenv.Load()
    token := os.Getenv("TOKEN")

    following, err := fetchPaginatedData(baseFollowingURL, token)
    if err != nil {
        fmt.Printf("Error fetching following data: %v\n", err)
        return
    }

    followingMap := make(map[string]bool)
    for _, user := range following {
        followingMap[user.Login] = false
    }

    followers, err := fetchPaginatedData(baseFollowersURL, token)
    if err != nil {
        fmt.Printf("Error fetching followers data: %v\n", err)
        return
    }

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

    if len(notFollowingBack) > 0 {
        fmt.Println("Users you follow who don't follow you back:")
        for _, user := range notFollowingBack {
            fmt.Println(user)
        }
    } else {
        fmt.Println("Everyone you follow is following you back!")
    }
}
