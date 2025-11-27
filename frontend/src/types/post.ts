import type { User } from './user';

export interface Post {
  id: number;
  title: string;
  slug: string;
  content: string;
  preview: string;
  author_id: number;
  author?: User;
  published: boolean;
  published_at?: string;
  created_at: string;
  updated_at: string;
  media?: MediaFile[];
}

export interface CreatePostRequest {
  title: string;
  content: string;
  preview: string;
  published: boolean;
}

export interface UpdatePostRequest {
  title?: string;
  content?: string;
  preview?: string;
  published?: boolean;
}

export interface PostsListResponse {
  posts: Post[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface MediaFile {
  id: number;
  filename: string;
  path: string;
  mime_type: string;
  size: number;
  uploader_id: number;
  created_at: string;
}
