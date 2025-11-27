import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from '../components/layout';
import { ProtectedRoute, AdminRoute } from '../components/auth';
import { HomePage, LoginPage, BlogPage, NotFoundPage } from '../pages';
import { PostDetailPage } from '../pages/PostDetailPage';

export const AppRoutes: React.FC = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          {/* Public routes */}
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />

          {/* Protected routes (require authentication) */}
          <Route
            path="/blog"
            element={
              <ProtectedRoute>
                <BlogPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/blog/:slug"
            element={
              <ProtectedRoute>
                <PostDetailPage />
              </ProtectedRoute>
            }
          />

          {/* Admin routes */}
          <Route
            path="/admin"
            element={
              <AdminRoute>
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                  <h1 className="text-3xl font-bold text-gray-900">Административная панель</h1>
                  <p className="text-gray-600 mt-4">
                    Placeholder для админ панели. Будет реализовано на этапе 8.
                  </p>
                </div>
              </AdminRoute>
            }
          />

          {/* 404 */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
};
