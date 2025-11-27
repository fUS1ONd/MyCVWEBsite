export const UserRole = {
  User: 'user',
  Admin: 'admin',
} as const;

export type UserRole = (typeof UserRole)[keyof typeof UserRole];

export interface User {
  id: number;
  email: string;
  name: string;
  avatar_url?: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}

export interface Session {
  user: User;
  token: string;
  expires_at: string;
}

export interface LoginResponse {
  user: User;
  token: string;
}
