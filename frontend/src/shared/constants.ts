/**
 * Shared constants
 */

export const APP_NAME = 'Quotio' as const;
export const APP_VERSION = '0.1.0' as const;

export const API_BASE_URL = 'http://localhost:8080/api' as const;

export const ROUTES = {
  DASHBOARD: '/',
  QUOTA: '/quota',
  PROVIDERS: '/providers',
  AGENTS: '/agents',
  SETTINGS: '/settings',
  ABOUT: '/about',
} as const;

export type Route = typeof ROUTES[keyof typeof ROUTES];

