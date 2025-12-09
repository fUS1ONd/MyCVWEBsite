import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Comment } from '@/lib/types';
import { useAuth } from '@/contexts/AuthContext';
import axiosInstance from '@/lib/axios';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { CommentForm } from './CommentForm';
import { format } from 'date-fns';
import { MessageSquare, Trash2, Heart } from 'lucide-react';
import { toast } from '@/hooks/use-toast';

interface CommentItemProps {
  comment: Comment;
  postSlug: string;
  replyingTo: number | null;
  setReplyingTo: (id: number | null) => void;
  depth: number;
}

export function CommentItem({
  comment,
  postSlug,
  replyingTo,
  setReplyingTo,
  depth,
}: CommentItemProps) {
  const { user, isAuthenticated } = useAuth();
  const queryClient = useQueryClient();
  const [isDeleting, setIsDeleting] = useState(false);

  const deleteMutation = useMutation({
    mutationFn: async () => {
      await axiosInstance.delete(`/api/v1/comments/${comment.id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      toast({ title: 'Comment deleted' });
    },
    onError: () => {
      toast({ title: 'Failed to delete comment', variant: 'destructive' });
    },
  });

  const likeMutation = useMutation({
    mutationFn: async () => {
      const response = await axiosInstance.post<{ is_liked: boolean; likes_count: number }>(
        `/api/v1/comments/${comment.id}/like`
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      queryClient.invalidateQueries({ queryKey: ['post', postSlug] });
    },
    onError: () => {
      toast({ title: 'Failed to like comment', variant: 'destructive' });
    },
  });

  const handleDelete = () => {
    if (confirm('Are you sure you want to delete this comment?')) {
      deleteMutation.mutate();
    }
  };

  const handleLike = () => {
    if (!isAuthenticated) {
      toast({ title: 'Please log in to like comments', variant: 'destructive' });
      return;
    }
    likeMutation.mutate();
  };

  const isAuthor = user?.id === comment.user_id;
  const canReply = isAuthenticated && depth < 10; // Limit nesting depth

  return (
    <div className="space-y-3">
      <div className="flex gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors">
        <Avatar className={depth === 0 ? 'h-10 w-10' : 'h-8 w-8'}>
          <AvatarImage src={comment.user?.avatar_url} />
          <AvatarFallback>{comment.user?.name.charAt(0).toUpperCase() || 'U'}</AvatarFallback>
        </Avatar>

        <div className="flex-1 space-y-2">
          <div className="flex items-center gap-2">
            <span className="font-medium text-sm">{comment.user?.name || 'User'}</span>
            <span className="text-xs text-muted-foreground">
              {format(new Date(comment.created_at), 'MMM dd, yyyy HH:mm')}
            </span>
          </div>

          <p className="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
            {comment.content}
          </p>

          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              className="h-7 px-2 text-xs"
              onClick={handleLike}
              disabled={likeMutation.isPending || !isAuthenticated}
            >
              <Heart
                className={`h-3 w-3 mr-1 ${comment.is_liked ? 'fill-red-500 text-red-500' : ''}`}
              />
              <span className={comment.is_liked ? 'text-red-500' : ''}>
                {comment.likes_count > 0 ? comment.likes_count : ''}
              </span>
            </Button>

            {canReply && (
              <Button
                variant="ghost"
                size="sm"
                className="h-7 px-2 text-xs"
                onClick={() => setReplyingTo(replyingTo === comment.id ? null : comment.id)}
              >
                <MessageSquare className="h-3 w-3 mr-1" />
                Reply
              </Button>
            )}

            {isAuthor && (
              <Button
                variant="ghost"
                size="sm"
                className="h-7 px-2 text-xs text-destructive hover:text-destructive"
                onClick={handleDelete}
                disabled={deleteMutation.isPending}
              >
                <Trash2 className="h-3 w-3 mr-1" />
                Delete
              </Button>
            )}
          </div>
        </div>
      </div>

      {replyingTo === comment.id && (
        <div className="ml-11">
          <CommentForm
            postSlug={postSlug}
            parentId={comment.id}
            onSuccess={() => setReplyingTo(null)}
            onCancel={() => setReplyingTo(null)}
          />
        </div>
      )}
    </div>
  );
}
