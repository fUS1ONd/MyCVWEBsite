import React from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { postsService } from '../services/posts.service';
import { MarkdownContent, CommentsSection } from '../components/blog';
import { SEO } from '../components/common/SEO';
import { PageSpinner } from '../components/common/Spinner';
import { Button } from '../components/common/Button';

export const PostDetailPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();

  const {
    data: post,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['post', slug],
    queryFn: () => postsService.getPostBySlug(slug!),
    enabled: !!slug,
  });

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (isLoading) {
    return <PageSpinner />;
  }

  if (error || !post) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Пост не найден</h1>
          <p className="text-gray-600 mb-8">
            К сожалению, запрашиваемый пост не существует или был удален.
          </p>
          <Button onClick={() => navigate('/blog')}>Вернуться к списку постов</Button>
        </div>
      </div>
    );
  }

  return (
    <>
      <SEO title={post.title} description={post.preview} />
      <article className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Back button */}
        <div className="mb-6">
          <Link
            to="/blog"
            className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
          >
            ← Назад к списку постов
          </Link>
        </div>

        {/* Post header */}
        <header className="mb-8">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">{post.title}</h1>

          {/* Metadata */}
          <div className="flex items-center gap-4 text-gray-600">
            {post.author && (
              <div className="flex items-center gap-2">
                {post.author.avatar_url && (
                  <img
                    src={post.author.avatar_url}
                    alt={post.author.name}
                    className="w-10 h-10 rounded-full"
                  />
                )}
                <span className="font-medium">{post.author.name}</span>
              </div>
            )}
            <span>•</span>
            <time dateTime={post.published_at || post.created_at}>
              {formatDate(post.published_at || post.created_at)}
            </time>
          </div>
        </header>

        {/* Post content */}
        <div className="mb-12">
          <MarkdownContent content={post.content} />
        </div>

        {/* Share buttons */}
        <div className="border-t border-gray-200 pt-8 mb-12">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Поделиться статьей</h3>
          <div className="flex gap-4">
            <a
              href={`https://vk.com/share.php?url=${encodeURIComponent(window.location.href)}&title=${encodeURIComponent(post.title)}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              ВКонтакте
            </a>
            <a
              href={`https://twitter.com/intent/tweet?url=${encodeURIComponent(window.location.href)}&text=${encodeURIComponent(post.title)}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600 transition-colors"
            >
              Twitter
            </a>
            <a
              href={`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(window.location.href)}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-4 py-2 bg-blue-700 text-white rounded-lg hover:bg-blue-800 transition-colors"
            >
              Facebook
            </a>
            <a
              href={`https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(window.location.href)}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-4 py-2 bg-blue-800 text-white rounded-lg hover:bg-blue-900 transition-colors"
            >
              LinkedIn
            </a>
          </div>
        </div>

        {/* Comments section */}
        <div className="border-t border-gray-200 pt-8">
          <CommentsSection postSlug={post.slug} />
        </div>
      </article>
    </>
  );
};
