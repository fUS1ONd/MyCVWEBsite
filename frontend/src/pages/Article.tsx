import { useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Post, Comment } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { CommentSection } from '@/components/comments/CommentSection';
import { format } from 'date-fns';
import { ChevronLeft, Heart, Share2, Clock } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/AuthContext';
import { SmartImage } from '@/components/ui/smart-image';

export default function Article() {
  const { toast } = useToast();
  const { slug } = useParams<{ slug: string }>();
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const {
    data: post,
    isLoading: postLoading,
    error: postError,
  } = useQuery({
    queryKey: ['post', slug],
    queryFn: async () => {
      const response = await axiosInstance.get<Post>(`/api/v1/posts/${slug}`);
      return response.data;
    },
    enabled: !!slug,
  });

  const likeMutation = useMutation({
    mutationFn: async () => {
      if (!post) return;
      const response = await axiosInstance.post<{ is_liked: boolean; likes_count: number }>(
        `/api/v1/posts/${post.id}/like`
      );
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['post', slug] });
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      toast({
        title: data?.is_liked ? 'Post liked' : 'Post unliked',
      });
    },
    onError: () => {
      toast({
        title: 'Failed to like post',
        variant: 'destructive',
      });
    },
  });

  const handleLike = () => {
    if (!user) {
      toast({
        title: 'Please log in to like posts',
        variant: 'destructive',
      });
      return;
    }
    likeMutation.mutate();
  };

  const handleShare = async () => {
    const url = window.location.href;
    try {
      await navigator.clipboard.writeText(url);
      toast({
        title: 'Link copied to clipboard',
        description: url,
      });
    } catch (error) {
      toast({
        title: 'Failed to copy link',
        variant: 'destructive',
      });
    }
  };

  const {
    data: comments,
    isLoading: commentsLoading,
    error: commentsError,
  } = useQuery({
    queryKey: ['comments', slug],
    queryFn: async () => {
      const response = await axiosInstance.get<Comment[]>(`/api/v1/posts/${slug}/comments`);
      return response.data;
    },
    enabled: !!slug,
  });

  useEffect(() => {
    if (postError) {
      toast({
        title: 'Ошибка загрузки статьи',
        description:
          postError instanceof Error ? postError.message : 'Не удалось подключиться к серверу',
        variant: 'destructive',
      });
    }
  }, [postError, toast]);

  useEffect(() => {
    if (commentsError) {
      toast({
        title: 'Ошибка загрузки комментариев',
        description:
          commentsError instanceof Error
            ? commentsError.message
            : 'Не удалось загрузить комментарии',
        variant: 'destructive',
      });
    }
  }, [commentsError, toast]);

  if (postLoading) {
    return (
      <div className="container max-w-4xl py-16 space-y-8">
        <Skeleton className="h-10 w-3/4" />
        <Skeleton className="h-4 w-32" />
        <div className="space-y-4">
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-3/4" />
        </div>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="container max-w-4xl py-16 text-center">
        <p className="text-muted-foreground">Post not found</p>
        <Button asChild variant="outline" className="mt-4">
          <Link to="/blog">
            <ChevronLeft className="mr-2 h-4 w-4" />
            Back to Blog
          </Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="container max-w-4xl py-16 space-y-8">
      <Button asChild variant="ghost" size="sm">
        <Link to="/blog">
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Blog
        </Link>
      </Button>

      <article className="space-y-6">
        {/* Cover Image */}
        {post.cover_image && (
          <SmartImage
            src={post.cover_image}
            alt={post.title}
            containerClassName="rounded-lg max-h-[60vh]"
            className="rounded-lg"
            loading="lazy"
          />
        )}

        <header className="space-y-4">
          <h1 className="text-2xl sm:text-4xl font-bold tracking-tight">{post.title}</h1>

          {/* Author Info */}
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10">
              <AvatarImage src={post.author?.avatar_url} alt={post.author?.name} />
              <AvatarFallback>{post.author?.name?.charAt(0).toUpperCase()}</AvatarFallback>
            </Avatar>
            <div className="flex-1 flex items-center gap-2 text-sm text-muted-foreground">
              <span className="font-medium">{post.author?.name}</span>
              <span>•</span>
              <time dateTime={post.created_at}>
                {format(new Date(post.created_at), 'MMMM dd, yyyy HH:mm')}
              </time>
              {post.read_time_minutes > 0 && (
                <>
                  <span>•</span>
                  <div className="flex items-center gap-1">
                    <Clock className="h-3 w-3" />
                    <span>{post.read_time_minutes} min read</span>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex items-center gap-4 pt-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleLike}
              disabled={likeMutation.isPending || !user}
              className="gap-2"
            >
              <Heart className={`h-4 w-4 ${post.is_liked ? 'fill-red-500 text-red-500' : ''}`} />
              <span className={post.is_liked ? 'text-red-500' : ''}>{post.likes_count}</span>
            </Button>

            <Button variant="outline" size="sm" onClick={handleShare} className="gap-2">
              <Share2 className="h-4 w-4" />
              Share
            </Button>
          </div>
        </header>

        <Card>
          <CardContent className="pt-6">
            <div
              className="prose prose-sm sm:prose-base prose-slate dark:prose-invert max-w-none"
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          </CardContent>
        </Card>
      </article>

      <div id="comments" className="pt-8 border-t">
        <CommentSection postSlug={slug!} comments={comments || []} isLoading={commentsLoading} />
      </div>
    </div>
  );
}
