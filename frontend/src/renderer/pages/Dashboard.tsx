import {
  ArrowUpDown,
  CheckCircle2,
  Copy,
  Hash,
  RefreshCw,
  Users,
} from "lucide-react";
import { ToastContextValue } from "../../types";
import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { useToast } from "../components/ui/toast-container";
import { getProviderIcon } from "../lib/provider-metadata";
import { useDashboard, useProxyStatus } from "../queries";

function Dashboard() {
  const { data: dashboard, isLoading } = useDashboard();
  const { data: proxyStatus } = useProxyStatus();
  const { showSuccess } = useToast() as ToastContextValue;

  const handleCopyEndpoint = () => {
    const endpoint = `http://localhost:${proxyStatus?.port || 8317}/v1`;
    navigator.clipboard.writeText(endpoint);
    showSuccess("Endpoint copied to clipboard");
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-lg font-bold">Loading...</div>
      </div>
    );
  }

  const successRate =
    dashboard?.success_rate != null
      ? (dashboard.success_rate * 100).toFixed(0)
      : 0;
  const failedRequests =
    (dashboard?.requests_today ?? 0) > 0
      ? Math.round(
          (dashboard?.requests_today ?? 0) -
            (dashboard?.success_rate ?? 0) * (dashboard?.requests_today ?? 0)
        )
      : 0;

  return (
    <div className="p-6 max-w-6xl mx-auto animate-fade-in">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-1">
          <h2 className="text-2xl font-black">Dashboard</h2>
          <Button
            variant="secondary"
            size="sm"
            onClick={() => window.location.reload()}
          >
            <RefreshCw className="w-3.5 h-3.5 mr-1.5" />
            Refresh
          </Button>
        </div>
        <p className="text-sm text-gray-500 font-medium">
          Monitor your AI proxy server status and activity
        </p>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <Card>
          <CardContent className="text-center pt-5">
            <Users className="w-6 h-6 mx-auto mb-2 text-neobrutal-blue" />
            <div className="text-3xl font-black text-neobrutal-blue mb-0.5">
              {dashboard?.active_accounts || 0}
            </div>
            <p className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              {dashboard?.active_accounts || 0} ready
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="text-center pt-5">
            <ArrowUpDown className="w-6 h-6 mx-auto mb-2 text-neobrutal-green" />
            <div className="text-3xl font-black text-neobrutal-green mb-0.5">
              {dashboard?.requests_today || 0}
            </div>
            <p className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              total requests
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="text-center pt-5">
            <Hash className="w-6 h-6 mx-auto mb-2 text-neobrutal-purple" />
            <div className="text-3xl font-black text-neobrutal-purple mb-0.5">
              {dashboard?.tokens_today || 0}
            </div>
            <p className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              processed
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="text-center pt-5">
            <CheckCircle2 className="w-6 h-6 mx-auto mb-2 text-neobrutal-orange" />
            <div className="text-3xl font-black text-neobrutal-orange mb-0.5">
              {successRate}%
            </div>
            <p className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
              {failedRequests} failed
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="flex items-center gap-2 text-base">
              <span>📋</span>
              Providers
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-1.5">
              {dashboard?.providers && dashboard.providers.length > 0 ? (
                dashboard.providers.map((p) => (
                  <Badge key={p.provider} variant="secondary">
                    {getProviderIcon(p.provider)} {p.provider} ×{p.accounts}
                  </Badge>
                ))
              ) : (
                <p className="text-gray-400 text-xs font-medium italic">
                  No providers connected
                </p>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="mb-2">
            <CardTitle className="flex items-center gap-2 text-base">
              <span>🔗</span>
              API Endpoint
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <code className="flex-1 px-3 py-1.5 bg-surface border-2 border-black rounded-base font-mono text-xs">
                http://localhost:{proxyStatus?.port || 8317}/v1
              </code>
              <Button
                variant="secondary"
                size="sm"
                onClick={handleCopyEndpoint}
              >
                <Copy className="w-3.5 h-3.5" />
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default Dashboard;
