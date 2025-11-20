import { apiClient } from './api'

export interface Download {
  id: string
  libraryItemId: string
  title: string
  status: 'queued' | 'downloading' | 'completed' | 'failed'
  progress: number
  speed?: number
  size?: number
  downloaded?: number
  error?: string
  createdAt: string
  updatedAt: string
}

export const downloadService = {
  getAll: async (): Promise<Download[]> => {
    return apiClient.get<Download[]>('/api/v1/downloads')
  },

  getById: async (id: string): Promise<Download> => {
    return apiClient.get<Download>(`/api/v1/downloads/${id}`)
  },

  start: async (libraryItemId: string, releaseId: string): Promise<Download> => {
    return apiClient.post<Download>('/api/v1/downloads', {
      libraryItemId,
      releaseId,
    })
  },

  cancel: async (id: string): Promise<void> => {
    return apiClient.delete(`/api/v1/downloads/${id}`)
  },
}

