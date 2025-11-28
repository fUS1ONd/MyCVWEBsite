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

const mockPost: Post = {
  id: 1,
  title: 'Введение в современную веб-разработку',
  slug: 'intro-to-modern-web-dev',
  content: `# Современная веб-разработка

Веб-разработка прошла долгий путь за последние годы. От простых HTML-страниц мы перешли к сложным, интерактивным приложениям, которые работают в браузере.

## Ключевые технологии

### React
React произвел революцию в том, как мы создаем пользовательские интерфейсы. Компонентный подход позволяет создавать переиспользуемые блоки кода.

### TypeScript
TypeScript добавляет статическую типизацию к JavaScript, что делает код более надежным и легким в поддержке.

### Современные инструменты
- Vite для быстрой сборки
- TanStack Query для управления состоянием сервера
- Tailwind CSS для стилизации

## Заключение
Современная веб-разработка - это увлекательная область, которая постоянно развивается. Главное - не бояться пробовать новые технологии и подходы.`,
  preview:
    'Обзор современных технологий и подходов в веб-разработке. Узнайте о React, TypeScript, и лучших практиках.',
  published: true,
  author_id: 1,
  author: {
    id: 1,
    email: 'alex@example.com',
    name: 'Александр Петров',
    role: 'admin',
    created_at: new Date().toISOString(),
  },
  created_at: new Date('2024-01-15').toISOString(),
  updated_at: new Date('2024-01-15').toISOString(),
};

const mockComments: Comment[] = [
  {
    id: 1,
    content: 'Отличная статья! Очень понравился раздел про TypeScript.',
    post_id: 1,
    user_id: 2,
    user: {
      id: 2,
      email: 'user@example.com',
      name: 'Мария Иванова',
      role: 'user',
      created_at: new Date().toISOString(),
    },
    created_at: new Date('2024-01-16').toISOString(),
    updated_at: new Date('2024-01-16').toISOString(),
    replies: [
      {
        id: 2,
        content: 'Спасибо за отзыв! Рад, что статья была полезной.',
        post_id: 1,
        user_id: 1,
        parent_id: 1,
        user: {
          id: 1,
          email: 'alex@example.com',
          name: 'Александр Петров',
          role: 'admin',
          created_at: new Date().toISOString(),
        },
        created_at: new Date('2024-01-16').toISOString(),
        updated_at: new Date('2024-01-16').toISOString(),
      },
    ],
  },
  {
    id: 3,
    content: 'Было бы интересно увидеть больше примеров кода с React hooks.',
    post_id: 1,
    user_id: 3,
    user: {
      id: 3,
      email: 'dev@example.com',
      name: 'Дмитрий Смирнов',
      role: 'user',
      created_at: new Date().toISOString(),
    },
    created_at: new Date('2024-01-17').toISOString(),
    updated_at: new Date('2024-01-17').toISOString(),
  },
];

export default function Article() {
  const { slug } = useParams<{ slug: string }>();

  const { data: post, isLoading: postLoading } = useQuery({
    queryKey: ['post', slug],
    queryFn: async () => {
      try {
        const response = await axiosInstance.get<Post>(`/api/v1/posts/${slug}`);
        return response.data;
      } catch (error) {
        return mockPost;
      }
    },
    enabled: !!slug,
  });

  const { data: comments, isLoading: commentsLoading } = useQuery({
    queryKey: ['comments', slug],
    queryFn: async () => {
      try {
        const response = await axiosInstance.get<Comment[]>(`/api/v1/posts/${slug}/comments`);
        return response.data;
      } catch (error) {
        return mockComments;
      }
    },
    enabled: !!slug,
  });

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
