import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { postsService } from '../services/posts.service';
import { PostCard, PostCardSkeleton } from '../components/blog';
import { Pagination } from '../components/common';
import { SEO } from '../components/common/SEO';

export const BlogPage: React.FC = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState('');
  const perPage = 10;

  const { data, isLoading, error } = useQuery({
    queryKey: ['posts', currentPage, searchQuery],
    queryFn: () =>
      postsService.getPosts({
        page: currentPage,
        per_page: perPage,
        published: true,
      }),
  });

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
    setCurrentPage(1); // Reset to first page on search
  };

  // Filter posts by search query on client side
  const filteredPosts =
    data?.posts.filter((post) => post.title.toLowerCase().includes(searchQuery.toLowerCase())) ||
    [];

  return (
    <>
      <SEO
        title="Блог об искусственном интеллекте"
        description="Статьи и заметки об искусственном интеллекте, машинном обучении и современных технологиях"
      />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Блог об искусственном интеллекте
          </h1>
          <p className="text-lg text-gray-600">
            Статьи и заметки об искусственном интеллекте, машинном обучении и современных
            технологиях
          </p>
        </div>

        {/* Search */}
        <div className="mb-8">
          <input
            type="text"
            placeholder="Поиск по заголовкам..."
            value={searchQuery}
            onChange={handleSearchChange}
            className="w-full max-w-md px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Error state */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-8">
            Ошибка при загрузке постов. Пожалуйста, попробуйте позже.
          </div>
        )}

        {/* Loading state */}
        {isLoading && (
          <div className="space-y-6">
            {Array.from({ length: 3 }).map((_, index) => (
              <PostCardSkeleton key={index} />
            ))}
          </div>
        )}

        {/* Empty state */}
        {!isLoading && !error && filteredPosts.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-600 text-lg">
              {searchQuery
                ? 'Посты не найдены по вашему запросу'
                : 'Пока нет опубликованных постов'}
            </p>
          </div>
        )}

        {/* Posts list */}
        {!isLoading && !error && filteredPosts.length > 0 && (
          <>
            <div className="space-y-6">
              {filteredPosts.map((post) => (
                <PostCard key={post.id} post={post} />
              ))}
            </div>

            {/* Pagination */}
            {data && data.total_pages > 1 && (
              <Pagination
                currentPage={currentPage}
                totalPages={data.total_pages}
                onPageChange={handlePageChange}
              />
            )}
          </>
        )}
      </div>
    </>
  );
};
