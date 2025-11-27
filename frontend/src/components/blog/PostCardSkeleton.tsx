import React from 'react';
import { Card } from '../common/Card';

export const PostCardSkeleton: React.FC = () => {
  return (
    <Card>
      <div className="space-y-3 animate-pulse">
        {/* Title skeleton */}
        <div className="h-8 bg-gray-200 rounded w-3/4"></div>

        {/* Metadata skeleton */}
        <div className="flex items-center gap-4">
          <div className="w-6 h-6 bg-gray-200 rounded-full"></div>
          <div className="h-4 bg-gray-200 rounded w-32"></div>
          <div className="h-4 bg-gray-200 rounded w-24"></div>
        </div>

        {/* Preview skeleton */}
        <div className="space-y-2">
          <div className="h-4 bg-gray-200 rounded"></div>
          <div className="h-4 bg-gray-200 rounded"></div>
          <div className="h-4 bg-gray-200 rounded w-2/3"></div>
        </div>

        {/* Read more link skeleton */}
        <div className="h-4 bg-gray-200 rounded w-32 mt-2"></div>
      </div>
    </Card>
  );
};
