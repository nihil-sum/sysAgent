package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// --- 1. PROXY & AUTH CONFIG (Reusable Logic) ---
type authTransport struct {
	transport http.RoundTripper
	apiKey    string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-goog-api-key", t.apiKey)
	return t.transport.RoundTrip(req)
}

func setupLLM() llms.Model {
	apiKey := os.Getenv("GEMINI_API_KEY")
	// Proxy Setup (Port 7897 based on your success)
	proxyUrl, _ := url.Parse("http://127.0.0.1:7897")

	httpClient := &http.Client{
		Transport: &authTransport{
			transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
			apiKey:    apiKey,
		},
	}

	ctx := context.Background()
	llm, err := googleai.New(ctx,
		googleai.WithAPIKey(apiKey),
		googleai.WithHTTPClient(httpClient),
		googleai.WithDefaultModel("gemini-2.5-flash-lite"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return llm
}

// --- 2. API HANDLERS ---

type ChatRequest struct {
	Message string `json:"message"`
}

func main() {
	// Initialize the Brain
	agentLLM := setupLLM()

	// Initialize the Server
	r := gin.Default()

	// CORS is critical for Next.js communication
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your Next.js URL
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Define Routes
	r.POST("/api/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prompt Engineering: Give the agent a System Persona
		systemPrompt := "You are SysAgent, a high-performance system optimizer. Answer briefly and technically."
		finalPrompt := fmt.Sprintf("%s\nUser: %s", systemPrompt, req.Message)

		ctx := context.Background()
		response, err := llms.GenerateFromSinglePrompt(ctx, agentLLM, finalPrompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Agent crashed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": response,
			"status":   "active",
		})
	})

	// Start on Port 8080
	log.Println("🚀 SysAgent Core online at http://localhost:8080")
	r.Run(":8080")
}
