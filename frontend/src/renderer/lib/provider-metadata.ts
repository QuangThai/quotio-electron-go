export const PROVIDER_METADATA: Record<
  string,
  { icon: string; label: string; description?: string }
> = {
  claude: { icon: '🤖', label: 'Claude Code', description: 'Anthropic Claude models' },
  openai: { icon: '⚡', label: 'OpenAI (Codex)', description: 'ChatGPT and Codex models' },
  gemini: { icon: '💎', label: 'Gemini CLI', description: 'Google Gemini models' },
  antigravity: { icon: '🚀', label: 'Antigravity', description: 'Unified Gateway (Claude via Gemini)' },
  copilot: { icon: '💻', label: 'GitHub Copilot' },
  qwen: { icon: '🌟', label: 'Qwen Code' },
  vertex: { icon: '☁️', label: 'Vertex AI' },
  iflow: { icon: '🌊', label: 'iFlow' },
  kiro: { icon: '✨', label: 'Kiro (CodeWhisperer)' },
  cursor: { icon: '🖱️', label: 'Cursor' },
  ampcode: { icon: '🧬', label: 'Ampcode' },
  'z.ai': { icon: '🧠', label: 'Z.AI' },
};

export const getProviderIcon = (provider: string): string =>
  PROVIDER_METADATA[provider]?.icon ?? '🔌';

export const getProviderLabel = (provider: string): string =>
  PROVIDER_METADATA[provider]?.label ?? provider;

export const getProviderDescription = (provider: string): string | undefined =>
  PROVIDER_METADATA[provider]?.description;
