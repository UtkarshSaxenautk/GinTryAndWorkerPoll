package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type wordCountRequest struct {
	Str string `json:"str"`
}

type breed struct {
	Breed   string `json:"breed"`
	Country string `json:"country"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

type catBreedsResponse struct {
	CurrentPage int     `json:"current_page"`
	Data        []breed `json:"data"`
	LastPage    int     `json:"last_page"`
	Total       int     `json:"total"`
}

type responseBreed struct {
	Breed   string `json:"breed"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

type catBreedsByCountry map[string][]responseBreed

func main() {
	r := gin.Default()

	r.GET("/cat-breeds", func(c *gin.Context) {
		// Send GET request to retrieve data from the URL
		resp, err := http.Get("https://catfact.ninja/breeds")
		if err != nil {
			log.Println("Error making request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
			return
		}
		defer resp.Body.Close()

		// Read the response body
		resposnseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Log the response AS-IS to a text file
		addLogsToFile(resposnseBody)

		// Parse the response body into CatBreedsResponse struct
		var catBreedsRes catBreedsResponse
		err = json.Unmarshal(resposnseBody, &catBreedsRes)
		if err != nil {
			log.Println("Error parsing response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		// Console log the number of pages of data available
		numPages := catBreedsRes.LastPage
		log.Println("Total Number of pages :", numPages)

		// Fetch data from all pages
		allBreeds := catBreedsRes.Data
		for i := 2; i <= numPages; i++ {
			// Send GET request to retrieve data from the next page
			pageURL := "https://catfact.ninja/breeds?page=" + strconv.Itoa(i)
			resp, err := http.Get(pageURL)
			if err != nil {
				log.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			// Read the response body of the next page
			pageBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response:", err)
				continue
			}

			// Parse the response body of the next page into CatBreedsResponse struct
			var nextPageResponse catBreedsResponse
			err = json.Unmarshal(pageBody, &nextPageResponse)
			if err != nil {
				log.Println("Error parsing response:", err)
				continue
			}

			// Append breeds from the next page to the existing breeds
			allBreeds = append(allBreeds, nextPageResponse.Data...)
		}

		// Group cat breeds by country
		catBreedsByCountry := groupBreedsByCountry(allBreeds)

		c.JSON(http.StatusOK, catBreedsByCountry)
	})

	r.POST("/", func(c *gin.Context) {
		// Bind the request payload to WordCountRequest struct
		var req wordCountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		// Count the number of words (sets of characters) in the string
		wordCount := countWords(req.Str)

		// Check if the word count is at least 8
		if wordCount >= 8 {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Not Acceptable as words are less than 8"})
		}
	})

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln("error in starting server at 8080")
		return
	}
}

//helpers functions

// Helper function to count words from given string
func countWords(str string) int {
	// Use regex to split the string into words (sets of characters)
	words := regexp.MustCompile(`\w+`).FindAllString(str, -1)

	// Count the number of words
	wordCount := len(words)
	return wordCount
}

// Helper function to log the response to a sort breeds by country
func groupBreedsByCountry(breeds []breed) catBreedsByCountry {
	catBreedsByCountries := make(catBreedsByCountry)
	for _, breed := range breeds {
		if breed.Country == "" {
			continue
		}
		resBreed := responseBreed{
			Breed:   breed.Breed,
			Origin:  breed.Breed,
			Coat:    breed.Coat,
			Pattern: breed.Pattern,
		}
		catBreedsByCountries[breed.Country] = append(catBreedsByCountries[breed.Country], resBreed)
	}
	return catBreedsByCountries
}

// Helper function to log the response to a  response text file
func addLogsToFile(body []byte) {
	filePath := "response.txt"
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		log.Println("Error writing response to file:", err)
		return
	}
	log.Println("Response logged to:", filePath)
}
