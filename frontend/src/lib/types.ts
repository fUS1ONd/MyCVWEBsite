export interface User {
  id: number;
  email: string;
  name: string;
  avatar_url?: string;
  role: 'user' | 'admin';
  created_at: string;
}

export interface Profile {
  id: number;
  name: string;
  description: string;
  photo_url?: string;
  activity: string;
  contacts: {
    email: string;
    github: string;
    linkedin: string;
  };
}

export interface MediaFile {
  id: number;
  filename: string;
  url: string;
  mime_type: string;
  size: number;
  uploader_id: number;
  created_at: string;
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
  cover_image?: string;
  read_time_minutes: number;
  likes_count: number;
  comments_count: number;
  is_liked: boolean;
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
  likes_count: number;
  is_liked: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  meta?: {
    page?: number;
    page_size?: number;
    total_pages?: number;
    total_count?: number;
  };
  error?: {
    code: string;
    message: string;
  };
}

export interface PostListResponse {
  posts: Post[];
  total_count: number;
  page: number;
  limit: number;
  total_pages: number;
}
