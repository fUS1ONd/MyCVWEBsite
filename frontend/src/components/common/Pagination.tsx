import React from 'react';
import { Button } from './Button';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
}) => {
  if (totalPages <= 1) {
    return null;
  }

  const pages: (number | string)[] = [];
  const maxVisiblePages = 7;

  if (totalPages <= maxVisiblePages) {
    // Show all pages if total is small
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i);
    }
  } else {
    // Always show first page
    pages.push(1);

    if (currentPage > 3) {
      pages.push('...');
    }

    // Show pages around current
    const start = Math.max(2, currentPage - 1);
    const end = Math.min(totalPages - 1, currentPage + 1);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (currentPage < totalPages - 2) {
      pages.push('...');
    }

    // Always show last page
    pages.push(totalPages);
  }

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      {/* Previous button */}
      <Button
        variant="secondary"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="px-4"
      >
        ← Назад
      </Button>

      {/* Page numbers */}
      <div className="flex gap-1">
        {pages.map((page, index) => {
          if (page === '...') {
            return (
              <span key={`ellipsis-${index}`} className="px-3 py-2 text-gray-400">
                ...
              </span>
            );
          }

          const pageNumber = page as number;
          const isActive = pageNumber === currentPage;

          return (
            <button
              key={pageNumber}
              onClick={() => onPageChange(pageNumber)}
              className={`
                min-w-[40px] h-10 px-3 rounded-lg font-medium transition-colors
                ${
                  isActive
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
                }
              `}
            >
              {pageNumber}
            </button>
          );
        })}
      </div>

      {/* Next button */}
      <Button
        variant="secondary"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="px-4"
      >
        Вперед →
      </Button>
    </div>
  );
};
