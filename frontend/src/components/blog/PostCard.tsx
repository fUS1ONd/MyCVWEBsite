import { Link } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Post } from '@/lib/types';
import { Card } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { toast } from '@/hooks/use-toast';
import { Heart, MessageSquare, Share2, Clock } from 'lucide-react';
import { format } from 'date-fns';
import { useAuth } from '@/contexts/AuthContext';
import { SmartImage } from '@/components/ui/smart-image';

interface PostCardProps {
  post: Post;
}

export function PostCard({ post }: PostCardProps) {
  const queryClient = useQueryClient();
  const { user } = useAuth();

  const likeMutation = useMutation({
    mutationFn: async () => {
      const response = await axiosInstance.post<{ is_liked: boolean; likes_count: number }>(
        `/api/v1/posts/${post.id}/like`
      );
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.slug] });
      toast({
        title: data.is_liked ? 'Post liked' : 'Post unliked',
      });
    },
    onError: () => {
      toast({
        title: 'Failed to like post',
        variant: 'destructive',
      });
    },
  });

  const handleShare = async () => {
    const url = `${window.location.origin}/blog/${post.slug}`;
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

  const stripHtml = (html: string) => {
    const tmp = document.createElement('div');
    tmp.innerHTML = html;
    return tmp.textContent || tmp.innerText || '';
  };

  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow">
      {/* Cover Image */}
      {post.cover_image && (
        <Link to={`/blog/${post.slug}`} className="block">
          <SmartImage src={post.cover_image} alt={post.title} loading="lazy" />
        </Link>
      )}

      <div className="p-6 space-y-4">
        {/* Title */}
        <Link to={`/blog/${post.slug}`}>
          <h2 className="text-xl sm:text-2xl font-bold hover:text-primary transition-colors line-clamp-2">
            {post.title}
          </h2>
        </Link>

        {/* Preview Text */}
        <p className="text-muted-foreground line-clamp-3">{stripHtml(post.preview)}</p>

        {/* Author Info */}
        <div className="flex items-center gap-3">
          <Avatar className="h-8 w-8">
            <AvatarImage src={post.author?.avatar_url} alt={post.author?.name} />
            <AvatarFallback>{post.author?.name?.charAt(0).toUpperCase()}</AvatarFallback>
          </Avatar>
          <div className="flex-1 flex items-center gap-2 text-sm text-muted-foreground">
            <span className="font-medium">{post.author?.name}</span>
            <span>•</span>
            <time dateTime={post.created_at}>
              {format(new Date(post.created_at), 'MMM dd, yyyy HH:mm')}
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

        {/* Action Bar */}
        <div className="flex flex-wrap items-center gap-2 sm:gap-4 pt-4 border-t">
          <Button
            variant="ghost"
            size="sm"
            onClick={handleLike}
            disabled={likeMutation.isPending || !user}
            className="gap-2"
          >
            <Heart className={`h-4 w-4 ${post.is_liked ? 'fill-red-500 text-red-500' : ''}`} />
            <span className={post.is_liked ? 'text-red-500' : ''}>{post.likes_count}</span>
          </Button>

          <Link to={`/blog/${post.slug}#comments`}>
            <Button variant="ghost" size="sm" className="gap-2">
              <MessageSquare className="h-4 w-4" />
              <span>{post.comments_count}</span>
            </Button>
          </Link>

          <Button variant="ghost" size="sm" onClick={handleShare} className="gap-2 ml-auto">
            <Share2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </Card>
  );
}
