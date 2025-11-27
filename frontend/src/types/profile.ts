export interface ProfileInfo {
  id: number;
  name: string;
  description: string;
  photo_url: string;
  activity: string;
  contacts: ProfileContacts;
  created_at: string;
  updated_at: string;
}

export interface ProfileContacts {
  email?: string;
  github?: string;
  linkedin?: string;
  vk?: string;
}

export interface UpdateProfileRequest {
  name?: string;
  description?: string;
  photo_url?: string;
  activity?: string;
  contacts?: ProfileContacts;
}
