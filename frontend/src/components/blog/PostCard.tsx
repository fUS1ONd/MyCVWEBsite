import React from 'react';
import { Link } from 'react-router-dom';
import type { Post } from '../../types';
import { Card } from '../common/Card';

interface PostCardProps {
  post: Post;
}

export const PostCard: React.FC<PostCardProps> = ({ post }) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <Card className="hover:shadow-lg transition-shadow duration-200">
      <Link to={`/blog/${post.slug}`} className="block">
        <div className="space-y-3">
          {/* Title */}
          <h2 className="text-2xl font-bold text-gray-900 hover:text-blue-600 transition-colors">
            {post.title}
          </h2>

          {/* Metadata */}
          <div className="flex items-center gap-4 text-sm text-gray-500">
            {post.author && (
              <span className="flex items-center gap-2">
                {post.author.avatar_url && (
                  <img
                    src={post.author.avatar_url}
                    alt={post.author.name}
                    className="w-6 h-6 rounded-full"
                  />
                )}
                <span>{post.author.name}</span>
              </span>
            )}
            <span>•</span>
            <time dateTime={post.published_at || post.created_at}>
              {formatDate(post.published_at || post.created_at)}
            </time>
          </div>

          {/* Preview */}
          <p className="text-gray-600 line-clamp-3">{post.preview}</p>

          {/* Read more link */}
          <div className="pt-2">
            <span className="text-blue-600 font-medium hover:text-blue-700">Читать далее →</span>
          </div>
        </div>
      </Link>
    </Card>
  );
};
