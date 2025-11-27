import React from 'react';

export const BlogPage: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Блог об искусственном интеллекте</h1>

      <p className="text-gray-600 mb-8">
        Здесь будет отображаться список постов. Эта страница требует авторизации.
      </p>

      <div className="bg-white rounded-lg shadow-md p-6">
        <p className="text-gray-700">
          Placeholder для списка постов. Будет реализовано на этапе 7.
        </p>
      </div>
    </div>
  );
};
