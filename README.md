# Eino ChatModel Practical Example

This is a practical example demonstrating how to use the CloudWeGo Eino framework to create a ChatModel-based AI assistant with tool calling capabilities.

## Overview

This repository contains two examples:

### Basic Example (`main.go`)
- Simple weather assistant with one tool
- Demonstrates basic ChatModelAgent setup
- Shows streaming response handling

### Advanced Example (`advanced_example.go`)
- Multi-tool city guide assistant
- Includes detailed weather and city information tools
- Demonstrates callbacks for logging streaming events
- Shows complex multi-tool coordination

Both examples demonstrate key Eino features:
- **ChatModelAgent**: Agents that can call tools automatically
- **Tool Creation**: Custom tools with structured input/output
- **Streaming**: Real-time response handling
- **Callbacks**: Event-driven monitoring and logging
- **Context Management**: Maintaining conversation context

## Prerequisites

- Go 1.21 or higher
- OpenAI API key (or compatible LLM provider)
- Basic understanding of Go programming

## Setup

1. Clone this repository:
   ```bash
   git clone https://github.com/SmartsYoung/eino-chatmodel-example.git
   cd eino-chatmodel-example
   ```

2. Set your OpenAI API key as an environment variable:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```
   
   Or copy the example env file:
   ```bash
   cp .env.example .env
   # Edit .env with your actual credentials
   source .env
   ```

3. Run the examples:

   **Basic Example:**
   ```bash
   go run main.go
   ```

   **Advanced Example:**
   ```bash
   go run advanced_example.go advanced
   ```

## Code Structure

### Basic Example (`main.go`)
- **Weather Tool**: Simple tool with basic weather data
- **Single Agent**: Configured with one tool and basic instructions
- **Simple Streaming**: Basic event processing

### Advanced Example (`advanced_example.go`)
- **Multiple Tools**: Weather tool + City info tool
- **Enhanced Data**: Detailed weather (temp, humidity, wind) and city info (population, country, language)
- **Callbacks**: Streaming event logging with custom handlers
- **Complex Queries**: Handles multi-aspect questions requiring multiple tools

## Key Features Demonstrated

### 1. Tool Calling Architecture
```go
// Create invokable tools
tool, err := tool.NewInvokableTool(
    "tool_name",
    "Tool description",
    yourFunction, // Function with structured input/output
)

// Configure agent with tools
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Model: chatModel,
    ToolsConfig: adk.ToolsConfig{
        ToolsNodeConfig: compose.ToolsNodeConfig{
            Tools: []tool.BaseTool{tool},
        },
    },
})
```

### 2. Streaming Response Handling
```go
// Enable streaming in runner
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           agent,
    EnableStreaming: true,
})

// Process streaming events
iter := runner.Query(ctx, "Your question here")
for {
    event, ok := iter.Next()
    if !ok { break }
    if event.Message != nil && event.Message.Content != "" {
        fmt.Print(event.Message.Content) // Real-time output
    }
}
```

### 3. Callbacks for Monitoring
```go
// Create streaming callbacks
callback := cbutils.NewHandlerHelper().ChatModel(&cbutils.ModelCallbackHandler{
    OnStartWithStreamInput: func(ctx context.Context, info *callbacks.RunInfo, input *tool.StreamingInput) context.Context {
        fmt.Printf("🚀 Starting streaming request...\n")
        return ctx
    },
    OnEndWithStreamOutput: func(ctx context.Context, info *callbacks.RunInfo, output *tool.StreamingOutput) {
        fmt.Printf("✅ Streaming completed\n")
    },
}).Handler()

// Add to agent configuration
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Callbacks: []callbacks.Handler{callback},
    // ... other config
})
```

## Customization

### Change LLM Provider
Modify the ChatModel initialization:
```go
// For OpenAI (default)
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model:  "gpt-4o",
    APIKey: os.Getenv("OPENAI_API_KEY"),
})

// For Azure OpenAI
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    APIKey:  os.Getenv("OPENAI_API_KEY"),
    Model:   os.Getenv("OPENAI_MODEL"),
    BaseURL: os.Getenv("OPENAI_BASE_URL"),
    ByAzure: true,
})

// For Ark (Volcano Engine) - see .env.example
```

### Add More Tools
Create additional tools following the existing patterns:
```go
func createYourTool() tool.InvokableTool {
    tool, err := tool.NewInvokableTool(
        "your_tool_name",
        "Tool description with specific capabilities",
        yourFunction, // Must have structured input/output types
    )
    // ... error handling
    return tool
}
```

### Modify Agent Behavior
Update the `Instruction` field to change agent personality and capabilities:
```go
Instruction: `You are a [role]. 
Use these tools: [tool descriptions].
Follow these guidelines: [specific instructions].`,
```

## Learning Resources

- [CloudWeGo Eino Documentation](https://www.cloudwego.io/docs/eino/)
- [Eino Examples Repository](https://github.com/cloudwego/eino-examples)
- [Eino Extensions](https://github.com/cloudwego/eino-ext)
- [Official Quick Start Guide](https://www.cloudwego.io/docs/eino/quick_start/)

## License

This example is provided under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

---

**Note**: These examples use mock data for demonstration purposes. In production applications:
- Replace mock functions with real API integrations
- Add proper error handling and retry logic
- Implement authentication and security best practices
- Consider rate limiting and cost management