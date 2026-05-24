package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

// WeatherInput represents the input for the weather tool
type WeatherInput struct {
	City string `json:"city" jsonschema_description:"The name of the city to get weather for"`
}

// WeatherOutput represents the output from the weather tool
type WeatherOutput struct {
	Temperature int    `json:"temperature"`
	Condition   string `json:"condition"`
}

// getWeather simulates a weather API call
func getWeather(ctx context.Context, input *WeatherInput) (*WeatherOutput, error) {
	// In a real implementation, this would call a weather API
	// For this example, we'll return mock data
	city := input.City
	
	// Mock weather data based on city
	weatherData := map[string]WeatherOutput{
		"beijing":   {Temperature: 25, Condition: "Sunny"},
		"shanghai":  {Temperature: 28, Condition: "Partly Cloudy"},
		"guangzhou": {Temperature: 32, Condition: "Hot and Humid"},
		"shenzhen":  {Temperature: 30, Condition: "Thunderstorms"},
		"hangzhou":  {Temperature: 26, Condition: "Rainy"},
	}
	
	if data, exists := weatherData[city]; exists {
		return &data, nil
	}
	
	// Default response for unknown cities
	return &WeatherOutput{Temperature: 20, Condition: "Unknown"}, nil
}

// createWeatherTool creates a tool for getting weather information
func createWeatherTool() tool.InvokableTool {
	tool, err := tool.NewInvokableTool(
		"get_weather",
		"Get current weather information for a specific city",
		getWeather,
	)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}
	return tool
}

func main() {
	ctx := context.Background()
	
	// Initialize OpenAI ChatModel
	// Make sure to set OPENAI_API_KEY environment variable
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-4o", // or "gpt-3.5-turbo"
		APIKey: os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create chat model: %v", err)
	}
	
	// Create a ChatModelAgent with the weather tool
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "WeatherAssistant",
		Description: "An AI assistant that can provide weather information",
		Instruction: `You are a helpful weather assistant. 
When users ask about weather, use the "get_weather" tool to fetch current weather information.
Always provide temperature in Celsius and describe the weather condition clearly.
If the user doesn't specify a city, ask them to provide one.`,
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{createWeatherTool()},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	
	// Create a runner to execute the agent
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true, // Enable streaming responses
	})
	
	// Query the agent with a weather question
	fmt.Println("Asking for weather in Beijing...")
	iter := runner.Query(ctx, "What's the weather like in Beijing today?")
	
	// Process the streaming response
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("Error: %v", event.Err)
			continue
		}
		
		// Print the content from the event
		if event.Message != nil && event.Message.Content != "" {
			fmt.Print(event.Message.Content)
		}
	}
	fmt.Println("\n")
	
	// Example with a different city
	fmt.Println("Asking for weather in Shanghai...")
	iter = runner.Query(ctx, "How's the weather in Shanghai?")
	
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("Error: %v", event.Err)
			continue
		}
		
		if event.Message != nil && event.Message.Content != "" {
			fmt.Print(event.Message.Content)
		}
	}
	fmt.Println("\n")
	
	// Example with ambiguous request (should trigger clarification)
	fmt.Println("Asking without specifying city...")
	iter = runner.Query(ctx, "What's the weather like?")
	
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("Error: %v", event.Err)
			continue
		}
		
		if event.Message != nil && event.Message.Content != "" {
			fmt.Print(event.Message.Content)
		}
	}
	fmt.Println()
}