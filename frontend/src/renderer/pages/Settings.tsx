import { Check, Copy, RotateCw } from "lucide-react";
import React, { useState } from "react";
import type { Settings as SettingsType, ToastContextValue } from "../../types";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { Input } from "../components/ui/input";
import { Select } from "../components/ui/select";
import { useToast } from "../components/ui/toast-container";
import { useSettings, useUpdateSettings } from "../queries";

function Settings() {
  const { data: settings, isLoading } = useSettings();
  const updateSettingsMutation = useUpdateSettings();
  const { showSuccess, showError } = useToast() as ToastContextValue;

  const [localSettings, setLocalSettings] = useState<SettingsType | null>(null);
  const [copiedKey, setCopiedKey] = useState(false);

  // Update local settings when data loads
  React.useEffect(() => {
    if (settings) {
      setLocalSettings(settings);
    }
  }, [settings]);

  const handleSave = async () => {
    if (!localSettings) return;

    updateSettingsMutation.mutate(localSettings, {
      onSuccess: () => {
        showSuccess("Settings saved successfully!");
      },
      onError: () => {
        showError("Failed to save settings");
      },
    });
  };

  const handleChange = <K extends keyof SettingsType>(
    field: K,
    value: SettingsType[K]
  ) => {
    setLocalSettings(
      (prev) => (prev ? { ...prev, [field]: value } : null) as SettingsType
    );
  };

  const generateNewAPIKey = () => {
    const newKey = Array.from({ length: 32 }, () =>
      Math.floor(Math.random() * 16).toString(16)
    ).join("");
    handleChange("api_key", newKey as SettingsType["api_key"]);
    showSuccess("New API key generated");
  };

  const copyAPIKey = () => {
    if (localSettings?.api_key) {
      navigator.clipboard.writeText(localSettings.api_key);
      setCopiedKey(true);
      showSuccess("API key copied to clipboard");
      setTimeout(() => setCopiedKey(false), 2000);
    }
  };

  if (isLoading || !localSettings) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-lg font-bold">Loading...</div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-4xl mx-auto animate-fade-in">
      <div className="mb-6">
        <h2 className="text-2xl font-black mb-1">Settings</h2>
        <p className="text-sm text-gray-500 font-medium">
          Configure proxy server and application preferences
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="text-lg">Proxy Configuration</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
                Proxy Port
              </label>
              <Input
                type="number"
                value={localSettings?.port || 8081}
                onChange={(e) => handleChange("port", parseInt(e.target.value))}
                min="1024"
                max="65535"
              />
              <p className="text-[10px] text-gray-400 font-bold uppercase mt-2">
                Port number for the proxy server (default: 8081)
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="text-lg">Routing Strategy</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
                Strategy
              </label>
              <Select
                value={localSettings?.routing_strategy || "round_robin"}
                onChange={(e) =>
                  handleChange("routing_strategy", e.target.value)
                }
              >
                <option value="round_robin">Round Robin</option>
                <option value="fill_first">Fill First</option>
              </Select>
              <p className="text-[10px] text-gray-400 font-bold uppercase mt-2">
                Round Robin: Distribute requests evenly across accounts
                <br />
                Fill First: Use accounts until quota exhausted, then move to
                next
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="text-lg">Startup</CardTitle>
          </CardHeader>
          <CardContent>
            <label className="flex items-center gap-3 cursor-pointer group">
              <input
                type="checkbox"
                checked={localSettings?.auto_start || false}
                onChange={(e) => handleChange("auto_start", e.target.checked)}
                className="w-5 h-5 border-2 border-black rounded-base cursor-pointer accent-black"
              />
              <span className="text-sm font-bold text-gray-700 group-hover:text-black transition-colors">
                Auto-start proxy when application launches
              </span>
            </label>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="text-lg">Proxy API Key</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="block text-[10px] font-bold text-gray-500 uppercase tracking-wider mb-2">
                API Key
              </label>
              <div className="flex gap-2 mb-3">
                <Input
                  type="text"
                  value={localSettings?.api_key || ""}
                  onChange={(e) => handleChange("api_key", e.target.value)}
                  placeholder="Generate or enter API key for proxy authentication"
                />
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={copyAPIKey}
                  disabled={!localSettings?.api_key}
                  title="Copy API key"
                >
                  {copiedKey ? (
                    <Check className="w-4 h-4" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </Button>
              </div>
              <p className="text-[10px] text-gray-400 font-bold uppercase mb-4">
                API key required for clients to authenticate with the proxy
                server. Leave empty to allow any client.
              </p>
              <Button
                size="sm"
                variant="secondary"
                onClick={generateNewAPIKey}
                className="flex items-center gap-2"
              >
                <RotateCw className="w-3.5 h-3.5" />
                Generate New Key
              </Button>
            </div>
          </CardContent>
        </Card>

        <div className="flex justify-end">
          <Button
            variant="primary"
            onClick={handleSave}
            disabled={updateSettingsMutation.isPending}
          >
            {updateSettingsMutation.isPending ? "Saving..." : "Save Settings"}
          </Button>
        </div>
      </div>
    </div>
  );
}

export default Settings;
