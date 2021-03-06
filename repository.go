package main

import (
	"log"
	"net/http"
	"strings"
)

// Repository - Stores the information about the repository
type Repository struct {
	name     string
	service  string
	owner    string
	url      string
	ownerURL string
}

// ParseRequestURL - parses the request URL to get the repository information
func (repo *Repository) ParseRequestURL(url string) (bool, *Error) {

	if len(url) == 0 {
		return false, &Error{"URL cannot be empty!"}
	}

	// Remove the forward slash from the beginning to avoid edge cases
	if url[0] == byte('/') {
		url = url[1:]
	}

	var tokens []string = strings.Split(url, "/")

	switch len(tokens) {
	case 1:
		return false, &Error{"Username and Repository is missing!"}
	case 2:
		return false, &Error{"Repository name is missing!"}
	default:
		log.Printf("Received valid request path - [%s]", url)
	}

	// Store the repository details
	repo.service = tokens[0]
	repo.owner = tokens[1]
	repo.name = tokens[2]

	return true, nil
}

// VerifyRepository - Verifies the repository details with the service
func (repo *Repository) VerifyRepository() (bool, *Error) {

	// Make sure URL is generated before using it
	if len(repo.url) == 0 {
		return false, &Error{"No URLs to work with!"}
	}

	// Search for the repository
	response, err := http.Get(repo.url)
	if err != nil || response.StatusCode != 200 {

		// Repository not found, search for the owner for precise error
		response, err := http.Get(repo.ownerURL)
		if err != nil || response.StatusCode != 200 {
			return false, &Error{"Owner not found!"}
		}

		return false, &Error{"No such public repository found under user!"}
	}

	// repository details are completely valid
	return true, nil
}
