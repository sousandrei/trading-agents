# trading-agents

trading-agents is a multi-agent financial trading framework re-implemented in Go, based on the original [TradingAgents](https://github.com/TauricResearch/TradingAgents/).

## Project Structure

The project is organized into the following key packages:

*   **`cmd/agents`**: Contains the main entry point for the application, responsible for parsing configuration and initializing core components.
*   **`internal/agents`**: Houses the logic and prompts for various financial agents, categorized by their roles (e.g., analysts, researchers).
*   **`internal/orchestrator`**: Manages the workflow, calling agents in sequence and passing results to subsequent stages of the analysis.
*   **`internal/tools`**: Provides interfaces and implementations for external services, including LLM clients (Gemini, Ollama) and financial data APIs (Finnhub, Simfin).

## Getting Started

### Prerequisites

*   **Go:** Ensure you have Go (version 1.23.2 or later) installed.
*   **Environment Variables:** Set up the necessary API keys and configuration in your environment.

### Configuration

The application is configured using environment variables, typically managed via an `.envrc` file (which is ignored by Git for security). An example `.envrc` is provided:

```bash
export PORT=3001
export LOG_FORMAT="text"

export LLM_PROVIDER="gemini" # or "ollama"
export LLM_MODEL="gemini-2.5-flash" # or "gemma3:4b" etc...

# Ollama specific
export LLM_API_URL="http://localhost:11434"

# Gemini specific (using Vertex AI backend)
# export GOOGLE_CLOUD_PROJECT="your-gcp-project-id"
# export GOOGLE_CLOUD_LOCATION="your-gcp-region"

export FINNHUB_API_KEY="YOUR_FINNHUB_API_KEY"
export SIMFIN_API_KEY="YOUR_SIMFIN_API_KEY"

export API_TIMEOUT="10s"
export API_CACHE_TTL="8760h"
```

**Note:** Replace `YOUR_FINNHUB_API_KEY` and `YOUR_SIMFIN_API_KEY` with your actual API keys. For Gemini, now we're using the default application credentials.


## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes. Ensure you follow the project's coding standards.