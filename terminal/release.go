// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CheckLatestVersion compares the current application version against the latest
// version available on GitHub. It fetches the latest release information from the
// repository specified by GitHubAPIURL and determines if an update is available.
//
// Parameters:
//
//	currentVersion string: The version string of the currently running application.
//
// Returns:
//
//	isLatest bool: A boolean indicating if the current version is the latest available.
//	latestVersion string: The tag name of the latest release, if newer than current; otherwise, an empty string.
//	err error: An error if the request fails or if there is an issue parsing the response.
func CheckLatestVersion(currentVersion string) (isLatest bool, latestVersion string, err error) {
	// Create an HTTP client with a timeout to prevent hanging requests.
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform an HTTP GET request to the GitHub API URL.
	resp, err := client.Get(GitHubAPIURL)
	if err != nil {
		// Log and return the error if the HTTP request fails.
		logger.Error(ErrorFailedToFetchReleaseInfo, err)
		return false, "", err
	}
	// Ensure the body of the response is closed when the function returns.
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		// Log and return an error if the status code is not 200 OK.
		errMsg := fmt.Sprintf(ErrorReceivedNon200StatusCode, resp.StatusCode)
		logger.Error(errMsg)
		return false, "", fmt.Errorf(errMsg)
	}

	// Decode the JSON response into a GitHubRelease struct.
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		// Log and return the error if JSON unmarshaling fails.
		logger.Error(ErrorFaileduUnmarshalTheReleaseData, err)
		return false, "", err
	}

	// Determine if the current version is the latest by comparing strings.
	isLatest = currentVersion == release.TagName
	// Return the comparison result, the tag name of the latest release, and no error.
	return isLatest, release.TagName, nil
}

// GetFullReleaseInfo retrieves detailed information about a specific release from GitHub.
// It constructs the request URL based on the provided tag name and fetches the data
// from the GitHub API.
//
// Parameters:
//
//	tagName: The name of the tag for which release information is requested.
//
// Returns:
//
//	release *GitHubRelease: A pointer to the GitHubRelease struct containing the release information.
//	err error: An error if the request fails or if there is an issue parsing the response.
func GetFullReleaseInfo(tagName string) (release *GitHubRelease, err error) {
	// Construct the full URL to the GitHub API for the given tag name.
	releaseURL := fmt.Sprintf(GitHubReleaseFUll, tagName)

	// Perform an HTTP GET request to the constructed URL.
	resp, err := http.Get(releaseURL)
	if err != nil {
		// Log and return the error if the HTTP request fails.
		logger.Error(ErrorFailedTagToFetchReleaseInfo, tagName, err)
		return nil, err // Return the original error without additional formatting
	}
	// Ensure the body of the response is closed when the function returns.
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		// Log and return an error if the status code is not 200 OK.
		errMsg := fmt.Sprintf(ErrorReceivedNon200StatusCode, resp.StatusCode)
		logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Decode the JSON response into a GitHubRelease struct.
	var checkVersion GitHubRelease // This should be a local variable, not a global one.
	if err := json.NewDecoder(resp.Body).Decode(&checkVersion); err != nil {
		// Log and return the error if JSON unmarshaling fails.
		logger.Error(ErrorFailedTagUnmarshalTheReleaseData, tagName, err)
		return nil, err // Return the original error without additional formatting
	}

	// Return a pointer to the filled GitHubRelease struct and no error.
	return &checkVersion, nil
}

