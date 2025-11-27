import type { User, LoginResponse } from '../types';
import { apiClient } from './api';

export const authService = {
  // Get current user info
  getCurrentUser: async (): Promise<User> => {
    return apiClient.get<User>('/auth/me');
  },

  // Logout
  logout: async (): Promise<void> => {
    await apiClient.post('/auth/logout');
    localStorage.removeItem('auth_token');
  },

  // OAuth login URL generation
  getOAuthUrl: (provider: 'vk' | 'google' | 'github'): string => {
    const baseUrl =
      import.meta.env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:8080';
    return `${baseUrl}/auth/${provider}`;
  },

  // Handle OAuth callback (this would typically be handled by backend redirect)
  handleOAuthCallback: async (token: string): Promise<LoginResponse> => {
    localStorage.setItem('auth_token', token);
    const user = await authService.getCurrentUser();
    return { user, token };
  },
};
