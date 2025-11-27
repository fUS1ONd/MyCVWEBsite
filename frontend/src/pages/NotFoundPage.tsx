import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../components/common';

export const NotFoundPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-[calc(100vh-theme(spacing.32))] flex items-center justify-center px-4">
      <div className="text-center">
        <h1 className="text-9xl font-bold text-gray-200">404</h1>
        <h2 className="text-3xl font-semibold text-gray-900 mt-4 mb-2">Страница не найдена</h2>
        <p className="text-gray-600 mb-8">Запрашиваемая страница не существует или была удалена.</p>
        <div className="flex gap-3 justify-center">
          <Button variant="secondary" onClick={() => navigate(-1)}>
            Назад
          </Button>
          <Button variant="primary" onClick={() => navigate('/')}>
            На главную
          </Button>
        </div>
      </div>
    </div>
  );
};
