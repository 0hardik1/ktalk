# ktalk

ktalk is a kubectl plugin that converts natural language descriptions into kubectl commands using OpenAI's GPT-4. It makes interacting with your Kubernetes cluster more intuitive by allowing you to describe what you want to do in plain English.

## Features

- Natural language to kubectl command conversion
- Interactive command execution with confirmation
- Powered by OpenAI's GPT-4 model
- Safe command execution with user confirmation
- Support for complex kubectl queries

## Prerequisites

- Go 1.16 or later
- kubectl installed and configured
- OpenAI API key

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

# For Windows (PowerShell as Administrator)
Copy-Item kubectl-ktalk C:\Windows\System32\
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
2. ktalk uses OpenAI's GPT-4 to convert your description into the appropriate kubectl command
3. The generated command is shown to you for confirmation
4. Upon confirmation (pressing Enter), the command is executed

## Safety Features

- All generated commands are shown to you before execution
- You must explicitly confirm by pressing Enter to execute any command
- The plugin uses GPT-4's understanding of Kubernetes to generate safe and appropriate commands

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- OpenAI for providing the GPT-4 API
- The Kubernetes community for kubectl and its plugin architecture
