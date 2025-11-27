import React, { useState } from 'react';
import type { Comment as CommentType } from '../../types';
import { Button } from '../common/Button';
import { useAuth } from '../../contexts/AuthContext';

interface CommentProps {
  comment: CommentType;
  onReply?: (parentId: number) => void;
  onEdit?: (commentId: number, content: string) => void;
  onDelete?: (commentId: number) => void;
}

export const Comment: React.FC<CommentProps> = ({ comment, onReply, onEdit, onDelete }) => {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(comment.content);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const isOwner = user?.id === comment.user_id;
  const isAdmin = user?.role === 'admin';
  const canEdit = isOwner;
  const canDelete = isOwner || isAdmin;

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'только что';
    if (diffMins < 60) return `${diffMins} мин. назад`;
    if (diffHours < 24) return `${diffHours} ч. назад`;
    if (diffDays < 7) return `${diffDays} дн. назад`;

    return date.toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const handleSaveEdit = () => {
    if (editContent.trim() && onEdit) {
      onEdit(comment.id, editContent.trim());
      setIsEditing(false);
    }
  };

  const handleCancelEdit = () => {
    setEditContent(comment.content);
    setIsEditing(false);
  };

  const handleDelete = () => {
    if (onDelete) {
      onDelete(comment.id);
      setShowDeleteConfirm(false);
    }
  };

  if (comment.deleted_at) {
    return <div className="py-4 text-gray-500 italic">Комментарий был удален</div>;
  }

  return (
    <div className="py-4">
      <div className="flex gap-3">
        {/* Avatar */}
        <div className="flex-shrink-0">
          {comment.user?.avatar_url ? (
            <img
              src={comment.user.avatar_url}
              alt={comment.user.name}
              className="w-10 h-10 rounded-full"
            />
          ) : (
            <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 font-semibold">
              {comment.user?.name.charAt(0).toUpperCase() || '?'}
            </div>
          )}
        </div>

        {/* Comment content */}
        <div className="flex-1 min-w-0">
          {/* Header */}
          <div className="flex items-center gap-2 mb-1">
            <span className="font-semibold text-gray-900">{comment.user?.name || 'Аноним'}</span>
            <span className="text-sm text-gray-500">{formatDate(comment.created_at)}</span>
            {comment.updated_at !== comment.created_at && (
              <span className="text-xs text-gray-400">(изменено)</span>
            )}
          </div>

          {/* Content */}
          {isEditing ? (
            <div className="space-y-2">
              <textarea
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                rows={3}
                maxLength={1000}
              />
              <div className="flex gap-2">
                <Button size="sm" onClick={handleSaveEdit}>
                  Сохранить
                </Button>
                <Button size="sm" variant="secondary" onClick={handleCancelEdit}>
                  Отмена
                </Button>
              </div>
            </div>
          ) : (
            <>
              <p className="text-gray-700 whitespace-pre-wrap break-words">{comment.content}</p>

              {/* Actions */}
              <div className="flex items-center gap-4 mt-2">
                {onReply && (
                  <button
                    onClick={() => onReply(comment.id)}
                    className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                  >
                    Ответить
                  </button>
                )}
                {canEdit && (
                  <button
                    onClick={() => setIsEditing(true)}
                    className="text-sm text-gray-600 hover:text-gray-700 font-medium"
                  >
                    Редактировать
                  </button>
                )}
                {canDelete && (
                  <>
                    {!showDeleteConfirm ? (
                      <button
                        onClick={() => setShowDeleteConfirm(true)}
                        className="text-sm text-red-600 hover:text-red-700 font-medium"
                      >
                        Удалить
                      </button>
                    ) : (
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-gray-600">Уверены?</span>
                        <button
                          onClick={handleDelete}
                          className="text-sm text-red-600 hover:text-red-700 font-medium"
                        >
                          Да
                        </button>
                        <button
                          onClick={() => setShowDeleteConfirm(false)}
                          className="text-sm text-gray-600 hover:text-gray-700 font-medium"
                        >
                          Нет
                        </button>
                      </div>
                    )}
                  </>
                )}
              </div>
            </>
          )}

          {/* Nested replies */}
          {comment.replies && comment.replies.length > 0 && (
            <div className="mt-4 ml-4 border-l-2 border-gray-200 pl-4 space-y-4">
              {comment.replies.map((reply) => (
                <Comment
                  key={reply.id}
                  comment={reply}
                  onReply={onReply}
                  onEdit={onEdit}
                  onDelete={onDelete}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
