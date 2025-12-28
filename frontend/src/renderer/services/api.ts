import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import type {
  Dashboard,
  Provider,
  ProxyStatus,
  QuotaInfo,
  Settings,
  AgentStatus,
  ProviderHealth,
  QuotaHistory,
} from "../../types";

const API_BASE_URL = "http://localhost:8080/api";

const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // Handle API errors silently
    return Promise.reject(error);
  }
);

export const healthCheck = () => api.get("/health");

export const getDashboard = (): Promise<Dashboard> =>
  api.get("/dashboard").then(res => res.data);

export const getProviders = (): Promise<Provider[]> =>
  api.get("/providers").then(res => res.data);

export const getProviderStatus = (): Promise<Provider[]> =>
  api.get("/providers/status").then(res => res.data);

export const detectProviders = (): Promise<Provider[]> =>
  api.get("/providers/detect").then(res => res.data);

export const addProvider = (provider: Partial<Provider>): Promise<Provider> =>
  api.post("/providers", provider).then(res => res.data);

export const updateProvider = (id: string, provider: Partial<Provider>): Promise<Provider> =>
  api.put(`/providers/${id}`, provider).then(res => res.data);

export const deleteProvider = (id: string): Promise<void> =>
  api.delete(`/providers/${id}`).then(res => res.data);

export const getQuota = (): Promise<QuotaInfo[]> =>
  api.get("/quota").then(res => res.data);

export const getAgents = (): Promise<AgentStatus[]> =>
  api.get("/agents").then(res => res.data);

export const configureAgent = (config: Partial<AgentStatus>): Promise<AgentStatus> =>
  api.post("/agents/configure", config).then(res => res.data);

export const startProxy = (): Promise<{ port: number }> =>
  api.post("/proxy/start").then(res => res.data);

export const stopProxy = (): Promise<void> =>
  api.post("/proxy/stop").then(res => res.data);

export const getProxyStatus = (): Promise<ProxyStatus> =>
  api.get("/proxy/status").then(res => res.data);

export const getSettings = (): Promise<Settings> =>
  api.get("/settings").then(res => res.data);

export const updateSettings = (settings: Settings): Promise<Settings> =>
  api.put("/settings", settings).then(res => res.data);

// Provider Health API
export const getProviderHealth = (): Promise<ProviderHealth[]> =>
  api.get("/providers/health").then(res => res.data);

export const checkProviderHealth = (id: string): Promise<ProviderHealth> =>
  api.post(`/providers/health/${id}`).then(res => res.data);

// Quota History & Model API
export const getQuotaHistory = (accountId: string, limit?: number): Promise<QuotaHistory[]> =>
  api.get(`/quota/history/${accountId}${limit ? `?limit=${limit}` : ''}`).then(res => res.data);

export const resetQuota = (id: string): Promise<void> =>
  api.post(`/quota/reset/${id}`).then(res => res.data);

export const getModels = (provider?: string): Promise<Record<string, unknown>[]> =>
  api.get(`/models${provider ? `?provider=${provider}` : ''}`).then(res => res.data);

// Routing Strategy API
export const updateRoutingStrategy = (strategy: 'round_robin' | 'fill_first'): Promise<{ message: string; strategy: string }> =>
  api.post("/routing-strategy", { strategy }).then(res => res.data);

// Agents Refresh API
export const refreshAgents = (): Promise<AgentStatus[]> =>
  api.post("/agents/refresh").then(res => res.data);

export const validateAgentConfig = (agentName: string): Promise<{ valid: boolean; error?: string }> =>
  api.post("/agents/validate", { agent_name: agentName }).then(res => res.data);

export default api;

