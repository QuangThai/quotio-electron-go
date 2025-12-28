import { Badge } from '../ui/badge';

/**
 * Render account/provider status badge
 * Handles: active | rate_limited | cooldown | disabled
 */
export function renderAccountStatusBadge(status?: string) {
  if (!status || status === 'active') {
    return <Badge variant="success">Active</Badge>;
  }
  if (status === 'rate_limited') {
    return <Badge variant="warning">Rate Limited</Badge>;
  }
  if (status === 'cooldown') {
    return <Badge variant="danger">Cooldown</Badge>;
  }
  return <Badge variant="secondary">Disabled</Badge>;
}

/**
 * Render agent status badge
 * Handles: configured | installed | not_installed | config_error
 */
export function renderAgentStatusBadge(
  installed: boolean,
  configured: boolean,
  hasError: boolean
) {
  if (hasError) {
    return <Badge variant="danger">Config Error</Badge>;
  }
  if (configured) {
    return <Badge variant="success">Configured</Badge>;
  }
  if (installed) {
    return <Badge variant="warning">Installed</Badge>;
  }
  return <Badge variant="secondary">Not Installed</Badge>;
}
