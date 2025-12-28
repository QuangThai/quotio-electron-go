import {
  AlertCircle,
  BookOpen,
  Check,
  Copy,
  RefreshCw,
  Settings,
} from "lucide-react";
import { useState } from "react";
import type { AgentStatus, ToastContextValue } from "../../types";
import { renderAgentStatusBadge } from "../components/shared/status-badges";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { useToast } from "../components/ui/toast-container";
import {
  useAgents,
  useConfigureAgent,
  useRefreshAgents,
  useSettings,
} from "../queries";

interface KnownAgent {
  name: string;
  label: string;
  icon: string;
  description: string;
}

function Agents() {
  const { data: agents = [], isLoading } = useAgents();
  const { data: settings } = useSettings();
  const configureAgent = useConfigureAgent();
  const refreshAgentsMutation = useRefreshAgents();
  const { showSuccess, showError } = useToast() as ToastContextValue;
  const [configuringAgent, setConfiguringAgent] = useState<string | null>(null);
  const [copiedField, setCopiedField] = useState<string | null>(null);

  const proxyPort = settings?.proxy_port ?? settings?.port ?? 8081;
  const proxyURL = `http://localhost:${proxyPort}`;
  const proxyAPIKey = settings?.api_key || "";

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text);
    setCopiedField(field);
    showSuccess("Copied to clipboard");
    setTimeout(() => setCopiedField(null), 2000);
  };

  const handleConfigure = async (agentName: string) => {
    if (!proxyURL || proxyURL === "http://localhost:undefined") {
      showError("Proxy configuration is missing. Check Settings.");
      return;
    }

    setConfiguringAgent(agentName);
    configureAgent.mutate(
      {
        agent_name: agentName,
        proxy_url: proxyURL,
      } as Partial<AgentStatus>,
      {
        onSuccess: () => {
          showSuccess("Agent configured successfully!");
          setConfiguringAgent(null);
        },
        onError: () => {
          showError("Failed to configure agent");
          setConfiguringAgent(null);
        },
      }
    );
  };

  const handleRefresh = async () => {
    refreshAgentsMutation.mutate(undefined, {
      onSuccess: () => {
        showSuccess("Agents refreshed successfully!");
      },
      onError: () => {
        showError("Failed to refresh agents");
      },
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-lg font-bold">Loading...</div>
      </div>
    );
  }

  const knownAgents: KnownAgent[] = [
    {
      name: "claude-code",
      label: "Claude Code",
      icon: "🧠",
      description: "Anthropic's official CLI for Claude models",
    },
    {
      name: "opencode",
      label: "OpenCode",
      icon: "⚡",
      description: "The open source AI coding agent",
    },
    {
      name: "gemini-cli",
      label: "Gemini CLI",
      icon: "💎",
      description: "Google's Gemini CLI for Gemini models",
    },
    {
      name: "droid",
      label: "Factory Droid",
      icon: "🤖",
      description: "Factory's AI coding agent",
    },
    {
      name: "codex",
      label: "Codex CLI",
      icon: "💻",
      description: "OpenAI's Codex CLI for GPT-5 models",
    },
    {
      name: "amp-cli",
      label: "Amp CLI",
      icon: "⚡",
      description: "Amp CLI for AI coding",
    },
  ];

  const installedAgents = knownAgents.filter((agent) => {
    const status = agents.find((a) => a.agent_name === agent.name);
    return status?.installed ?? false;
  });

  const notInstalledAgents = knownAgents.filter((agent) => {
    const status = agents.find((a) => a.agent_name === agent.name);
    return !(status?.installed ?? false);
  });

  const configuredCount = agents.filter((a) => a.configured).length;
  const installedCount = agents.filter((a) => a.installed).length;

  return (
    <div className="p-6 max-w-6xl mx-auto animate-fade-in">
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-black mb-1">AI Agent Setup</h2>
            <p className="text-sm text-gray-500 font-medium">
              Configure CLI agents to use CLIProxyAPI
            </p>
          </div>
          <Button
            variant="secondary"
            size="sm"
            onClick={handleRefresh}
            disabled={refreshAgentsMutation.isPending}
            className="flex items-center gap-2"
          >
            <RefreshCw
              className={`w-3.5 h-3.5 ${refreshAgentsMutation.isPending ? "animate-spin" : ""}`}
            />
            Refresh Agents
          </Button>
        </div>
      </div>

      {/* Proxy Configuration Info */}
      <Card className="mb-6">
        <CardHeader className="mb-2">
          <CardTitle className="text-lg">Proxy Configuration</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 pt-3">
          <div>
            <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
              Proxy URL
            </label>
            <div className="flex gap-2">
              <code className="flex-1 bg-surface p-2.5 rounded-base text-xs font-mono break-all border-2 border-black">
                {proxyURL}/v1
              </code>
              <Button
                size="sm"
                variant="secondary"
                onClick={() => copyToClipboard(`${proxyURL}/v1`, "url")}
                className="shrink-0"
              >
                {copiedField === "url" ? (
                  <Check className="w-3.5 h-3.5" />
                ) : (
                  <Copy className="w-3.5 h-3.5" />
                )}
              </Button>
            </div>
          </div>
          {proxyAPIKey && (
            <div>
              <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
                API Key (for client authentication)
              </label>
              <div className="flex gap-2">
                <code className="flex-1 bg-surface p-2.5 rounded-base text-xs font-mono break-all border-2 border-black">
                  {proxyAPIKey.substring(0, 10)}...
                </code>
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => copyToClipboard(proxyAPIKey, "key")}
                  className="shrink-0"
                >
                  {copiedField === "key" ? (
                    <Check className="w-3.5 h-3.5" />
                  ) : (
                    <Copy className="w-3.5 h-3.5" />
                  )}
                </Button>
              </div>
              <p className="text-[10px] text-gray-400 font-bold uppercase mt-2">
                Clients must send:{" "}
                <code className="bg-surface px-1.5 py-0.5 rounded border border-black/10">
                  Authorization: Bearer {proxyAPIKey.substring(0, 8)}...
                </code>
              </p>
            </div>
          )}
          {!proxyAPIKey && (
            <div className="p-3 bg-blue-50 rounded border border-blue-200">
              <p className="text-xs text-blue-700">
                ℹ️ No API key set. Proxy accepts all clients. Set one in
                Settings to require authentication.
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="flex gap-4 mb-6">
        <Card className="flex-1">
          <CardContent className="pt-5 text-center">
            <div className="text-2xl font-black mb-1">✓ {installedCount}</div>
            <div className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              Installed
            </div>
          </CardContent>
        </Card>
        <Card className="flex-1">
          <CardContent className="pt-5 text-center">
            <div className="text-2xl font-black mb-1">⚙ {configuredCount}</div>
            <div className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              Configured
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="mb-6">
        <h3 className="text-lg font-black mb-4 uppercase tracking-wider">
          Installed
        </h3>
        <div className="space-y-4">
          {installedAgents.map((agent) => {
            const status = agents.find((a) => a.agent_name === agent.name);

            return (
              <Card key={agent.name}>
                <CardContent className="pt-6">
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-4 flex-1">
                      <div className="w-12 h-12 rounded-full bg-gray-100 border-4 border-black flex items-center justify-center text-2xl">
                        {agent.icon}
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <h3 className="text-base font-black">
                            {agent.label}
                          </h3>
                          {status &&
                            renderAgentStatusBadge(
                              status.installed || false,
                              status.configured || false,
                              !!status.validation_error
                            )}
                        </div>
                        <p className="text-sm text-gray-600 mb-2">
                          {agent.description}
                        </p>
                        {status?.config_path && (
                          <p className="text-xs text-gray-500 font-mono mb-1">
                            {status.config_path}
                          </p>
                        )}
                        {status?.validation_error && (
                          <div className="flex items-center gap-2 text-xs text-red-600">
                            <AlertCircle className="w-3 h-3" />
                            <span>{status.validation_error}</span>
                          </div>
                        )}
                        {status?.proxy_url && status?.last_configured && (
                          <p className="text-xs text-gray-500">
                            Last configured:{" "}
                            {new Date(status.last_configured).toLocaleString()}
                          </p>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {status?.config_path_exists && (
                        <Button
                          variant="secondary"
                          size="sm"
                          className="flex items-center gap-2"
                          title="View documentation"
                        >
                          <BookOpen className="w-4 h-4" />
                        </Button>
                      )}
                      <Button
                        variant={status?.configured ? "secondary" : "primary"}
                        size="sm"
                        onClick={() => handleConfigure(agent.name)}
                        disabled={
                          configuringAgent === agent.name &&
                          configureAgent.isPending
                        }
                        className="flex items-center gap-2"
                      >
                        {configuringAgent === agent.name &&
                        configureAgent.isPending ? (
                          <RefreshCw className="w-4 h-4 animate-spin" />
                        ) : status?.configured ? (
                          <>
                            <RefreshCw className="w-4 h-4" />
                            Reconfigure
                          </>
                        ) : (
                          <>
                            <Settings className="w-4 h-4" />
                            Configure
                          </>
                        )}
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>

      {notInstalledAgents.length > 0 && (
        <div className="animate-fade-in delay-100">
          <h3 className="text-lg font-black mb-4 uppercase tracking-wider">
            Not Installed
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {notInstalledAgents.map((agent) => (
              <Card
                key={agent.name}
                className="cursor-pointer hover:shadow-neobrutal transition-all"
              >
                <CardContent className="pt-6 text-center">
                  <div className="text-4xl mb-3">{agent.icon}</div>
                  <div className="font-bold text-sm">{agent.label}</div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

export default Agents;
