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

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Error(ErrorFaileduUnmarshalTheReleaseData, err)
		return false, "", err
	}

	isLatest := currentVersion == release.TagName
	return isLatest, release.TagName, nil
}
