export const NotificationType = {
  NewPost: 'new_post',
  NewComment: 'new_comment',
  System: 'system',
} as const;

export type NotificationType = (typeof NotificationType)[keyof typeof NotificationType];

export interface Notification {
  id: number;
  user_id: number;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  created_at: string;
}

export interface NotificationSettings {
  user_id: number;
  email_enabled: boolean;
  push_enabled: boolean;
  new_posts_enabled: boolean;
}

export interface UpdateNotificationSettingsRequest {
  email_enabled?: boolean;
  push_enabled?: boolean;
  new_posts_enabled?: boolean;
}
