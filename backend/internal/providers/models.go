package providers

type ModelInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	Description string   `json:"description"`
	ContextWindow int    `json:"context_window"`
}

var SupportedModels = []ModelInfo{
	// Gemini
	{
		ID:            "gemini-2.0-pro-exp",
		Name:          "Gemini 2.0 Pro Experimental",
		Provider:      "gemini",
		Description:   "Google's most capable model for high intelligence tasks.",
		ContextWindow: 2000000,
	},
	{
		ID:            "gemini-2.0-flash",
		Name:          "Gemini 2.0 Flash",
		Provider:      "gemini",
		Description:   "Fast and versatile model from Google.",
		ContextWindow: 1000000,
	},
	{
		ID:            "gemini-1.5-pro",
		Name:          "Gemini 1.5 Pro",
		Provider:      "gemini",
		Description:   "High-intelligence model with a massive context window.",
		ContextWindow: 2000000,
	},
	{
		ID:            "gemini-1.5-flash",
		Name:          "Gemini 1.5 Flash",
		Provider:      "gemini",
		Description:   "Fast, cost-efficient model for scaling.",
		ContextWindow: 1000000,
	},

	// OpenAI
	{
		ID:            "gpt-4o",
		Name:          "GPT-4o",
		Provider:      "openai",
		Description:   "OpenAI's most advanced multimodal model.",
		ContextWindow: 128000,
	},
	{
		ID:            "gpt-4-turbo",
		Name:          "GPT-4 Turbo",
		Provider:      "openai",
		Description:   "High-capability GPT-4 model.",
		ContextWindow: 128000,
	},
	{
		ID:            "gpt-3.5-turbo",
		Name:          "GPT-3.5 Turbo",
		Provider:      "openai",
		Description:   "Fast and reliable model for common tasks.",
		ContextWindow: 16385,
	},

	// Claude
	{
		ID:            "claude-3-5-sonnet-latest",
		Name:          "Claude 3.5 Sonnet",
		Provider:      "claude",
		Description:   "Anthropic's most intelligent model.",
		ContextWindow: 200000,
	},
	{
		ID:            "claude-3-opus-latest",
		Name:          "Claude 3 Opus",
		Provider:      "claude",
		Description:   "The original high-intelligence Claude 3 model.",
		ContextWindow: 200000,
	},
	{
		ID:            "claude-3-haiku-20240307",
		Name:          "Claude 3 Haiku",
		Provider:      "claude",
		Description:   "Fastest and most compact model from Anthropic.",
		ContextWindow: 200000,
	},

	// Qwen
	{
		ID:            "qwen-coder-plus",
		Name:          "Qwen Coder Plus",
		Provider:      "qwen",
		Description:   "Alibaba's high-performance coding model.",
		ContextWindow: 128000,
	},
	{
		ID:            "qwen-coder-turbo",
		Name:          "Qwen Coder Turbo",
		Provider:      "qwen",
		Description:   "Fast coding model from Alibaba.",
		ContextWindow: 128000,
	},

	// Cursor
	{
		ID:            "claude-3-5-sonnet",
		Name:          "Claude 3.5 Sonnet (Cursor)",
		Provider:      "cursor",
		Description:   "Anthropic's Sonnet 3.5 via Cursor.",
		ContextWindow: 200000,
	},
	{
		ID:            "gpt-4o",
		Name:          "GPT-4o (Cursor)",
		Provider:      "cursor",
		Description:   "OpenAI's GPT-4o via Cursor.",
		ContextWindow: 128000,
	},
	{
		ID:            "cursor-small",
		Name:          "Cursor Small",
		Provider:      "cursor",
		Description:   "Cursor's fast, efficient model.",
		ContextWindow: 32000,
	},
}

func GetModelsByProvider(provider string) []ModelInfo {
	var results []ModelInfo
	for _, m := range SupportedModels {
		if m.Provider == provider {
			results = append(results, m)
		}
	}
	return results
}
