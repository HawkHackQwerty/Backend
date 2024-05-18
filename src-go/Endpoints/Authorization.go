package Endpoints

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://zdgvghjjvbphcovfayyv.supabase.co/"
	apiKey  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InpkZ3ZnaGpqdmJwaGNvdmZheXl2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTYwNTk4NTMsImV4cCI6MjAzMTYzNTg1M30.fbYoUIyJkEhMImMmEgSWAdsLWAqE1-F31s2URGPqtkQ"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

// Helper Functions
func httpPost(url string, payload interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}

// Private Functions for Authentication ________________________________________________________

// Gets Refresh Token from header, to use:
func refreshTokenHandler(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	url := baseURL + "/auth/v1/token?grant_type=refresh_token"
	payload := map[string]string{
		"refresh_token": refreshToken,
	}
	headers := map[string]string{
		"apikey":       apiKey,
		"Content-Type": "application/json",
	}

	resp, err := httpPost(url, payload, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(responseBody)})
		return
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// End of Private Functions for Authentication ________________________________________________________

// Public Authentication ___________________________________________________________________-
func SignUp(c *gin.Context) {
	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	url := baseURL + "/auth/v1/signup"
	headers := map[string]string{
		"apikey":       apiKey,
		"Content-Type": "application/json",
	}

	resp, err := httpPost(url, loginDetails, headers)
	if err != nil || resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(responseBody)})
		return
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signed up successfully", "access_token": authResp.AccessToken})
}

func SignIn(c *gin.Context) {
	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	url := baseURL + "/auth/v1/token?grant_type=password"
	headers := map[string]string{
		"apikey":       apiKey,
		"Content-Type": "application/json",
	}

	resp, err := httpPost(url, loginDetails, headers)
	if err != nil || resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(responseBody)})
		return
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signed in successfully", "access_token": authResp.AccessToken})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// End of Public Authentication ___________________________________________________________________-
