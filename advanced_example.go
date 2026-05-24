package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	cbutils "github.com/cloudwego/eino/utils/callbacks"
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
	Humidity    int    `json:"humidity"`
	WindSpeed   int    `json:"wind_speed"`
}

// getWeather simulates a detailed weather API call
func getWeather(ctx context.Context, input *WeatherInput) (*WeatherOutput, error) {
	time.Sleep(500 * time.Millisecond) // Simulate API delay
	
	city := strings.ToLower(input.City)
	
	// Mock detailed weather data
	weatherData := map[string]WeatherOutput{
		"beijing":   {Temperature: 25, Condition: "Sunny", Humidity: 45, WindSpeed: 12},
		"shanghai":  {Temperature: 28, Condition: "Partly Cloudy", Humidity: 65, WindSpeed: 8},
		"guangzhou": {Temperature: 32, Condition: "Hot and Humid", Humidity: 80, WindSpeed: 5},
		"shenzhen":  {Temperature: 30, Condition: "Thunderstorms", Humidity: 75, WindSpeed: 15},
		"hangzhou":  {Temperature: 26, Condition: "Rainy", Humidity: 70, WindSpeed: 10},
		"chengdu":   {Temperature: 22, Condition: "Foggy", Humidity: 85, WindSpeed: 6},
		"wuhan":     {Temperature: 29, Condition: "Hot", Humidity: 70, WindSpeed: 8},
	}
	
	if data, exists := weatherData[city]; exists {
		return &data, nil
	}
	
	return &WeatherOutput{Temperature: 20, Condition: "Unknown", Humidity: 50, WindSpeed: 10}, nil
}

// CityInfoInput represents input for city information tool
type CityInfoInput struct {
	City string `json:"city" jsonschema_description:"The name of the city to get information about"`
}

// CityInfoOutput represents output from city information tool
type CityInfoOutput struct {
	Population int    `json:"population"`
	Country    string `json:"country"`
	Language   string `json:"language"`
}

// getCityInfo provides basic city information
func getCityInfo(ctx context.Context, input *CityInfoInput) (*CityInfoOutput, error) {
	time.Sleep(300 * time.Millisecond) // Simulate API delay
	
	city := strings.ToLower(input.City)
	
	cityData := map[string]CityInfoOutput{
		"beijing":   {Population: 21540000, Country: "China", Language: "Mandarin"},
		"shanghai":  {Population: 26320000, Country: "China", Language: "Mandarin"},
		"guangzhou": {Population: 15300000, Country: "China", Language: "Mandarin/Cantonese"},
		"shenzhen":  {Population: 17560000, Country: "China", Language: "Mandarin"},
		"hangzhou":  {Population: 11940000, Country: "China", Language: "Mandarin"},
		"tokyo":     {Population: 37400000, Country: "Japan", Language: "Japanese"},
		"seoul":     {Population: 25680000, Country: "South Korea", Language: "Korean"},
	}
	
	if data, exists := cityData[city]; exists {
		return &data, nil
	}
	
	return &CityInfoOutput{Population: 1000000, Country: "Unknown", Language: "Unknown"}, nil
}

// createWeatherTool creates a detailed weather tool
func createWeatherTool() tool.InvokableTool {
	tool, err := tool.NewInvokableTool(
		"get_weather",
		"Get detailed current weather information including temperature, condition, humidity, and wind speed for a specific city",
		getWeather,
	)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}
	return tool
}

// createCityInfoTool creates a city information tool
func createCityInfoTool() tool.InvokableTool {
	tool, err := tool.NewInvokableTool(
		"get_city_info",
		"Get basic information about a city including population, country, and primary language",
		getCityInfo,
	)
	if err != nil {
		log.Fatalf("Failed to create city info tool: %v", err)
	}
	return tool
}

// createStreamingCallback creates a callback to log streaming events
func createStreamingCallback() callbacks.Handler {
	return cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
		OnStartWithStreamInput: func(ctx context.Context, info *callbacks.RunInfo, input *tool.StreamingInput) context.Context {
			fmt.Printf("\n🚀 Starting streaming request to %s...\n", info.Name)
			return ctx
		},
		OnEndWithStreamOutput: func(ctx context.Context, info *callbacks.RunInfo, output *tool.StreamingOutput) {
			fmt.Printf("✅ Streaming completed for %s\n", info.Name)
		},
	}).Handler()
}

func advancedExample() {
	ctx := context.Background()
	
	// Initialize OpenAI ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-4o",
		APIKey: os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create chat model: %v", err)
	}
	
	// Create agent with multiple tools and streaming callback
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "CityGuideAssistant",
		Description: "An AI assistant that provides comprehensive city information including weather and demographics",
		Instruction: `You are a helpful city guide assistant. 
You have access to two tools:
1. get_weather - Get detailed weather information (temperature, condition, humidity, wind speed)
2. get_city_info - Get basic city information (population, country, language)

When users ask about cities, use the appropriate tools to provide comprehensive information.
Always present weather data in a clear, readable format with all available details.
For city information, provide population in millions when possible and mention the primary language.
If a user asks about multiple aspects of a city, use both tools if needed.`,
		Model: chatModel,
		Callbacks: []callbacks.Handler{
			createStreamingCallback(),
		},
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					createWeatherTool(),
					createCityInfoTool(),
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	
	// Example 1: Comprehensive city query
	fmt.Println("=== Example 1: Comprehensive City Query ===")
	iter := runner.Query(ctx, "Tell me about Beijing - what's the weather like and some basic info about the city?")
	
	var fullResponse strings.Builder
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
			fullResponse.WriteString(event.Message.Content)
		}
	}
	fmt.Println("\n")
	
	// Example 2: Multi-city comparison
	fmt.Println("=== Example 2: Multi-City Weather Comparison ===")
	iter = runner.Query(ctx, "Compare the weather in Shanghai and Guangzhou right now.")
	
	fullResponse.Reset()
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
			fullResponse.WriteString(event.Message.Content)
		}
	}
	fmt.Println("\n")
	
	// Example 3: Complex multi-tool query
	fmt.Println("=== Example 3: Complex Multi-Tool Query ===")
	iter = runner.Query(ctx, "I'm planning to visit Hangzhou next week. What should I know about the city and what clothes should I pack based on the current weather?")
	
	fullResponse.Reset()
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
			fullResponse.WriteString(event.Message.Content)
		}
	}
	fmt.Println()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "advanced" {
		advancedExample()
	} else {
		// Run the basic example from main.go
		// This is just a placeholder since we can't include both main functions
		fmt.Println("Run 'go run advanced_example.go advanced' for the advanced example")
		fmt.Println("Run 'go run main.go' for the basic example")
	}
}