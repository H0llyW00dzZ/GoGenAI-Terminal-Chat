// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GitHubRelease represents the metadata of a software release from GitHub.
// It includes information such as the tag name, release name, and a description body,
// typically containing the changelog or release notes.
type GitHubRelease struct {
	TagName string `json:"tag_name"` // The tag associated with the release, e.g., "v1.2.3"
	Name    string `json:"name"`     // The official name of the release
	Body    string `json:"body"`     // Detailed description or changelog for the release
}

// CheckLatestVersion compares the current application version against the latest
// version available on GitHub. It fetches the latest release information from the
// repository specified by GitHubAPIURL and determines if an update is available.
//
// Parameters:
// - currentVersion: The version string of the currently running application.
//
// Returns:
// - isLatest: A boolean indicating if the current version is the latest available.
// - latestVersion: The tag name of the latest release, if newer than current; otherwise, an empty string.
// - err: An error if the request fails or if there is an issue parsing the response.
func CheckLatestVersion(currentVersion string) (isLatest bool, latestVersion string, err error) {
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

	isLatest = currentVersion == checkVersion.TagName
	return isLatest, checkVersion.TagName, nil
}

// GetFullReleaseInfo retrieves detailed information about a specific release from GitHub.
// It constructs the request URL based on the provided tag name and fetches the data
// from the GitHub API.
//
// Parameters:
// - tagName: The name of the tag for which release information is requested.
//
// Returns:
// - release: A pointer to the GitHubRelease struct containing the release information.
// - err: An error if the request fails or if there is an issue parsing the response.
func GetFullReleaseInfo(tagName string) (release *GitHubRelease, err error) {
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
