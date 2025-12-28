import {
  useQuery,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import type {
  Provider,
  AgentStatus,
  Settings,
} from '../../types'
import {
  getDashboard,
  getProviders,
  detectProviders,
  addProvider as apiAddProvider,
  updateProvider as apiUpdateProvider,
  deleteProvider as apiDeleteProvider,
  getQuota,
  getAgents,
  configureAgent,
  startProxy as apiStartProxy,
  stopProxy as apiStopProxy,
  getProxyStatus as apiGetProxyStatus,
  getSettings,
  updateSettings as apiUpdateSettings,
  getProviderHealth,
  checkProviderHealth,
  getQuotaHistory,
  resetQuota,
  getModels,
  updateRoutingStrategy,
  refreshAgents as apiRefreshAgents,
} from '../services/api'

// Dashboard Queries
export const useDashboard = () =>
  useQuery({
    queryKey: ['dashboard'],
    queryFn: () => getDashboard(),
    refetchInterval: 5000, // Real-time updates
  })

export const useProxyStatus = () =>
  useQuery({
    queryKey: ['proxy-status'],
    queryFn: () => apiGetProxyStatus(),
    refetchInterval: 5000,
  })

// Provider Queries
export const useProviders = () =>
  useQuery({
    queryKey: ['providers'],
    queryFn: () => getProviders(),
  })

export const useAddProvider = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (provider: Partial<Provider>) => apiAddProvider(provider),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
    },
  })
}

export const useUpdateProvider = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, provider }: { id: string; provider: Partial<Provider> }) =>
      apiUpdateProvider(id, provider),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
    },
  })
}

export const useDeleteProvider = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => apiDeleteProvider(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
    },
  })
}

export const useDetectProviders = () =>
  useMutation({
    mutationFn: () => detectProviders(),
  })

// Quota Queries
export const useQuota = () =>
  useQuery({
    queryKey: ['quota'],
    queryFn: () => getQuota(),
    refetchInterval: 5000, // Real-time updates
  })

// Agent Queries
export const useAgents = () =>
  useQuery({
    queryKey: ['agents'],
    queryFn: () => getAgents(),
  })

export const useConfigureAgent = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (config: Partial<AgentStatus>) => configureAgent(config),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['agents'] })
    },
  })
}

// Settings Queries
export const useSettings = () =>
  useQuery({
    queryKey: ['settings'],
    queryFn: () => getSettings(),
  })

export const useUpdateSettings = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (settings: Settings) => apiUpdateSettings(settings),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings'] })
    },
  })
}

// Proxy Mutations
export const useStartProxy = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => apiStartProxy(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['proxy-status'] })
    },
  })
}

export const useStopProxy = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => apiStopProxy(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['proxy-status'] })
    },
  })
}

// Provider Health Queries
export const useProviderHealth = () =>
  useQuery({
    queryKey: ['provider-health'],
    queryFn: () => getProviderHealth(),
    refetchInterval: 10000, // 10 seconds
  })

export const useCheckHealth = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => checkProviderHealth(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
      queryClient.invalidateQueries({ queryKey: ['provider-health'] })
    },
  })
}

// Quota History & Model Queries
export const useQuotaHistory = (accountId: string | null, limit?: number) =>
  useQuery({
    queryKey: ['quota-history', accountId, limit],
    queryFn: () => getQuotaHistory(accountId!, limit),
    enabled: !!accountId,
  })

export const useResetQuota = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => resetQuota(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['quota'] })
      queryClient.invalidateQueries({ queryKey: ['providers'] })
    },
  })
}

export const useModels = (provider?: string) =>
  useQuery({
    queryKey: ['models', provider],
    queryFn: () => getModels(provider),
  })

// Routing Strategy Queries
export const useRoutingStrategy = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (strategy: 'round_robin' | 'fill_first') => updateRoutingStrategy(strategy),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings'] })
    },
  })
}

// Agent Refresh Queries
export const useRefreshAgents = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => apiRefreshAgents(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['agents'] })
    },
  })
}

