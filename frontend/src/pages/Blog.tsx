import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import axiosInstance from '@/lib/axios';
import { Post, PaginatedResponse } from '@/lib/types';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { format } from 'date-fns';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

export default function Blog() {
  const { toast } = useToast();
  const [page, setPage] = useState(1);
  const limit = 10;

  const { data, isLoading, error } = useQuery({
    queryKey: ['posts', page, limit],
    queryFn: async () => {
      const response = await axiosInstance.get<PaginatedResponse<Post>>(
        `/api/v1/posts?page=${page}&limit=${limit}&published=true`
      );
      return response.data;
    },
  });

  useEffect(() => {
    if (error) {
      toast({
        title: 'Ошибка загрузки постов',
        description: error instanceof Error ? error.message : 'Не удалось подключиться к серверу',
        variant: 'destructive',
      });
    }
  }, [error, toast]);

  if (isLoading) {
    return (
      <div className="container max-w-4xl py-16 space-y-8">
        <div className="space-y-2">
          <Skeleton className="h-10 w-64" />
          <Skeleton className="h-4 w-96" />
        </div>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-3/4" />
                <Skeleton className="h-4 w-full" />
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  const posts = data?.data || [];
  const totalPages = data?.total_pages || 1;

  return (
    <div className="container max-w-4xl py-16 space-y-8">
      <div className="space-y-2">
        <h1 className="text-4xl font-bold tracking-tight">Blog</h1>
        <p className="text-lg text-muted-foreground">Thoughts, stories and ideas</p>
      </div>

      <div className="space-y-6">
        {posts.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">No posts published yet</p>
            </CardContent>
          </Card>
        ) : (
          posts.map((post) => (
            <Link key={post.id} to={`/blog/${post.slug}`}>
              <Card className="hover:border-foreground/20 transition-colors cursor-pointer">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-2xl">{post.title}</CardTitle>
                    <time className="text-sm text-muted-foreground">
                      {format(new Date(post.created_at), 'MMM dd, yyyy')}
                    </time>
                  </div>
                  <CardDescription className="text-base mt-2">{post.preview}</CardDescription>
                </CardHeader>
              </Card>
            </Link>
          ))
        )}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
          >
            <ChevronLeft className="h-4 w-4 mr-1" />
            Previous
          </Button>
          <span className="text-sm text-muted-foreground px-4">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
          >
            Next
            <ChevronRight className="h-4 w-4 ml-1" />
          </Button>
        </div>
      )}
    </div>
  );
}
