package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// 1. CONFIG
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiKey = strings.TrimSpace(apiKey)
	proxyAddr := "http://127.0.0.1:7897" // Confirm your proxy port

	if apiKey == "" {
		log.Fatal("❌ Error: GEMINI_API_KEY is not set.")
	}

	// 2. SETUP PROXY TRANSPORT
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 3. CONSTRUCT RAW REQUEST
	// We hit the endpoint directly. No libraries.
	// Note the ?key= query parameter. This is the explicit authentication method.
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	jsonBody := []byte(`{
		"contents": [{
			"parts": [{"text": "If you receive this, reply with 'SYSTEM_CHECK_PASSED'"}]
		}]
	}`)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("--- SYSTEM DIAGNOSTIC ---")
	fmt.Printf("1. Proxy Target: %s\n", proxyAddr)
	fmt.Printf("2. Target URL: https://generativelanguage.googleapis.com/... (Key Hidden)\n")
	fmt.Println("3. Sending Raw HTTP Request...")

	// 4. EXECUTE
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n❌ NETWORK ERROR: %v\n", err)
		fmt.Println("Suggestion: Check if your VPN/Proxy is running and allows LAN connections.")
		return
	}
	defer resp.Body.Close()

	// 5. ANALYZE RESPONSE
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\nStatus Code: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode == 200 {
		fmt.Println("\n✅ DIAGNOSIS: Success! Your Key and Proxy are working.")
		fmt.Println("Result: The issue lies within the LangChainGo configuration.")
	} else if resp.StatusCode == 403 {
		fmt.Println("\n❌ DIAGNOSIS: Authorization Failed.")
		fmt.Println("Reason: The server rejected the key. If the key is new, wait 5 mins or check 'Google AI Studio' vs 'Vertex AI' settings.")
	} else {
		fmt.Println("\n⚠️ DIAGNOSIS: Unknown Error. Read the body above.")
	}
}
