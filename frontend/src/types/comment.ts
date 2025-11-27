import type { User } from './user';

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  user?: User;
  content: string;
  parent_id?: number;
  replies?: Comment[];
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface CreateCommentRequest {
  content: string;
  parent_id?: number;
}

export interface UpdateCommentRequest {
  content: string;
}

export interface CommentsListResponse {
  comments: Comment[];
  total: number;
}
