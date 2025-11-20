import axios, { AxiosInstance, AxiosError } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8686'

export interface APIResponse<T> {
  success: boolean
  data?: T
  error?: string
}

class APIClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add API key if available
        const apiKey = localStorage.getItem('apiKey')
        if (apiKey) {
          config.headers['X-API-Key'] = apiKey
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Handle unauthorized - API key invalid or missing
          localStorage.removeItem('apiKey')
          // Could redirect to settings page to enter API key
          console.error('API key invalid or missing')
        }
        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, params?: unknown): Promise<T> {
    const response = await this.client.get<APIResponse<T>>(url, { params })
    return response.data.data as T
  }

  async post<T>(url: string, data?: unknown): Promise<T> {
    const response = await this.client.post<APIResponse<T>>(url, data)
    return response.data.data as T
  }

  async put<T>(url: string, data?: unknown): Promise<T> {
    const response = await this.client.put<APIResponse<T>>(url, data)
    return response.data.data as T
  }

  async delete<T>(url: string): Promise<T> {
    const response = await this.client.delete<APIResponse<T>>(url)
    return response.data.data as T
  }
}

export const apiClient = new APIClient()

