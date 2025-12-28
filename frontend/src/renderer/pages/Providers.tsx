import {
  AlertCircle,
  BarChart3,
  CheckCircle,
  Edit2,
  Plus,
  RefreshCw,
  Search,
  Trash2,
  TrendingUp,
  Users as UsersIcon,
  XCircle,
} from "lucide-react";
import { useState } from "react";
import type { Provider, ProviderType, ToastContextValue } from "../../types";
import { renderAccountStatusBadge } from "../components/shared/status-badges";
import { AlertDialog } from "../components/ui/alert-dialog";
import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { StatsCard } from "../components/ui/charts";
import { Input } from "../components/ui/input";
import { Modal } from "../components/ui/modal";
import { Progress } from "../components/ui/progress";
import { Select } from "../components/ui/select";
import { useToast } from "../components/ui/toast-container";
import { getProviderIcon, getProviderLabel } from "../lib/provider-metadata";
import {
  useAddProvider,
  useCheckHealth,
  useDeleteProvider,
  useDetectProviders,
  useProviders,
  useUpdateProvider,
} from "../queries";

interface FormData {
  provider: ProviderType;
  name: string;
  api_key: string;
  quota_limit: number;
}

function Providers() {
  const { data: providers = [], isLoading } = useProviders();
  const [showAddModal, setShowAddModal] = useState(false);
  const [showDetectModal, setShowDetectModal] = useState(false);
  const [detectedProviders, setDetectedProviders] = useState<Provider[]>([]);
  const [validatingProviders, setValidatingProviders] = useState<Set<number>>(
    new Set()
  );
  const [expandedChartId, setExpandedChartId] = useState<number | null>(null);
  const [formData, setFormData] = useState<FormData>({
    provider: "claude",
    name: "",
    api_key: "",
    quota_limit: 0,
  });
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [providerToDelete, setProviderToDelete] = useState<
    string | number | null
  >(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
  const [editQuotaLimit, setEditQuotaLimit] = useState<number>(0);
  const { showSuccess, showError } = useToast() as ToastContextValue;

  const addProviderMutation = useAddProvider();
  const updateProviderMutation = useUpdateProvider();
  const deleteProviderMutation = useDeleteProvider();
  const detectProvidersMutation = useDetectProviders();
  const checkHealthMutation = useCheckHealth();

  // Helper function to render health indicator using provider's own health data
  const renderHealthIndicator = (provider: Provider) => {
    if (provider.is_healthy === undefined) {
      return <AlertCircle className="w-4 h-4 text-gray-400" />;
    }
    return provider.is_healthy ? (
      <CheckCircle className="w-4 h-4 text-green-500" />
    ) : (
      <XCircle className="w-4 h-4 text-red-500" />
    );
  };

  const handleAdd = () => {
    setFormData({
      provider: "claude",
      name: "",
      api_key: "",
      quota_limit: 0,
    });
    setShowAddModal(true);
  };

  const handleDelete = async (id: string | number) => {
    setProviderToDelete(id);
    setShowDeleteDialog(true);
  };

  const confirmDelete = async () => {
    if (providerToDelete) {
      deleteProviderMutation.mutate(String(providerToDelete), {
        onSuccess: () => {
          showSuccess("Provider deleted successfully!");
          setProviderToDelete(null);
        },
        onError: () => {
          showError("Failed to delete provider");
        },
      });
    }
  };

  const handleValidate = (accountId: number) => {
    setValidatingProviders((prev) => new Set(prev).add(accountId));
    checkHealthMutation.mutate(String(accountId), {
      onSuccess: () => {
        showSuccess("Provider validated");
        setValidatingProviders((prev) => {
          const next = new Set(prev);
          next.delete(accountId);
          return next;
        });
      },
      onError: () => {
        showError("Failed to validate provider");
        setValidatingProviders((prev) => {
          const next = new Set(prev);
          next.delete(accountId);
          return next;
        });
      },
    });
  };

  const handleDetect = async () => {
    detectProvidersMutation.mutate(undefined, {
      onSuccess: (data) => {
        setDetectedProviders(data || []);
        setShowDetectModal(true);
      },
      onError: () => {
        showError("Failed to detect providers from environment");
      },
    });
  };

  const handleAddDetected = async (detected: Provider) => {
    addProviderMutation.mutate(
      {
        provider: detected.provider,
        name: detected.name.replace(" (detected)", ""),
        api_key: detected.api_key || "",
        quota_limit: detected.quota_limit || 0,
        auto_detected: true,
      } as unknown as Provider,
      {
        onSuccess: () => {
          showSuccess("Provider added successfully!");
          setShowDetectModal(false);
        },
        onError: () => {
          showError("Failed to add provider");
        },
      }
    );
  };

  const handleEditQuota = (provider: Provider) => {
    setEditingProvider(provider);
    setEditQuotaLimit(provider.quota_limit || 0);
    setShowEditModal(true);
  };

  const handleSaveQuota = () => {
    if (!editingProvider) return;
    updateProviderMutation.mutate(
      {
        id: String(editingProvider.id),
        provider: {
          quota_limit: editQuotaLimit,
          quota_manual: editQuotaLimit > 0,
        },
      },
      {
        onSuccess: () => {
          showSuccess("Quota limit updated successfully!");
          setShowEditModal(false);
          setEditingProvider(null);
        },
        onError: () => {
          showError("Failed to update quota limit");
        },
      }
    );
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-lg font-bold">Loading...</div>
      </div>
    );
  }

  // Group providers by provider type
  const groupedByProvider: Record<string, Provider[]> = providers.reduce(
    (acc, p) => {
      (acc[p.provider] = acc[p.provider] || []).push(p);
      return acc;
    },
    {} as Record<string, Provider[]>
  );

  // Separate connected and auto-detected
  const connectedAccounts = providers.filter((p) => !p.auto_detected);
  const autoDetected = providers.filter((p) => p.auto_detected);

  const totalQuotaUsed = providers.reduce(
    (sum, p) => sum + (p.quota_used || 0),
    0
  );
  const healthyCount = providers.filter((p) => p.is_healthy === true).length;

  return (
    <div className="p-6 max-w-6xl mx-auto animate-fade-in">
      <div className="mb-6">
        <h2 className="text-2xl font-black mb-1">Providers</h2>
        <p className="text-sm text-gray-500 font-medium">
          Manage your AI provider accounts and credentials
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <StatsCard
          label="Total Providers"
          value={providers.length}
          icon={<UsersIcon className="w-4 h-4" />}
          color="blue"
        />
        <StatsCard
          label="Healthy"
          value={healthyCount}
          icon={<CheckCircle className="w-4 h-4" />}
          color="green"
        />
        <StatsCard
          label="Total Tokens Used"
          value={totalQuotaUsed.toLocaleString()}
          icon={<TrendingUp className="w-4 h-4" />}
          color="purple"
        />
        <StatsCard
          label="Provider Types"
          value={Object.keys(groupedByProvider).length}
          icon={<BarChart3 className="w-4 h-4" />}
          color="orange"
        />
      </div>

      <Card className="mb-6">
        <CardHeader className="mb-2">
          <CardTitle className="text-lg">
            Connected Accounts ({connectedAccounts.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-1.5 mb-3">
            {Object.entries(groupedByProvider).map(([provider, accounts]) => (
              <Badge key={provider} variant="secondary">
                {getProviderIcon(provider)} {getProviderLabel(provider)} ×
                {accounts.length}
              </Badge>
            ))}
          </div>
          {connectedAccounts.length > 0 && (
            <p className="text-[10px] text-gray-400 font-medium uppercase">
              Click chart icon to toggle usage details
            </p>
          )}
        </CardContent>
      </Card>

      {autoDetected.length > 0 && (
        <Card className="mb-6">
          <CardHeader className="mb-2">
            <CardTitle className="text-lg">
              Auto-detected ({autoDetected.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {autoDetected.map((provider) => (
                <div
                  key={provider.id}
                  className="flex items-center justify-between p-3 bg-surface border-2 border-black rounded-base"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-xl">
                      {getProviderIcon(provider.provider)}
                    </span>
                    <div>
                      <div className="font-bold text-sm">
                        {provider.name || provider.provider}
                      </div>
                      <div className="text-xs text-gray-500">
                        {provider.provider} • Auto-detected
                      </div>
                    </div>
                  </div>
                  <Button variant="secondary" size="sm">
                    <BarChart3 className="w-4 h-4" />
                  </Button>
                </div>
              ))}
            </div>
            <p className="text-[10px] text-gray-400 font-medium uppercase mt-2">
              Add these providers to your connected accounts
            </p>
          </CardContent>
        </Card>
      )}

      <div className="flex gap-3 mb-6">
        <Button
          variant="secondary"
          onClick={handleDetect}
          disabled={detectProvidersMutation.isPending}
        >
          <Search className="w-4 h-4 mr-2" />
          {detectProvidersMutation.isPending
            ? "Detecting..."
            : "Detect Providers"}
        </Button>
        <Button variant="primary" onClick={handleAdd}>
          <Plus className="w-4 h-4 mr-2" />
          Add Provider
        </Button>
      </div>

      <div className="space-y-4">
        {connectedAccounts.map((provider) => {
          const quotaPercentage =
            provider.quota_limit > 0
              ? (provider.quota_used / provider.quota_limit) * 100
              : 0;
          const progressVariant =
            quotaPercentage >= 80
              ? "danger"
              : quotaPercentage >= 50
                ? "warning"
                : "success";

          return (
            <Card key={provider.id}>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <span className="text-2xl group-hover:scale-110 transition-transform">
                      {getProviderIcon(provider.provider)}
                    </span>
                    <div>
                      <div className="font-black text-base flex items-center gap-2 flex-wrap">
                        {provider.name || provider.provider}
                        {renderAccountStatusBadge(provider.status)}
                        {provider.is_healthy !== undefined && (
                          <Badge
                            variant={provider.is_healthy ? "success" : "danger"}
                          >
                            {provider.is_healthy ? "✓ Verified" : "✗ Invalid"}
                          </Badge>
                        )}
                      </div>
                      <div className="text-[10px] text-gray-400 font-bold uppercase tracking-wider flex items-center gap-2 pt-2">
                        {provider.provider}
                        {provider.is_healthy !== undefined && (
                          <>
                            <span>•</span>
                            <span className="flex items-center gap-1">
                              {renderHealthIndicator(provider)}
                              {provider.response_time_ms &&
                                provider.response_time_ms > 0 && (
                                  <span>{provider.response_time_ms}ms</span>
                                )}
                            </span>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => handleValidate(provider.id)}
                      disabled={validatingProviders.has(provider.id)}
                      title="Validate credentials"
                    >
                      <RefreshCw
                        className={`w-4 h-4 ${validatingProviders.has(provider.id) ? "animate-spin" : ""}`}
                      />
                    </Button>
                    <Button
                      variant={
                        expandedChartId === provider.id
                          ? "primary"
                          : "secondary"
                      }
                      size="sm"
                      onClick={() =>
                        setExpandedChartId(
                          expandedChartId === provider.id ? null : provider.id
                        )
                      }
                      title="Toggle quota chart"
                    >
                      <BarChart3 className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => handleDelete(provider.id)}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
                {provider.quota_limit > 0 && (
                  <div className="mt-3">
                    <Progress
                      value={quotaPercentage}
                      max={100}
                      variant={progressVariant}
                      className="mb-1.5"
                    />
                    <div className="text-[10px] text-gray-500 font-bold uppercase">
                      {provider.quota_used.toLocaleString()} /{" "}
                      {provider.quota_limit.toLocaleString()} tokens
                    </div>
                  </div>
                )}

                {/* Expanded Chart Section */}
                {expandedChartId === provider.id && (
                  <div className="mt-4 p-4 bg-surface border-2 border-black rounded-base animate-slide-up">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="font-black text-xs uppercase tracking-widest flex items-center gap-2">
                        <BarChart3 className="w-4 h-4" />
                        Usage Statistics
                      </h4>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 border-0"
                        onClick={() => setExpandedChartId(null)}
                      >
                        ✕
                      </Button>
                    </div>
                    <div className="grid grid-cols-3 gap-3 mb-4">
                      <div className="text-center p-2 bg-white border-2 border-black rounded-base shadow-neobrutal-sm">
                        <div className="text-lg font-black text-neobrutal-blue">
                          {(provider.quota_used || 0).toLocaleString()}
                        </div>
                        <div className="text-[10px] text-gray-500 font-bold uppercase">
                          Used
                        </div>
                      </div>
                      <button
                        onClick={() => handleEditQuota(provider)}
                        className="text-center p-2 bg-white border-2 border-black rounded-base shadow-neobrutal-sm hover:shadow-neobrutal hover:bg-gray-50 transition-all cursor-pointer group"
                        title="Click to edit quota limit"
                      >
                        <div className="text-lg font-black text-neobrutal-green flex items-center justify-center gap-1">
                          {(provider.quota_limit || 0).toLocaleString()}
                          <Edit2 className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" />
                        </div>
                        <div className="text-[10px] text-gray-500 font-bold uppercase">
                          Limit
                        </div>
                      </button>
                      <div className="text-center p-2 bg-white border-2 border-black rounded-base shadow-neobrutal-sm">
                        <div
                          className={`text-lg font-black ${
                            quotaPercentage >= 80
                              ? "text-red-500"
                              : quotaPercentage >= 50
                                ? "text-amber-500"
                                : "text-neobrutal-green"
                          }`}
                        >
                          {quotaPercentage.toFixed(1)}%
                        </div>
                        <div className="text-[10px] text-gray-500 font-bold uppercase">
                          Usage
                        </div>
                      </div>
                    </div>
                    {provider.quota_limit > 0 && (
                      <Progress
                        value={quotaPercentage}
                        max={100}
                        variant={progressVariant}
                      />
                    )}
                    {provider.response_time_ms &&
                      provider.response_time_ms > 0 && (
                        <div className="mt-3 text-sm text-gray-600 flex items-center gap-2">
                          <span className="font-medium">Response Time:</span>
                          <span className="font-mono bg-white px-2 py-0.5 border border-black/20 rounded">
                            {provider.response_time_ms}ms
                          </span>
                        </div>
                      )}
                  </div>
                )}
              </CardContent>
            </Card>
          );
        })}
      </div>

      {connectedAccounts.length === 0 && autoDetected.length === 0 && (
        <Card>
          <CardContent className="pt-6 text-center text-gray-500">
            <p>No providers connected. Add a provider to get started.</p>
          </CardContent>
        </Card>
      )}

      <Card className="mt-6">
        <CardHeader className="mb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Plus className="w-5 h-5" />
            Add Provider
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {[
              { label: "OpenAI (Codex)", provider: "openai" as ProviderType },
              { label: "Claude Code", provider: "claude" as ProviderType },
              { label: "Gemini CLI", provider: "gemini" as ProviderType },
              { label: "Antigravity", provider: "antigravity" as ProviderType },
              { label: "Cursor", provider: "cursor" as ProviderType },
              { label: "Qwen Code", provider: "qwen" as ProviderType },
              { label: "iFlow", provider: "iflow" as ProviderType },
              { label: "Vertex AI", provider: "vertex" as ProviderType },
              {
                label: "Kiro (CodeWhisperer)",
                provider: "kiro" as ProviderType,
              },
              { label: "GitHub Copilot", provider: "copilot" as ProviderType },
              { label: "Ampcode", provider: "ampcode" as ProviderType },
              { label: "Z.AI", provider: "z.ai" as ProviderType },
            ].map((tile) => {
              const count = providers.filter(
                (p) => p.provider === tile.provider
              ).length;

              return (
                <button
                  key={tile.label}
                  onClick={() => {
                    setFormData((prev) => ({
                      ...prev,
                      provider: tile.provider,
                    }));
                    setShowAddModal(true);
                  }}
                  className="flex items-center justify-between p-3 bg-white border-2 border-black rounded-base shadow-neobrutal-sm hover:shadow-neobrutal hover:translate-x-[1px] hover:translate-y-[1px] transition-all group"
                >
                  <div className="flex items-center gap-2">
                    <span className="text-lg group-hover:scale-110 transition-transform">
                      {getProviderIcon(tile.provider)}
                    </span>
                    <span className="font-bold text-xs truncate">
                      {tile.label}
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 min-w-max">
                    {count > 0 && (
                      <Badge variant="success" className="h-4 px-1">
                        {count}
                      </Badge>
                    )}
                    <Plus className="w-3.5 h-3.5" />
                  </div>
                </button>
              );
            })}
          </div>
        </CardContent>
      </Card>

      <Modal
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
        title="Add Provider"
      >
        <form
          onSubmit={(e: React.FormEvent) => {
            e.preventDefault();
            addProviderMutation.mutate(formData, {
              onSuccess: () => {
                showSuccess("Provider added successfully!");
                setShowAddModal(false);
              },
              onError: () => {
                showError("Failed to save provider");
              },
            });
          }}
          className="space-y-4"
        >
          <div>
            <label className="block text-sm font-bold mb-2">Provider</label>
            <Select
              value={formData.provider}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  provider: e.target.value as ProviderType,
                })
              }
              required
            >
              <option value="openai">OpenAI (inc. Codex)</option>
              <option value="claude">Claude Code</option>
              <option value="gemini">Gemini CLI</option>
              <option value="antigravity">
                Antigravity (Claude via Gemini)
              </option>
              <option value="cursor">Cursor</option>
              <option value="copilot">GitHub Copilot</option>
              <option value="qwen">Qwen Code</option>
              <option value="vertex">Vertex AI</option>
              <option value="iflow">iFlow</option>
              <option value="kiro">Kiro (CodeWhisperer)</option>
              <option value="ampcode">Ampcode</option>
              <option value="z.ai">Z.AI</option>
            </Select>
          </div>
          <div>
            <label className="block text-sm font-bold mb-2">Name</label>
            <Input
              type="text"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              placeholder="Account name"
            />
          </div>
          <div>
            <label className="block text-sm font-bold mb-2">API Key</label>
            <Input
              type="password"
              value={formData.api_key}
              onChange={(e) =>
                setFormData({ ...formData, api_key: e.target.value })
              }
              placeholder="Enter API key"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-bold mb-2">
              Quota Limit (0 for unlimited)
            </label>
            <Input
              type="number"
              value={formData.quota_limit}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  quota_limit: parseInt(e.target.value) || 0,
                })
              }
              min="0"
            />
          </div>
          <div className="flex gap-3 justify-end pt-4">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setShowAddModal(false)}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary">
              Add Provider
            </Button>
          </div>
        </form>
      </Modal>

      <Modal
        isOpen={showDetectModal}
        onClose={() => setShowDetectModal(false)}
        title="Detected Providers"
      >
        {detectedProviders.length > 0 ? (
          <div className="space-y-3 mb-4">
            {detectedProviders.map((detected, idx) => (
              <div
                key={idx}
                className="flex items-center justify-between p-3 bg-gray-50 border-2 border-black"
              >
                <div className="flex items-center gap-3">
                  <span className="text-2xl">
                    {getProviderIcon(detected.provider)}
                  </span>
                  <div>
                    <div className="font-semibold">{detected.provider}</div>
                    <div className="text-sm text-gray-600">{detected.name}</div>
                  </div>
                </div>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => handleAddDetected(detected)}
                >
                  Add
                </Button>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-center text-gray-500 py-4">
            No providers detected from environment variables
          </p>
        )}
        <div className="flex justify-end">
          <Button variant="secondary" onClick={() => setShowDetectModal(false)}>
            Close
          </Button>
        </div>
      </Modal>

      <AlertDialog
        isOpen={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        onConfirm={confirmDelete}
        title="Delete Provider"
        message="Are you sure you want to delete this provider? This action cannot be undone."
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
      />

      {/* Edit Quota Modal */}
      <Modal
        isOpen={showEditModal}
        onClose={() => {
          setShowEditModal(false);
          setEditingProvider(null);
        }}
        title="Edit Quota Limit"
      >
        <div className="space-y-4">
          <p className="text-sm text-gray-600">
            Set a custom quota limit for{" "}
            <strong>{editingProvider?.name || "this provider"}</strong>. This is
            useful for providers that don't expose rate limits via headers (like
            Gemini).
          </p>
          <div>
            <label className="block text-sm font-bold mb-2">
              Quota Limit (tokens)
            </label>
            <Input
              type="number"
              value={editQuotaLimit}
              onChange={(e) => setEditQuotaLimit(parseInt(e.target.value) || 0)}
              min="0"
              placeholder="0 for unlimited"
            />
            <p className="text-xs text-gray-500 mt-1">
              Set to 0 for unlimited quota
            </p>
          </div>
          <div className="flex gap-3 justify-end pt-4">
            <Button
              variant="secondary"
              onClick={() => {
                setShowEditModal(false);
                setEditingProvider(null);
              }}
            >
              Cancel
            </Button>
            <Button
              variant="primary"
              onClick={handleSaveQuota}
              disabled={updateProviderMutation.isPending}
            >
              {updateProviderMutation.isPending ? "Saving..." : "Save Quota"}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default Providers;
