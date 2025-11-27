import type { ProfileInfo, UpdateProfileRequest } from '../types';
import { apiClient } from './api';

export const profileService = {
  // Get profile info (public)
  getProfile: async (): Promise<ProfileInfo> => {
    return apiClient.get<ProfileInfo>('/profile');
  },

  // Admin: Update profile
  updateProfile: async (data: UpdateProfileRequest): Promise<ProfileInfo> => {
    return apiClient.put<ProfileInfo>('/admin/profile', data);
  },
};
