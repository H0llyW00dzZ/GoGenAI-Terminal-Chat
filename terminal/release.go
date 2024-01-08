// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GitHubRelease represents the structure of a release as returned by the GitHub API.
type GitHubRelease struct {
	TagName string `json:"tag_name"` // The name of the tag for this release
	Name    string `json:"name"`     // The name of the release
	Body    string `json:"body"`     // The body of the release, typically includes changelog
}

// checkLatestVersion fetches the latest release from the GitHub repository and compares it
// with the current version of the application.
//
// Parameters:
//
//	currentVersion string: The current version of the application.
//
// Returns:
//
//	bool: Indicates whether the current version is the latest.
//	string: The latest version tag name if a newer version is available.
//	error: An error if the request to the GitHub API fails or parsing fails.
func checkLatestVersion(currentVersion string) (bool, string, error) {
	resp, err := http.Get(GitHubAPIURL)
	if err != nil {
		logger.Error(ErrorFailedToFetchReleaseInfo, err)
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf(ErrorReceivedNon200StatusCode, resp.StatusCode)
		logger.Error(errMsg)
		return false, "", fmt.Errorf(errMsg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ErrorFailedToReadTheResponseBody, err)
		return false, "", err
	}

	if err := json.Unmarshal(body, &checkVersion); err != nil {
		logger.Error(ErrorFaileduUnmarshalTheReleaseData, err)
		return false, "", err
	}

	isLatest := currentVersion == checkVersion.TagName
	return isLatest, checkVersion.TagName, nil
}

// getFullReleaseInfo fetches the full release info for the given tag name from the GitHub API.
func getFullReleaseInfo(tagName string) (*GitHubRelease, error) {
	releaseURL := fmt.Sprintf(GitHubReleaseFUll, tagName)

	resp, err := http.Get(releaseURL)
	if err != nil {
		logger.Error(ErrorFailedTagToFetchReleaseInfo, tagName, err)
		return nil, err // Return the original error without additional formatting
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf(ErrorReceivedNon200StatusCode, resp.StatusCode)
		logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := json.NewDecoder(resp.Body).Decode(&checkVersion); err != nil {
		logger.Error(ErrorFailedTagUnmarshalTheReleaseData, tagName, err)
		return nil, err // Return the original error without additional formatting
	}

	return &checkVersion, nil
}
