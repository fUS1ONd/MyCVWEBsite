import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { toast } from '@/hooks/use-toast';

interface CommentFormProps {
  postSlug: string;
  parentId?: number;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function CommentForm({ postSlug, parentId, onSuccess, onCancel }: CommentFormProps) {
  const [content, setContent] = useState('');
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (data: { content: string; parent_id?: number }) => {
      await axiosInstance.post(`/api/v1/posts/${postSlug}/comments`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      setContent('');
      toast({ title: 'Comment posted' });
      onSuccess?.();
    },
    onError: () => {
      toast({ title: 'Failed to post comment', variant: 'destructive' });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;

    mutation.mutate({
      content: content.trim(),
      parent_id: parentId,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <Textarea
        placeholder={parentId ? 'Write a reply...' : 'Write a comment...'}
        value={content}
        onChange={(e) => setContent(e.target.value)}
        className="min-h-[80px]"
      />
      <div className="flex gap-2">
        <Button type="submit" disabled={!content.trim() || mutation.isPending} size="sm">
          {mutation.isPending ? 'Posting...' : parentId ? 'Reply' : 'Comment'}
        </Button>
        {onCancel && (
          <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </form>
  );
}
