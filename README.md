# ktalk

ktalk is a kubectl plugin that converts natural language descriptions into kubectl commands using OpenAI's GPT-3.5-turbo. It makes interacting with your Kubernetes cluster more intuitive by allowing you to describe what you want to do in plain English.

## Features

- Natural language to kubectl command conversion
- Interactive command execution with confirmation
- Powered by OpenAI's GPT-3.5-turbo model
- Safe command execution with user confirmation
- Support for complex kubectl queries
- Robust error handling and validation
- Command format validation
- HTTP status code checking
- Detailed error messages

## Prerequisites

- Go 1.22.0 or later
- kubectl installed and configured
- OpenAI API key
- Access to a Kubernetes cluster

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ktalk.git
cd ktalk
```

2. Build the plugin:
```bash
go build -o kubectl-ktalk
```

3. Install the plugin:
```bash
# For macOS/Linux
cp kubectl-ktalk /usr/local/bin/
```

4. Set up your OpenAI API key:
```bash
export OPENAI_API_KEY='your-api-key-here'
```

## Usage

Basic syntax:
```bash
kubectl ktalk <your natural language query>
```

### Examples

1. List containers in a namespace:
```bash
kubectl ktalk show me all containers in the kube-system namespace
```

2. Check pod status:
```bash
kubectl ktalk what's the status of pods in default namespace
```

3. Complex queries:
```bash
kubectl ktalk find all pods that are running as root user across all namespaces
```

## How It Works

1. You provide a natural language description of what you want to do with your Kubernetes cluster
2. ktalk uses OpenAI's GPT-3.5-turbo to convert your description into the appropriate kubectl command
3. The generated command is validated to ensure it:
   - Starts with "kubectl"
   - Is properly formatted
   - Is safe to execute
4. The command is shown to you for confirmation
5. Upon confirmation (pressing Enter), the command is executed

## Safety Features

- All generated commands are shown to you before execution
- You must explicitly confirm by pressing Enter to execute any command
- Commands are validated to ensure they start with "kubectl"
- The plugin uses GPT-3.5-turbo's understanding of Kubernetes to generate safe and appropriate commands
- HTTP status codes are checked for API responses
- Response format is thoroughly validated
- Empty or malformed responses are caught and reported

## Error Handling

The plugin includes comprehensive error handling for:
- Missing or invalid OpenAI API key
- API request failures
- Invalid API responses
- Malformed commands
- Command execution errors
- User input errors

## Technical Details

- Uses OpenAI's GPT-3.5-turbo model for optimal performance and cost
- Implements robust JSON response parsing
- Includes thorough type checking and validation
- Provides detailed error messages for debugging
- Uses Go's standard library for command execution
- Implements proper signal handling for command interruption

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. When contributing:

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## Development

To set up the development environment:

1. Install Go 1.22.0 or later
2. Clone the repository
3. Install dependencies:
```bash
go mod download
```
4. Run tests:
```bash
go test ./...
```
5. Build the plugin:
```bash
go build -o kubectl-ktalk
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- OpenAI for providing the GPT-3.5-turbo API
- The Kubernetes community for kubectl and its plugin architecture
- The Go community for excellent tools and libraries

## Support

If you encounter any issues or have questions:
1. Check the error messages for detailed information
2. Review the documentation
3. Open an issue in the GitHub repository
