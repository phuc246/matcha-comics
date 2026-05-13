import { useRuntimeConfig } from '#app'

export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = 'http://localhost:8080/api' // In production, this would come from config.public.apiBase

  const get = async <T>(endpoint: string, options = {}) => {
    return await $fetch<T>(`${baseURL}${endpoint}`, {
      method: 'GET',
      ...options
    })
  }

  const post = async <T>(endpoint: string, body: any, options = {}) => {
    return await $fetch<T>(`${baseURL}${endpoint}`, {
      method: 'POST',
      body,
      ...options
    })
  }

  const put = async <T>(endpoint: string, body: any, options = {}) => {
    return await $fetch<T>(`${baseURL}${endpoint}`, {
      method: 'PUT',
      body,
      ...options
    })
  }

  const del = async <T>(endpoint: string, options = {}) => {
    return await $fetch<T>(`${baseURL}${endpoint}`, {
      method: 'DELETE',
      ...options
    })
  }

  return {
    get,
    post,
    put,
    del
  }
}
