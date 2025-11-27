import type {
  Comment,
  CommentsListResponse,
  CreateCommentRequest,
  UpdateCommentRequest,
} from '../types';
import { apiClient } from './api';

export const commentsService = {
  // Get comments for a post
  getCommentsByPostSlug: async (slug: string): Promise<CommentsListResponse> => {
    return apiClient.get<CommentsListResponse>(`/posts/${slug}/comments`);
  },

  // Create comment
  createComment: async (slug: string, data: CreateCommentRequest): Promise<Comment> => {
    return apiClient.post<Comment>(`/posts/${slug}/comments`, data);
  },

  // Update comment
  updateComment: async (id: number, data: UpdateCommentRequest): Promise<Comment> => {
    return apiClient.put<Comment>(`/comments/${id}`, data);
  },

  // Delete comment
  deleteComment: async (id: number): Promise<void> => {
    return apiClient.delete(`/comments/${id}`);
  },
};
