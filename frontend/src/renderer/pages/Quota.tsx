import {
  AlertCircle,
  AlertTriangle,
  Circle,
  Clock,
  RefreshCcw,
  RotateCcw,
  TrendingUp,
  XCircle,
  Zap,
} from "lucide-react";
import { useEffect, useState } from "react";
import type { QuotaInfo, ToastContextValue } from "../../types";
import { renderAccountStatusBadge } from "../components/shared/status-badges";
import { AlertDialog } from "../components/ui/alert-dialog";
import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { Progress } from "../components/ui/progress";
import { useToast } from "../components/ui/toast-container";
import { getProviderIcon, getProviderLabel } from "../lib/provider-metadata";
import { useQuota, useQuotaHistory, useResetQuota } from "../queries";
import api from "../services/api";

interface FailedRequest {
  id: number;
  account_id: number;
  provider: string;
  account_name: string;
  model?: string;
  status_code: number;
  tokens_used: number;
  timestamp: string;
}

interface QuotaInfoExtended extends QuotaInfo {
  auto_detected_limit?: number;
  is_manual_quota?: boolean;
  rate_limit_requests?: number;
  rate_limit_requests_reset?: string;
  rate_limit_tokens?: number;
  rate_limit_tokens_reset?: string;
}

function Quota() {
  const { data: accountsData = [], isLoading } = useQuota();
  const accounts = accountsData as QuotaInfoExtended[];
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [selectedAccount, setSelectedAccount] = useState<string | null>(null);
  const [failedRequests, setFailedRequests] = useState<FailedRequest[]>([]);
  const [showManualInput, setShowManualInput] = useState<boolean>(false);
  const { showSuccess, showError } = useToast() as ToastContextValue;

  // Fetch quota history for selected account
  const { data: history = [] } = useQuotaHistory(selectedAccount);
  const resetQuotaMutation = useResetQuota();

  // PHASE C: Fetch failed requests on mount
  useEffect(() => {
    const fetchFailed = async () => {
      try {
        const response = await api.get("/quota/failed");
        setFailedRequests(response.data || []);
      } catch {
        // Handle fetch error silently
      }
    };
    fetchFailed();
  }, []);

  // Set first provider as default when data loads
  useEffect(() => {
    if (accounts.length > 0) {
      const firstProvider = accounts[0].provider;
      if (
        !selectedProvider ||
        !accounts.some((a) => a.provider === selectedProvider)
      ) {
        setSelectedProvider(firstProvider);
      }
    }
  }, [accounts]);

  // Set first account of selected provider
  useEffect(() => {
    if (selectedProvider) {
      const providerAccounts = accounts.filter(
        (a) => a.provider === selectedProvider
      );
      if (
        providerAccounts.length > 0 &&
        (!selectedAccount ||
          !providerAccounts.some((a) => String(a.id) === selectedAccount))
      ) {
        setSelectedAccount(String(providerAccounts[0].id));
      }
    }
  }, [selectedProvider, accounts]);

  const [pendingReset, setPendingReset] = useState<{
    id: string;
    name: string;
  } | null>(null);

  const handleResetQuota = (accountId: string, accountName: string) => {
    setPendingReset({ id: accountId, name: accountName });
  };

  const confirmReset = () => {
    if (!pendingReset) return;

    resetQuotaMutation.mutate(pendingReset.id, {
      onSuccess: () => {
        showSuccess("Quota reset successfully!");
        setPendingReset(null);
      },
      onError: () => {
        showError("Failed to reset quota");
        setPendingReset(null);
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

  // Group accounts by provider
  const groupedByProvider: Record<string, QuotaInfo[]> = accounts.reduce(
    (acc, account) => {
      if (!acc[account.provider]) {
        acc[account.provider] = [];
      }
      acc[account.provider].push(account);
      return acc;
    },
    {} as Record<string, QuotaInfo[]>
  );

  const providers = Object.keys(groupedByProvider);
  const activeProvider = selectedProvider || providers[0] || null;
  const activeAccounts = activeProvider
    ? groupedByProvider[activeProvider]
    : [];

  // Calculate total stats for selected provider
  const totalQuotaUsed = activeAccounts.reduce(
    (sum, acc) => sum + (acc.quota_used || 0),
    0
  );
  const totalQuotaLimit = activeAccounts.reduce(
    (sum, acc) => sum + (acc.quota_limit > 0 ? acc.quota_limit : 0),
    0
  );
  const totalPercentage =
    totalQuotaLimit > 0 ? (totalQuotaUsed / totalQuotaLimit) * 100 : 0;

  return (
    <div className="p-6 max-w-6xl mx-auto animate-fade-in">
      <div className="mb-6">
        <h2 className="text-2xl font-black mb-1">Quota</h2>
        <p className="text-sm text-gray-500 font-medium">
          Track quota usage and account status
        </p>
      </div>

      {providers.length > 0 && (
        <div className="flex gap-1.5 mb-6 flex-wrap">
          {providers.map((provider) => {
            const providerAccounts = groupedByProvider[provider];
            const hasActiveAccount = providerAccounts.some((acc) => {
              const pct =
                acc.quota_limit > 0
                  ? (acc.quota_used / acc.quota_limit) * 100
                  : 0;
              return acc.status !== "disabled" && pct < 100;
            });

            return (
              <button
                key={provider}
                onClick={() => setSelectedProvider(provider)}
                className={`px-3 py-1.5 font-bold border-2 border-black rounded-base transition-all ${
                  activeProvider === provider
                    ? "bg-neobrutal-blue text-white shadow-neobrutal-sm"
                    : "bg-white text-black shadow-neobrutal-sm hover:shadow-neobrutal hover:translate-x-[1px] hover:translate-y-[1px]"
                }`}
              >
                <div className="flex items-center gap-2">
                  <span>{getProviderIcon(provider)}</span>
                  <span className="text-xs">{getProviderLabel(provider)}</span>
                  <Badge
                    variant={
                      activeProvider === provider ? "default" : "secondary"
                    }
                    className="h-4 px-1 text-[10px]"
                  >
                    {providerAccounts.length}
                  </Badge>
                  <Circle
                    className={`w-2 h-2 ${hasActiveAccount ? "fill-green-500 text-green-500" : "fill-red-500 text-red-500"}`}
                  />
                </div>
              </button>
            );
          })}
        </div>
      )}

      {activeAccounts.length > 0 ? (
        <div className="space-y-6">
          {/* Summary Card */}
          <Card className="mb-6">
            <CardContent className="pt-5">
              <h3 className="text-lg font-black mb-4 uppercase tracking-wider">
                Total Usage
              </h3>
              <div className="grid grid-cols-3 gap-3">
                <div className="p-3 bg-white border-2 border-black rounded-base shadow-neobrutal-sm">
                  <div className="text-[10px] text-gray-500 font-bold uppercase mb-1">
                    Tokens Used
                  </div>
                  <div className="text-xl font-black">
                    {totalQuotaUsed.toLocaleString()}
                  </div>
                </div>
                <div className="p-3 bg-white border-2 border-black rounded-base shadow-neobrutal-sm">
                  <div className="text-[10px] text-gray-500 font-bold uppercase mb-1">
                    Total Limit
                  </div>
                  <div className="text-xl font-black">
                    {totalQuotaLimit.toLocaleString()}
                  </div>
                </div>
                <div className="p-3 bg-white border-2 border-black rounded-base shadow-neobrutal-sm">
                  <div className="text-[10px] text-gray-500 font-bold uppercase mb-1">
                    Usage
                  </div>
                  <div className="text-xl font-black">
                    {totalPercentage.toFixed(1)}%
                  </div>
                </div>
              </div>
              {totalQuotaLimit > 0 && (
                <Progress
                  value={totalPercentage}
                  max={100}
                  variant={
                    totalPercentage >= 80
                      ? "danger"
                      : totalPercentage >= 50
                        ? "warning"
                        : "success"
                  }
                  className="mt-4"
                />
              )}
            </CardContent>
          </Card>

          {/* Account Details */}
          {activeAccounts.map((account) => {
            const quotaPercentage =
              account.quota_limit > 0
                ? (account.quota_used / account.quota_limit) * 100
                : 0;
            const progressVariant =
              quotaPercentage >= 80
                ? "danger"
                : quotaPercentage >= 50
                  ? "warning"
                  : "success";
            const accountHistory = history.filter(
              (h) => h.account_id === Number(account.id)
            );

            return (
              <Card key={account.id}>
                <CardContent className="pt-6">
                  {/* Account Header */}
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-2xl">
                          {getProviderIcon(account.provider)}
                        </span>
                        <div>
                          <div className="text-base font-black">
                            {account.name || account.provider}
                          </div>
                          <div className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
                            {account.provider} • ID: {account.id}
                          </div>
                        </div>
                      </div>
                      {/* Status Badges */}
                      <div className="flex gap-1.5 mb-4">
                        {renderAccountStatusBadge(account.status)}
                        {account.auto_detected && (
                          <Badge variant="secondary">Auto-detected</Badge>
                        )}
                        {account.is_healthy !== undefined && (
                          <Badge
                            variant={account.is_healthy ? "success" : "danger"}
                          >
                            {account.is_healthy ? "Healthy" : "Unhealthy"}
                          </Badge>
                        )}
                        {/* NEW: Manual/Auto toggle */}
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setShowManualInput(!showManualInput)}
                          className="ml-auto"
                          title={
                            account.is_manual_quota
                              ? "Switch to auto-detect"
                              : "Set manual limit"
                          }
                        >
                          <RefreshCcw className="w-3.5 h-3.5" />
                        </Button>
                      </div>
                    </div>
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() =>
                        handleResetQuota(String(account.id), account.name)
                      }
                      disabled={resetQuotaMutation.isPending}
                      title="Reset quota"
                      className="h-8 w-8 p-0"
                    >
                      <RotateCcw className="w-3.5 h-3.5" />
                    </Button>
                  </div>

                  {/* Auto-Detected Limit */}
                  {account.auto_detected_limit &&
                    !account.is_manual_quota &&
                    (account.auto_detected_limit > 0 ||
                      account.auto_detected_limit === -1) && (
                      <div className="bg-blue-50 border-l-4 border-blue-500 rounded p-3 mb-4">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-blue-600 font-bold text-xs uppercase tracking-wider">
                            Auto-Detected Limit
                          </span>
                          <Badge variant="primary">From Provider</Badge>
                        </div>
                        <div className="text-2xl font-black">
                          {account.auto_detected_limit === -1
                            ? "Unknown"
                            : account.auto_detected_limit.toLocaleString()}
                        </div>
                        <div className="text-xs text-blue-700">
                          {account.auto_detected_limit === -1
                            ? "Provider does not report accurate rate limits via API"
                            : "Limit based on provider headers (anthropic-ratelimit, x-ratelimit)"}
                        </div>
                      </div>
                    )}

                  {/* Manual Override */}
                  {account.is_manual_quota && account.quota_limit > 0 && (
                    <div className="bg-purple-50 border-l-4 border-purple-500 rounded p-3 mb-4">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-purple-600 font-bold text-xs uppercase tracking-wider">
                          Manual Override
                        </span>
                        <Badge variant="warning">User Set</Badge>
                      </div>
                      <div className="text-2xl font-black">
                        {account.quota_limit.toLocaleString()}
                      </div>
                      <div className="text-xs text-purple-700">
                        Override auto-detected limit
                      </div>
                    </div>
                  )}

                  {/* Quota Progress */}
                  {(account.quota_limit > 0 || account.quota_limit === -1) && (
                    <div className="mb-4">
                      {account.quota_limit === -1 ? (
                        <div className="p-3 bg-gray-50 border border-gray-200 rounded text-center">
                          <div className="text-xs font-bold uppercase tracking-wider text-gray-500 mb-1">
                            Quota Usage
                          </div>
                          <div className="text-lg font-black text-gray-700">
                            {account.quota_used.toLocaleString()} tokens
                          </div>
                          <div className="text-[10px] text-gray-400 font-bold uppercase">
                            Limit: Unknown / Unlimited
                          </div>
                        </div>
                      ) : (
                        <>
                          <div className="flex justify-between mb-1.5">
                            <span className="text-xs font-bold uppercase tracking-wider text-gray-500">
                              Quota Usage
                            </span>
                            <span
                              className={`text-xs font-black ${progressVariant === "danger" ? "text-red-500" : progressVariant === "warning" ? "text-yellow-600" : "text-green-500"}`}
                            >
                              {quotaPercentage.toFixed(1)}%
                            </span>
                          </div>
                          <Progress
                            value={quotaPercentage}
                            max={100}
                            variant={progressVariant}
                            className="h-3"
                          />
                          <div className="text-[10px] text-gray-400 font-bold uppercase mt-1.5">
                            {account.quota_used.toLocaleString()} /{" "}
                            {account.quota_limit.toLocaleString()} tokens
                          </div>
                        </>
                      )}
                    </div>
                  )}

                  {/* Model Usage Breakdown */}
                  {account.model_usage &&
                    Object.keys(account.model_usage).length > 0 && (
                      <div className="mt-6">
                        <h4 className="font-black text-xs uppercase tracking-widest mb-4 flex items-center gap-2">
                          <TrendingUp className="w-4 h-4" />
                          Model Usage
                        </h4>
                        <div className="space-y-4">
                          {Object.entries(account.model_usage || {}).map(
                            ([model, tokens]) => {
                              const totalUsed = Object.values(
                                account.model_usage || {}
                              ).reduce((sum, val) => sum + (val as number), 0);
                              const modelPercentage =
                                totalUsed > 0 ? (tokens / totalUsed) * 100 : 0;
                              return (
                                <div key={model}>
                                  <div className="flex justify-between mb-1.5 text-[10px] font-bold uppercase tracking-wider">
                                    <span className="font-mono text-gray-500">
                                      {model}
                                    </span>
                                    <span className="text-black">
                                      {tokens.toLocaleString()} tokens (
                                      {modelPercentage.toFixed(1)}%)
                                    </span>
                                  </div>
                                  <Progress
                                    value={modelPercentage}
                                    max={100}
                                    variant="success"
                                    className="h-2"
                                  />
                                </div>
                              );
                            }
                          )}
                        </div>
                      </div>
                    )}

                  {/* Recent History */}
                  {accountHistory.length > 0 && (
                    <div className="mt-6">
                      <h4 className="font-black text-xs uppercase tracking-widest mb-4 flex items-center gap-2">
                        <Clock className="w-4 h-4" />
                        Recent Requests
                      </h4>
                      <div className="space-y-2">
                        {accountHistory.slice(0, 5).map((entry) => (
                          <div
                            key={entry.id}
                            className="flex items-center justify-between p-2.5 bg-surface border-2 border-black rounded-base text-xs font-bold"
                          >
                            <div className="flex items-center gap-2">
                              {entry.success ? (
                                <Circle className="w-2 h-2 fill-green-500 text-green-500" />
                              ) : (
                                <AlertTriangle className="w-2 h-2 text-red-500" />
                              )}
                              {entry.model && (
                                <span className="font-mono">{entry.model}</span>
                              )}
                              <span className="text-gray-600">
                                {entry.tokens_used} tokens
                              </span>
                            </div>
                            <span className="text-xs text-gray-500">
                              {new Date(entry.timestamp).toLocaleString()}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Response Time */}
                  {account.response_time_ms && (
                    <div className="mt-4 flex items-center gap-2 text-sm text-gray-600">
                      <Zap className="w-4 h-4" />
                      <span>Response time: {account.response_time_ms}ms</span>
                      <span className="text-gray-400">
                        • Last checked: {account.last_checked}
                      </span>
                    </div>
                  )}
                </CardContent>
              </Card>
            );
          })}
        </div>
      ) : (
        <Card>
          <CardContent className="pt-6 text-center text-gray-500">
            <p>
              No accounts configured. Add providers to start tracking quota.
            </p>
          </CardContent>
        </Card>
      )}

      {/* PHASE C: Show disabled accounts separately */}
      {accounts.some((a) => a.status === "disabled") && (
        <div className="mt-8">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2 text-red-600">
            <AlertCircle className="w-5 h-5" />
            Disabled Accounts (Invalid Credentials)
          </h2>
          <div className="space-y-4">
            {accounts
              .filter((a) => a.status === "disabled")
              .map((account) => (
                <Card key={account.id} className="border-red-400">
                  <CardContent className="pt-6">
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-3">
                        <XCircle className="w-5 h-5 text-red-500" />
                        <div>
                          <div className="font-bold">
                            {account.name || account.provider}
                          </div>
                          <div className="text-sm text-gray-600">
                            {account.provider} • ID: {account.id}
                          </div>
                        </div>
                      </div>
                      <Badge variant="danger">Invalid</Badge>
                    </div>
                    <p className="text-sm text-gray-600 mt-3">
                      This account's credentials failed validation. Check your
                      API key or refresh the credentials.
                    </p>
                  </CardContent>
                </Card>
              ))}
          </div>
        </div>
      )}

      {/* PHASE C: Show failed requests section */}
      {failedRequests.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2 text-red-600">
            <AlertTriangle className="w-5 h-5" />
            Failed Requests ({failedRequests.length})
          </h2>
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-gray-600 mb-4">
                These requests failed, likely due to invalid credentials
                (401/403) or provider errors (429, 5xx). Click the Validate
                button on the Providers page to re-check credentials.
              </p>
              <div className="space-y-2">
                {failedRequests.slice(0, 20).map((req) => {
                  const getBadgeVariant = (statusCode: number) => {
                    if (statusCode === 401 || statusCode === 403)
                      return "danger";
                    if (statusCode === 429) return "warning";
                    if (statusCode >= 500) return "danger";
                    return "secondary";
                  };

                  return (
                    <div
                      key={req.id}
                      className="flex items-center justify-between p-3 bg-red-50 border border-red-200 rounded text-sm"
                    >
                      <div className="flex items-center gap-3 flex-1">
                        <span className="text-lg">
                          {getProviderIcon(req.provider)}
                        </span>
                        <div className="flex-1">
                          <div className="font-semibold">
                            {req.account_name || req.provider}
                          </div>
                          <div className="text-xs text-gray-600">
                            {req.model && (
                              <span className="font-mono">{req.model}</span>
                            )}
                            {req.model && req.tokens_used && " • "}
                            {req.tokens_used && (
                              <span>{req.tokens_used} tokens</span>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant={getBadgeVariant(req.status_code)}>
                          {req.status_code === 401 || req.status_code === 403
                            ? "Auth Failed"
                            : req.status_code === 429
                              ? "Rate Limited"
                              : req.status_code >= 500
                                ? "Server Error"
                                : "Other Error"}
                        </Badge>
                        <span className="text-xs text-gray-600">
                          {new Date(req.timestamp).toLocaleString()}
                        </span>
                      </div>
                    </div>
                  );
                })}
                {failedRequests.length > 20 && (
                  <div className="text-xs text-gray-500 p-3 bg-gray-50 text-center">
                    ... and {failedRequests.length - 20} more failed requests
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <AlertDialog
        isOpen={!!pendingReset}
        onClose={() => setPendingReset(null)}
        onConfirm={confirmReset}
        title="Reset Quota Usage"
        message={`Are you sure you want to reset the quota usage for ${pendingReset?.name}? This action cannot be undone.`}
        variant="danger"
        confirmText="Reset Quota"
        cancelText="Cancel"
      />
    </div>
  );
}

export default Quota;
