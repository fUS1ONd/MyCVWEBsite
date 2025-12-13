import { useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { toast } from '@/hooks/use-toast';

interface CommentFormProps {
  postSlug: string;
  value: string;
  onChange: (value: string) => void;
  onSuccess?: () => void;
}

export function CommentForm({ postSlug, value, onChange, onSuccess }: CommentFormProps) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (data: { content: string }) => {
      await axiosInstance.post(`/api/v1/posts/${postSlug}/comments`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postSlug] });
      onChange('');
      toast({ title: 'Comment posted' });
      onSuccess?.();
    },
    onError: () => {
      toast({ title: 'Failed to post comment', variant: 'destructive' });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!value.trim()) return;

    mutation.mutate({ content: value.trim() });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <Textarea
        id="main-comment-textarea"
        placeholder="Write a comment..."
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="min-h-[80px]"
      />
      <div className="flex gap-2">
        <Button type="submit" disabled={!value.trim() || mutation.isPending} size="sm">
          {mutation.isPending ? 'Posting...' : 'Comment'}
        </Button>
      </div>
    </form>
  );
}
