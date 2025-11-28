export interface User {
  id: number;
  email: string;
  name: string;
  avatar_url?: string;
  role: 'user' | 'admin';
  created_at: string;
}

export interface Profile {
  name: string;
  title?: string;
  description?: string;
  bio?: string;
  avatar_url?: string;
  github_url?: string;
  linkedin_url?: string;
  twitter_url?: string;
  website_url?: string;
  email?: string;
}

export interface Post {
  id: number;
  title: string;
  slug: string;
  content: string;
  preview: string;
  published: boolean;
  created_at: string;
  updated_at: string;
  author_id: number;
  author?: User;
}

export interface Comment {
  id: number;
  content: string;
  post_id: number;
  user_id: number;
  parent_id?: number;
  created_at: string;
  updated_at: string;
  user?: User;
  replies?: Comment[];
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}
