// Provider Types
export type ProviderType =
  | 'openai'
  | 'claude'
  | 'gemini'
  | 'copilot'
  | 'antigravity'
  | 'codex'
  | 'cursor'
  | 'kiro'
  | 'vertex'
  | 'qwen'
  | 'iflow'
  | 'ampcode'
  | 'z.ai';

export interface Provider {
  id: number;
  provider: ProviderType;
  name: string;
  api_key?: string;
  quota_limit: number;
  quota_used: number;
  quota_manual?: boolean;
  status?: 'active' | 'rate_limited' | 'cooldown' | 'disabled';
  auto_detected?: boolean;
  supports_manual_auth?: boolean;
  model_access?: string[];
  priority?: number;
  is_healthy?: boolean;
  response_time_ms?: number;
  last_checked?: string;
  refresh_token?: string;
  token_expires_at?: string;
}

export interface ProviderHealth {
  id: number;
  account_id: number;
  provider_name: string;
  account_name: string;
  is_healthy: boolean;
  response_time_ms: number;
  last_checked: string;
  consecutive_failures: number;
}

export interface ModelQuota {
  model: string;
  tokens_used: number;
  requests: number;
  last_used: string;
}

export interface QuotaHistory {
  id: number;
  account_id: number;
  model?: string;
  tokens_used: number;
  requests_count: number;
  status_code: number;
  success: boolean;
  timestamp: string;
}

export interface ProviderStats {
  provider: string;
  accounts: number;
}

// Dashboard Types
export interface Dashboard {
  active_accounts: number;
  requests_today: number;
  tokens_today: number;
  success_rate: number;
  providers: ProviderStats[];
}

// Proxy Types
export interface ProxyStatus {
  running: boolean;
  port: number;
}

// Agent Types
export interface Agent {
  id: string;
  name: string;
  provider: ProviderType;
  config: Record<string, unknown>;
}

export interface AgentStatus {
  id?: string;
  agent_name: string;
  installed?: boolean;
  configured?: boolean;
  config_path?: string;
  config_path_exists?: boolean;
  validation_error?: string;
  proxy_url?: string;
  auto_configured?: boolean;
  last_configured?: string;
}

// Settings Types
export interface Settings {
  port?: number;
  routing_strategy?: string;
  auto_start?: boolean;
  api_key?: string;
  proxy_port?: number;
  [key: string]: unknown;
}

// Quota Types
export interface QuotaInfo {
  id: number;
  provider: string;
  name: string;
  quota_limit: number;
  quota_used: number;
  model_usage?: Record<string, number>; // Per-model token usage
  is_healthy?: boolean;
  response_time_ms?: number;
  last_checked?: string;
  status?: 'active' | 'rate_limited' | 'cooldown' | 'disabled';
  auto_detected?: boolean;
  // NEW: Auto-detected limits
  auto_detected_limit?: number;
  is_manual_quota?: boolean;
  rate_limit_requests?: number;
  rate_limit_requests_reset?: string;
  rate_limit_tokens?: number;
  rate_limit_tokens_reset?: string;
}

// API Response Types
export interface ApiResponse<T = unknown> {
  data: T;
  status: number;
  statusText: string;
}

export interface ApiError {
  message: string;
  code?: number;
}

// Toast Types
export type ToastVariant = 'info' | 'success' | 'error' | 'warning';

export interface Toast {
  id: string;
  message: string;
  variant: ToastVariant;
  duration: number;
}

export interface ToastContextValue {
  showToast: (message: string, variant?: ToastVariant, duration?: number) => string;
  showSuccess: (message: string, duration?: number) => string;
  showError: (message: string, duration?: number) => string;
  showWarning: (message: string, duration?: number) => string;
  showInfo: (message: string, duration?: number) => string;
  removeToast: (id: string) => void;
}

// Electron Types
export interface ElectronAPI {
  platform: string;
  versions: {
    node?: string;
    chrome?: string;
    electron?: string;
  };
}

// UI Component Types
export type ButtonVariant = 'default' | 'primary' | 'secondary' | 'danger' | 'success' | 'warning' | 'ghost' | 'purple';
export type ButtonSize = 'sm' | 'default' | 'lg';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  children: React.ReactNode;
}

export type ProgressVariant = 'default' | 'success' | 'warning' | 'danger';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  className?: string;
}

export interface AlertDialogProps extends React.HTMLAttributes<HTMLDivElement> {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  variant?: ButtonVariant;
  confirmText?: string;
  cancelText?: string;
}

export interface TabProps {
  id: string;
  label: string;
  children?: React.ReactNode;
}

export interface TabsProps {
  tabs: TabProps[];
  activeTab: string;
  onTabChange: (tabId: string) => void;
}

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: ButtonVariant;
  children: React.ReactNode;
}

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
}

export interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export interface Alert {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message?: string;
  dismissible?: boolean;
  duration?: number;
}

// Navigation Types
export interface MenuItem {
  path: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}