// checkLatestVersionWithBackoff wraps the CheckLatestVersion call with retry logic.
// It attempts to determine if the current version is the latest and retries on failure
// with exponential backoff.
//
// Returns:
//
//	isLatest bool: A boolean indicating if the current version is the latest available.
//	latestVersion string: The tag name of the latest release, if newer than current; otherwise, an empty string.
//	err error: An error if the request fails after retries or if there is an issue parsing the response.
func checkLatestVersionWithBackoff() (isLatest bool, latestVersion string, err error) {
	// Define a retryable operation with a function that checks the latest version.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Call CheckLatestVersion to compare the current version with the latest release.
			var err error
			isLatest, latestVersion, err = CheckLatestVersion(CurrentVersion)
			// The operation is successful if there is no error.
			return err == nil, err
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardOtherAPIErrorHandler)

	// If an error occurs or the operation is not successful after retries, return the error.
	if err != nil || !success {
		return false, "", err
	}

	// Return the results of the version check.
	return isLatest, latestVersion, nil
}

// fetchAndFormatReleaseInfo retrieves and formats the release information.
// It fetches the release information for the given tag name and formats it for display.
//
// Parameters:
//
//	latestVersion string: The tag name of the latest release.
//
// Returns:
//
//	aiPrompt string: A formatted string containing release information.
//	err error: An error if the request fails after retries or if there is an issue with formatting.
func fetchAndFormatReleaseInfo(latestVersion string) (aiPrompt string, err error) {
	// Fetch the release information with retries in case of transient errors.
	releaseInfo, err := fetchReleaseWithBackoff(latestVersion)
	if err != nil {
		return "", err
	}

	// Format the release date into a more readable format.
	formattedDate, err := formatReleaseDate(releaseInfo.Date)
	if err != nil {
		return "", err
	}
	releaseInfo.Date = formattedDate

	// Format the release information into a prompt to be displayed to the user.
	aiPrompt = formatReleasePrompt(releaseInfo)
	return aiPrompt, nil
}

// fetchReleaseWithBackoff tries to fetch the release information with exponential backoff.
// It attempts to retrieve detailed release information and retries on failure with exponential backoff.
//
// Parameters:
//
//	latestVersion string: The tag name of the latest release.
//
// Returns:
//
//	releaseInfo *GitHubRelease: A pointer to the GitHubRelease struct containing the release information.
//	err error: An error if the request fails after retries or if there is an issue parsing the response.
func fetchReleaseWithBackoff(latestVersion string) (*GitHubRelease, error) {
	var releaseInfo *GitHubRelease

	// Define a retryable operation with a function that fetches the release information.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Call GetFullReleaseInfo to fetch detailed release information for the given tag name.
			var err error
			releaseInfo, err = GetFullReleaseInfo(latestVersion)
			// The operation is successful if there is no error.
			return err == nil, err
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardOtherAPIErrorHandler)

	// If an error occurs or the operation is not successful after retries, return the error.
	if err != nil || !success {
		return nil, err
	}

	// Return the fetched release information.
	return releaseInfo, nil
}

// formatReleaseDate takes a date string and returns it in the desired format.
// It parses a date string and reformats it according to OtherTimeFormat.
//
// Parameters:
//
//	dateStr string: The date string to be reformatted.
//
// Returns:
//
//	string: The date string in the new format.
//	error: An error if the date string cannot be parsed.
func formatReleaseDate(dateStr string) (string, error) {
	// Parse the date string according to RFC3339 format.
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", err
	}
	// Format the time according to the OtherTimeFormat constant.
	return t.Format(OtherTimeFormat), nil
}

// formatReleasePrompt formats the release information into a prompt.
// It constructs a string containing formatted release information based on the provided data.
//
// Parameters:
//
//	releaseInfo *GitHubRelease: A struct containing the release information to be formatted.
//
// Returns:
//
//	string: A formatted string containing release information suitable for display.
func formatReleasePrompt(releaseInfo *GitHubRelease) string {
	// Use the ReleaseNotesPrompt format string to construct the release information prompt.
	return fmt.Sprintf(ReleaseNotesPrompt,
		VersionCommand,
		ApplicationName,
		CurrentVersion,
		releaseInfo.TagName,
		releaseInfo.Name,
		releaseInfo.Date,
		releaseInfo.Body,
	)

}
