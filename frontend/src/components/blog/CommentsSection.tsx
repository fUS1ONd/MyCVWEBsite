import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { commentsService } from '../../services/comments.service';
import { Comment } from './Comment';
import { CommentForm } from './CommentForm';
import { toast } from '../../utils/toast';

interface CommentsSectionProps {
  postSlug: string;
}

export const CommentsSection: React.FC<CommentsSectionProps> = ({ postSlug }) => {
  const [replyingTo, setReplyingTo] = useState<number | null>(null);
  const queryClient = useQueryClient();

  // Fetch comments
  const { data, isLoading, error } = useQuery({
    queryKey: ['comments', postSlug],
    queryFn: () => commentsService.getCommentsByPostSlug(postSlug),
  });

  // Create comment mutation
  const createCommentMutation = useMutation({
    mutationFn: ({ content, parentId }: { content: string; parentId?: number }) =>
      commentsService.createComment(postSlug, { content, parent_id: parentId }),
    onSuccess: (newComment) => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      toast.success('Комментарий добавлен');
      setReplyingTo(null);

      // Scroll to the new comment
      setTimeout(() => {
        const element = document.getElementById(`comment-${newComment.id}`);
        element?.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }, 100);
    },
    onError: () => {
      toast.error('Не удалось добавить комментарий');
    },
  });

  // Update comment mutation
  const updateCommentMutation = useMutation({
    mutationFn: ({ id, content }: { id: number; content: string }) =>
      commentsService.updateComment(id, { content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      toast.success('Комментарий обновлен');
    },
    onError: () => {
      toast.error('Не удалось обновить комментарий');
    },
  });

  // Delete comment mutation
  const deleteCommentMutation = useMutation({
    mutationFn: (id: number) => commentsService.deleteComment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      toast.success('Комментарий удален');
    },
    onError: () => {
      toast.error('Не удалось удалить комментарий');
    },
  });

  const handleCreateComment = async (content: string, parentId?: number) => {
    await createCommentMutation.mutateAsync({ content, parentId });
  };

  const handleEditComment = (commentId: number, content: string) => {
    updateCommentMutation.mutate({ id: commentId, content });
  };

  const handleDeleteComment = (commentId: number) => {
    deleteCommentMutation.mutate(commentId);
  };

  const handleReply = (parentId: number) => {
    setReplyingTo(replyingTo === parentId ? null : parentId);
  };

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
        Не удалось загрузить комментарии
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h3 className="text-2xl font-bold text-gray-900">Комментарии {data && `(${data.total})`}</h3>

      {/* New comment form */}
      <div className="bg-gray-50 rounded-lg p-4">
        <CommentForm onSubmit={handleCreateComment} />
      </div>

      {/* Loading state */}
      {isLoading && (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, index) => (
            <div key={index} className="animate-pulse">
              <div className="flex gap-3">
                <div className="w-10 h-10 bg-gray-200 rounded-full"></div>
                <div className="flex-1 space-y-2">
                  <div className="h-4 bg-gray-200 rounded w-32"></div>
                  <div className="h-4 bg-gray-200 rounded"></div>
                  <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Comments list */}
      {!isLoading && data && (
        <>
          {data.comments.length === 0 ? (
            <p className="text-gray-600 text-center py-8">Пока нет комментариев. Будьте первым!</p>
          ) : (
            <div className="divide-y divide-gray-200">
              {data.comments.map((comment) => (
                <div key={comment.id} id={`comment-${comment.id}`}>
                  <Comment
                    comment={comment}
                    onReply={handleReply}
                    onEdit={handleEditComment}
                    onDelete={handleDeleteComment}
                  />

                  {/* Reply form */}
                  {replyingTo === comment.id && (
                    <div className="ml-14 mb-4 bg-gray-50 rounded-lg p-4">
                      <CommentForm
                        onSubmit={handleCreateComment}
                        parentId={comment.id}
                        onCancel={() => setReplyingTo(null)}
                        placeholder={`Ответить ${comment.user?.name || 'пользователю'}...`}
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
};
