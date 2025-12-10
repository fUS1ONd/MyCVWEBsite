import { useEffect } from 'react';
import { Toaster } from '@/components/ui/toaster';
import { Toaster as Sonner } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from '@/contexts/AuthContext';
import { ThemeProvider } from '@/contexts/ThemeContext';
import { Layout } from '@/components/layout/Layout';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';
import { usePageTracking } from '@/hooks/usePageTracking';
import { CookieConsent } from '@/components/CookieConsent';
import { initializeGTM } from '@/lib/gtm';

// Pages
import Home from './pages/Home';
import Blog from './pages/Blog';
import Article from './pages/Article';
import Login from './pages/Login';
import NotFound from './pages/NotFound';

// Admin Pages
import Dashboard from './pages/admin/Dashboard';
import ProfileEditor from './pages/admin/ProfileEditor';
import PostManager from './pages/admin/PostManager';
import PostEditor from './pages/admin/PostEditor';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
});

// Helper component to use hook inside Router context
function PageTracker() {
  usePageTracking();
  return null;
}

const App = () => {
  useEffect(() => {
    initializeGTM();
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <TooltipProvider>
            <Toaster />
            <Sonner />
            <CookieConsent />
            <BrowserRouter>
              <PageTracker />
              <Routes>
                {/* Public Routes */}
                <Route
                  path="/"
                  element={
                    <Layout>
                      <Home />
                    </Layout>
                  }
                />

                {/* Protected Blog Routes */}
                <Route element={<ProtectedRoute />}>
                  <Route
                    path="/blog"
                    element={
                      <Layout>
                        <Blog />
                      </Layout>
                    }
                  />
                  <Route
                    path="/blog/:slug"
                    element={
                      <Layout>
                        <Article />
                      </Layout>
                    }
                  />
                </Route>

                <Route
                  path="/login"
                  element={
                    <Layout>
                      <Login />
                    </Layout>
                  }
                />

                {/* Admin Routes */}
                <Route path="/admin" element={<Dashboard />}>
                  <Route index element={<ProfileEditor />} />
                  <Route path="posts" element={<PostManager />} />
                  <Route path="posts/:id" element={<PostEditor />} />
                </Route>

                {/* Catch-all */}
                <Route path="*" element={<NotFound />} />
              </Routes>
            </BrowserRouter>
          </TooltipProvider>
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
};

export default App;
