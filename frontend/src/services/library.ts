import { apiClient } from './api'

export interface LibraryItem {
  id: string
  title: string
  author: string
  narrator?: string
  series?: string
  seriesPosition?: number
  status: 'available' | 'downloading' | 'processing' | 'error'
  filePath?: string
  coverArt?: string
  duration?: number
  createdAt: string
  updatedAt: string
}

export const libraryService = {
  getAll: async (): Promise<LibraryItem[]> => {
    return apiClient.get<LibraryItem[]>('/api/v1/library')
  },

  getById: async (id: string): Promise<LibraryItem> => {
    return apiClient.get<LibraryItem>(`/api/v1/library/${id}`)
  },

  add: async (data: {
    title: string
    author: string
    isbn?: string
  }): Promise<LibraryItem> => {
    return apiClient.post<LibraryItem>('/api/v1/library', data)
  },

  remove: async (id: string): Promise<void> => {
    return apiClient.delete(`/api/v1/library/${id}`)
  },
}

