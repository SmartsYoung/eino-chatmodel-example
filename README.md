# Eino ChatModel Practical Example

This is a practical example demonstrating how to use the CloudWeGo Eino framework to create a ChatModel-based AI assistant with tool calling capabilities.

## Overview

This example creates a weather assistant that can:
- Answer questions about weather in specific cities
- Handle ambiguous requests by asking for clarification
- Use streaming responses for real-time interaction

The example demonstrates key Eino features:
- **ChatModelAgent**: A simple agent that can call tools
- **Tool Creation**: Custom tool for weather information
- **Streaming**: Real-time response handling
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

3. Run the example:
   ```bash
   go run main.go
   ```

## Code Structure

### Main Components

- **`main.go`**: Contains the complete implementation
- **Weather Tool**: A custom tool that simulates weather API calls
- **ChatModelAgent**: Configured with instructions and the weather tool
- **Runner**: Handles execution and streaming responses

### Key Features Demonstrated

1. **Tool Calling**: The agent automatically decides when to call the weather tool
2. **Streaming Responses**: Real-time output as the LLM generates responses
3. **Error Handling**: Proper error handling for API calls and tool execution
4. **Flexible Configuration**: Easy to modify for different LLM providers

## How It Works

1. The `ChatModelAgent` is initialized with:
   - An OpenAI ChatModel
   - A custom weather tool
   - Instructions for behavior

2. When queried, the agent:
   - Analyzes the user's request
   - Decides whether to call the weather tool
   - Processes the tool response
   - Generates a natural language response

3. The runner handles:
   - Streaming the response in real-time
   - Error handling
   - Event processing

## Customization

You can easily customize this example:

### Change LLM Provider
Modify the ChatModel initialization in `main.go`:
```go
// For OpenAI
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model:  "gpt-4o",
    APIKey: os.Getenv("OPENAI_API_KEY"),
})

// For other providers (when supported)
// chatModel, err := otherProvider.NewChatModel(...)
```

### Add More Tools
Create additional tools following the weather tool pattern:
```go
func createAnotherTool() tool.InvokableTool {
    tool, err := tool.NewInvokableTool(
        "tool_name",
        "Tool description",
        yourFunction,
    )
    // ... error handling
    return tool
}
```

Then add it to the agent configuration:
```go
Tools: []tool.BaseTool{createWeatherTool(), createAnotherTool()},
```

### Modify Agent Behavior
Update the `Instruction` field in the `ChatModelAgentConfig` to change how the agent behaves.

## Learning Resources

- [CloudWeGo Eino Documentation](https://www.cloudwego.io/docs/eino/)
- [Eino Examples Repository](https://github.com/cloudwego/eino-examples)
- [Eino Extensions](https://github.com/cloudwego/eino-ext)

## License

This example is provided under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

---

**Note**: This example uses mock weather data. In a production application, you would integrate with a real weather API service.