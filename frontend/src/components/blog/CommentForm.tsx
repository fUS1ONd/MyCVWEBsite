import React, { useState } from 'react';
import { Button } from '../common/Button';

interface CommentFormProps {
  onSubmit: (content: string, parentId?: number) => Promise<void>;
  parentId?: number;
  onCancel?: () => void;
  placeholder?: string;
}

export const CommentForm: React.FC<CommentFormProps> = ({
  onSubmit,
  parentId,
  onCancel,
  placeholder = 'Напишите комментарий...',
}) => {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const maxLength = 1000;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!content.trim() || isSubmitting) {
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit(content.trim(), parentId);
      setContent('');
    } catch (error) {
      console.error('Failed to submit comment:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder={placeholder}
        className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
        rows={4}
        maxLength={maxLength}
        disabled={isSubmitting}
      />
      <div className="flex items-center justify-between">
        <span className="text-sm text-gray-500">
          {content.length} / {maxLength}
        </span>
        <div className="flex gap-2">
          {onCancel && (
            <Button type="button" variant="secondary" onClick={onCancel} disabled={isSubmitting}>
              Отмена
            </Button>
          )}
          <Button type="submit" disabled={!content.trim() || isSubmitting}>
            {isSubmitting ? 'Отправка...' : parentId ? 'Ответить' : 'Отправить'}
          </Button>
        </div>
      </div>
    </form>
  );
};
