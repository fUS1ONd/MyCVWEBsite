import { useState, useEffect } from 'react';
import { Comment } from '@/lib/types';
import { CommentThread } from './CommentThread';
import { CommentForm } from './CommentForm';
import { CommentSort, SortOption } from './CommentSort';
import { useAuth } from '@/contexts/AuthContext';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { MessageSquare } from 'lucide-react';

interface CommentSectionProps {
  postSlug: string;
  comments: Comment[];
  isLoading: boolean;
}

export function CommentSection({ postSlug, comments, isLoading }: CommentSectionProps) {
  const { isAuthenticated } = useAuth();
  const [replyingTo, setReplyingTo] = useState<number | null>(null);
  const [sortOption, setSortOption] = useState<SortOption>(() => {
    const saved = localStorage.getItem('commentSort');
    return (saved as SortOption) || 'newest';
  });

  useEffect(() => {
    localStorage.setItem('commentSort', sortOption);
  }, [sortOption]);

  // Organize comments into threads
  const topLevelComments = comments.filter((c) => !c.parent_id);

  // Sort comments
  const sortedComments = [...topLevelComments].sort((a, b) => {
    switch (sortOption) {
      case 'newest':
        return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
      case 'oldest':
        return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
      case 'most_liked':
        return (b.likes_count || 0) - (a.likes_count || 0);
      default:
        return 0;
    }
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent className="space-y-4">
          {[...Array(2)].map((_, i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <MessageSquare className="h-5 w-5" />
            Comments ({comments.length})
          </CardTitle>
          {topLevelComments.length > 0 && (
            <CommentSort value={sortOption} onChange={setSortOption} />
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {!isAuthenticated && (
          <Alert>
            <AlertDescription>Please sign in to leave a comment</AlertDescription>
          </Alert>
        )}

        {isAuthenticated && (
          <CommentForm postSlug={postSlug} onSuccess={() => setReplyingTo(null)} />
        )}

        <div className="space-y-6">
          {sortedComments.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              No comments yet. Be the first to comment!
            </p>
          ) : (
            sortedComments.map((comment) => (
              <CommentThread
                key={comment.id}
                comment={comment}
                postSlug={postSlug}
                allComments={comments}
                replyingTo={replyingTo}
                setReplyingTo={setReplyingTo}
              />
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
}
