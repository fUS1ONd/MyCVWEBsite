import { useEffect } from 'react';
import { useInfiniteQuery } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { PostListResponse } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useToast } from '@/hooks/use-toast';
import { useInView } from 'react-intersection-observer';
import { PostCard } from '@/components/blog/PostCard';
import { Loader2 } from 'lucide-react';

export default function Blog() {
  const { toast } = useToast();
  const limit = 10;
  const { ref, inView } = useInView();

  const { data, isLoading, error, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery({
      queryKey: ['posts'],
      queryFn: async ({ pageParam = 1 }) => {
        const response = await axiosInstance.get<PostListResponse>(
          `/api/v1/posts?page=${pageParam}&limit=${limit}&published=true`
        );
        return response.data;
      },
      getNextPageParam: (lastPage) => {
        return lastPage.page < lastPage.total_pages ? lastPage.page + 1 : undefined;
      },
      initialPageParam: 1,
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

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  const posts = data?.pages.flatMap((page) => page.posts) || [];

  if (isLoading) {
    return (
      <div className="container max-w-4xl py-16 space-y-8">
        <div className="space-y-2">
          <Skeleton className="h-10 w-64" />
          <Skeleton className="h-4 w-96" />
        </div>
        <div className="space-y-6">
          {[...Array(3)].map((_, i) => (
            <Card key={i} className="overflow-hidden">
              <Skeleton className="h-48 w-full" />
              <div className="p-6 space-y-4">
                <Skeleton className="h-8 w-3/4" />
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-2/3" />
              </div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="container max-w-4xl py-16 space-y-8">
      <div className="space-y-2">
        <h1 className="text-2xl sm:text-4xl font-bold tracking-tight">Blog</h1>
        <p className="text-base sm:text-lg text-muted-foreground">
          Thoughts, stories and ideas about AI
        </p>
      </div>

      <div className="space-y-6">
        {posts.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">No posts published yet</p>
            </CardContent>
          </Card>
        ) : (
          <>
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}

            {/* Infinite scroll trigger */}
            {hasNextPage && (
              <div ref={ref} className="flex justify-center py-8">
                {isFetchingNextPage && (
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Loader2 className="h-5 w-5 animate-spin" />
                    <span>Loading more posts...</span>
                  </div>
                )}
              </div>
            )}

            {!hasNextPage && posts.length > 0 && (
              <div className="text-center py-8 text-muted-foreground">
                <p>You've reached the end</p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
