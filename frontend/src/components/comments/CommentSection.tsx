import { useState, useEffect } from 'react';
import { Comment } from '@/lib/types';
import { CommentItem } from './CommentItem';
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
  const [commentContent, setCommentContent] = useState('');
  const [sortOption, setSortOption] = useState<SortOption>(() => {
    const saved = localStorage.getItem('commentSort');
    return (saved as SortOption) || 'newest';
  });

  useEffect(() => {
    localStorage.setItem('commentSort', sortOption);
  }, [sortOption]);

  // Sort comments
  const sortedComments = [...comments].sort((a, b) => {
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

  const handleReply = (authorName: string) => {
    const mention = `@${authorName}, `;
    setCommentContent((prev) => (prev.startsWith(mention) ? prev : mention + prev));

    // Focus the textarea
    const textarea = document.getElementById('main-comment-textarea');
    if (textarea) {
      textarea.focus();
      textarea.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  };

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
          {comments.length > 0 && <CommentSort value={sortOption} onChange={setSortOption} />}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {!isAuthenticated && (
          <Alert>
            <AlertDescription>Please sign in to leave a comment</AlertDescription>
          </Alert>
        )}

        {isAuthenticated && (
          <CommentForm postSlug={postSlug} value={commentContent} onChange={setCommentContent} />
        )}

        <div className="space-y-6">
          {sortedComments.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              No comments yet. Be the first to comment!
            </p>
          ) : (
            sortedComments.map((comment) => (
              <CommentItem
                key={comment.id}
                comment={comment}
                postSlug={postSlug}
                onReply={handleReply}
              />
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
}
