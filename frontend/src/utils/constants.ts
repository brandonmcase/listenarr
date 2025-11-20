export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8686'

export const ROUTES = {
  DASHBOARD: '/',
  LIBRARY: '/library',
  DOWNLOADS: '/downloads',
  PROCESSING: '/processing',
  SEARCH: '/search',
  SETTINGS: '/settings',
} as const

