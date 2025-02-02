package main

import (
    "testing"
)


// Test generated using Keploy
func TestFetchPaginatedData_EmptyResponse(t *testing.T) {
    mockFetchGitHubData := func(url, token string) ([]User, error) {
        return []User{}, nil
    }

    originalFetchGitHubData := fetchGitHubData
    fetchGitHubData = mockFetchGitHubData
    defer func() { fetchGitHubData = originalFetchGitHubData }()

    users, err := fetchPaginatedData(baseFollowingURL, "mockToken")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if len(users) != 0 {
        t.Errorf("Expected 0 users, got %d", len(users))
    }
}
