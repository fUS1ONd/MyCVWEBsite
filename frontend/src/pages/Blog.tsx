import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import axiosInstance from '@/lib/axios';
import { Post, PaginatedResponse } from '@/lib/types';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { format } from 'date-fns';
import { ChevronLeft, ChevronRight } from 'lucide-react';

const mockPosts: Post[] = [
  {
    id: 1,
    title: 'Введение в современную веб-разработку',
    slug: 'intro-to-modern-web-dev',
    content:
      '# Современная веб-разработка\n\nВеб-разработка прошла долгий путь за последние годы...',
    preview:
      'Обзор современных технологий и подходов в веб-разработке. Узнайте о React, TypeScript, и лучших практиках.',
    published: true,
    author_id: 1,
    created_at: new Date('2024-01-15').toISOString(),
    updated_at: new Date('2024-01-15').toISOString(),
  },
  {
    id: 2,
    title: 'Оптимизация производительности React приложений',
    slug: 'react-performance-optimization',
    content: '# Оптимизация React\n\nПроизводительность критична для пользовательского опыта...',
    preview:
      'Практические советы по улучшению производительности React приложений. Мемоизация, виртуализация и другие техники.',
    published: true,
    author_id: 1,
    created_at: new Date('2024-01-10').toISOString(),
    updated_at: new Date('2024-01-10').toISOString(),
  },
  {
    id: 3,
    title: 'TypeScript: От основ к продвинутым техникам',
    slug: 'typescript-advanced-techniques',
    content: '# TypeScript продвинутый уровень\n\nРазберем сложные паттерны TypeScript...',
    preview:
      'Глубокое погружение в TypeScript. Дженерики, условные типы, mapped types и другие продвинутые возможности.',
    published: true,
    author_id: 1,
    created_at: new Date('2024-01-05').toISOString(),
    updated_at: new Date('2024-01-05').toISOString(),
  },
];

const mockData: PaginatedResponse<Post> = {
  data: mockPosts,
  page: 1,
  limit: 10,
  total: 3,
  total_pages: 1,
};

export default function Blog() {
  const [page, setPage] = useState(1);
  const limit = 10;

  const { data, isLoading } = useQuery({
    queryKey: ['posts', page, limit],
    queryFn: async () => {
      try {
        const response = await axiosInstance.get<PaginatedResponse<Post>>(
          `/api/v1/posts?page=${page}&limit=${limit}&published=true`
        );
        return response.data;
      } catch (error) {
        return mockData;
      }
    },
  });

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
