import { createFileRoute } from "@tanstack/react-router";
import { Copy, Key, RefreshCw, Save, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import type { ToastContextValue } from "../../types";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { useToast } from "../components/ui/toast-container";
import { useSettings, useUpdateSettings } from "../queries";

export const Route = createFileRoute("/api-keys")({
  component: ApiKeys,
});

function ApiKeys() {
  const { data: settings, isLoading } = useSettings();
  const updateSettings = useUpdateSettings();
  const { showSuccess, showError } = useToast() as ToastContextValue;
  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    if (settings) {
      setApiKey(settings.api_key || "");
    }
  }, [settings]);

  const generateKey = () => {
    const chars =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    const newKey = Array.from(
      { length: 32 },
      () => chars[Math.floor(Math.random() * chars.length)]
    ).join("");
    setApiKey(newKey);
  };

  const handleSave = () => {
    if (!settings) return;

    updateSettings.mutate(
      { ...settings, api_key: apiKey },
      {
        onSuccess: () => {
          showSuccess("API Key updated successfully");
        },
        onError: () => {
          showError("Failed to update API Key");
        },
      }
    );
  };

  const handleClear = () => {
    setApiKey("");
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(apiKey);
    showSuccess("Copied to clipboard");
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-lg font-bold">Loading...</div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-4xl mx-auto animate-fade-in">
      <div className="mb-6">
        <h2 className="text-2xl font-black mb-1">API Keys</h2>
        <p className="text-sm text-gray-500 font-medium">
          Manage authentication keys for your AI proxy
        </p>
      </div>

      <Card>
        <CardHeader className="mb-2">
          <CardTitle className="text-lg flex items-center gap-2">
            <Key className="w-5 h-5" />
            Proxy API Key
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-3 space-y-6">
          <div className="p-4 bg-blue-50 border border-blue-200 rounded-base mb-4">
            <p className="text-xs text-blue-800 leading-relaxed font-medium">
              This key is used by CLI agents and other clients to authenticate
              with your local proxy server. If set, clients must include the{" "}
              <code className="bg-white/50 px-1 py-0.5 rounded">
                Authorization: Bearer YOUR_KEY
              </code>{" "}
              header. Leave empty to disable authentication.
            </p>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
                Your Secret API Key
              </label>
              <div className="flex gap-2">
                <input
                  type="text"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="No API key set (authentication disabled)"
                  className="flex-1 bg-surface px-3 py-2 border-2 border-black rounded-base font-mono text-sm focus:outline-none focus:ring-2 focus:ring-black"
                />
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={copyToClipboard}
                  title="Copy to clipboard"
                >
                  <Copy className="w-4 h-4" />
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={generateKey}
                  title="Generate new key"
                >
                  <RefreshCw className="w-4 h-4" />
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={handleClear}
                  title="Clear key"
                >
                  <Trash2 className="w-4 h-4 text-red-500" />
                </Button>
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <Button
                variant="primary"
                onClick={handleSave}
                disabled={updateSettings.isPending}
                className="flex items-center gap-2"
              >
                <Save className="w-4 h-4" />
                Save Changes
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="mt-8">
        <h3 className="text-base font-black mb-4 uppercase tracking-wider">
          How to use
        </h3>
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-4">
              <div>
                <p className="text-sm font-bold mb-2">For cURL requests:</p>
                <code className="block p-3 bg-surface border-2 border-black rounded-base text-xs font-mono">
                  curl http://localhost:8081/v1/chat/completions \<br />
                  &nbsp;&nbsp;-H "Authorization: Bearer {apiKey || "YOUR_KEY"}"
                  \<br />
                  &nbsp;&nbsp;-H "Content-Type: application/json" \<br />
                  &nbsp;&nbsp;-d '{"{"}"model": "gpt-4o", "messages": [{"{"}
                  "role": "user", "content": "Hello!"{"}"}]'
                </code>
              </div>
              <div>
                <p className="text-sm font-bold mb-2">
                  For OpenAI SDK (Node.js):
                </p>
                <code className="block p-3 bg-surface border-2 border-black rounded-base text-xs font-mono">
                  const openai = new OpenAI({"{"}
                  <br />
                  &nbsp;&nbsp;apiKey: "{apiKey || "YOUR_KEY"}",
                  <br />
                  &nbsp;&nbsp;baseURL: "http://localhost:8081/v1"
                  <br />
                  {"}"});
                </code>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
