import { useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Post, Comment } from '@/lib/types';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { CommentSection } from '@/components/comments/CommentSection';
import { format } from 'date-fns';
import { ChevronLeft } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useToast } from '@/hooks/use-toast';

export default function Article() {
  const { toast } = useToast();
  const { slug } = useParams<{ slug: string }>();

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
        <header className="space-y-4">
          <h1 className="text-4xl font-bold tracking-tight">{post.title}</h1>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <time>{format(new Date(post.created_at), 'MMMM dd, yyyy')}</time>
            {post.author && <span>by {post.author.name}</span>}
          </div>
        </header>

        <Card>
          <CardContent className="pt-6">
            <div className="prose prose-slate dark:prose-invert max-w-none">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>{post.content}</ReactMarkdown>
            </div>
          </CardContent>
        </Card>
      </article>

      <div className="pt-8 border-t">
        <CommentSection postSlug={slug!} comments={comments || []} isLoading={commentsLoading} />
      </div>
    </div>
  );
}
