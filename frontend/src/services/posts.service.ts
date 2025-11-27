import type {
  Post,
  PostsListResponse,
  CreatePostRequest,
  UpdatePostRequest,
  PaginationParams,
} from '../types';
import { apiClient } from './api';

export const postsService = {
  // Get list of posts (with pagination)
  getPosts: async (
    params?: PaginationParams & { published?: boolean }
  ): Promise<PostsListResponse> => {
    return apiClient.get<PostsListResponse>('/posts', { params });
  },

  // Get single post by slug
  getPostBySlug: async (slug: string): Promise<Post> => {
    return apiClient.get<Post>(`/posts/${slug}`);
  },

  // Admin: Create new post
  createPost: async (data: CreatePostRequest): Promise<Post> => {
    return apiClient.post<Post>('/admin/posts', data);
  },

  // Admin: Update post
  updatePost: async (id: number, data: UpdatePostRequest): Promise<Post> => {
    return apiClient.put<Post>(`/admin/posts/${id}`, data);
  },

  // Admin: Delete post
  deletePost: async (id: number): Promise<void> => {
    return apiClient.delete(`/admin/posts/${id}`);
  },
};
