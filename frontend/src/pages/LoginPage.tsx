import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { authService } from '../services';
import { Button, Card } from '../components/common';

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as { from?: Location })?.from?.pathname || '/blog';

  const handleOAuthLogin = (provider: 'vk' | 'google' | 'github') => {
    // Save return URL in localStorage
    localStorage.setItem('auth_redirect', from);

    // Redirect to OAuth provider
    const oauthUrl = authService.getOAuthUrl(provider);
    window.location.href = oauthUrl;
  };

  return (
    <div className="min-h-[calc(100vh-theme(spacing.32))] flex items-center justify-center px-4">
      <Card className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Вход в систему</h1>
          <p className="text-gray-600">Войдите с помощью одного из провайдеров</p>
        </div>

        <div className="space-y-3">
          <Button variant="primary" className="w-full" onClick={() => handleOAuthLogin('vk')}>
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M15.07 2H8.93C3.33 2 2 3.33 2 8.93v6.14C2 20.67 3.33 22 8.93 22h6.14c5.6 0 6.93-1.33 6.93-6.93V8.93C22 3.33 20.67 2 15.07 2zm3.18 14.49h-1.6c-.46 0-.6-.37-1.43-1.2-.72-.67-1.04-.76-1.22-.76-.25 0-.32.07-.32.42v1.1c0 .3-.09.47-1.01.47-1.64 0-3.45-.99-4.73-2.84-1.92-2.59-2.44-4.54-2.44-4.93 0-.18.07-.35.42-.35h1.6c.31 0 .43.14.55.48.65 1.94 1.75 3.64 2.2 3.64.17 0 .25-.08.25-.5v-1.98c-.06-.97-.56-1.05-.56-1.4 0-.15.12-.3.32-.3h2.5c.26 0 .36.14.36.45v2.68c0 .26.12.36.19.36.17 0 .31-.1.62-.41 1.03-1.16 1.77-2.96 1.77-2.96.09-.21.23-.35.54-.35h1.6c.38 0 .46.19.38.45-.16.72-1.81 3.29-1.81 3.29-.14.23-.19.33 0 .59.13.19.56.55.85.88.53.56 1.07 1.03 1.19 1.36.13.33-.07.5-.4.5z" />
            </svg>
            Войти через ВКонтакте
          </Button>

          <Button variant="secondary" className="w-full" onClick={() => handleOAuthLogin('google')}>
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
            </svg>
            Войти через Google
          </Button>

          <Button variant="ghost" className="w-full" onClick={() => handleOAuthLogin('github')}>
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path
                fillRule="evenodd"
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                clipRule="evenodd"
              />
            </svg>
            Войти через GitHub
          </Button>
        </div>

        <div className="mt-6 text-center">
          <button
            onClick={() => navigate('/')}
            className="text-sm text-gray-600 hover:text-gray-900 transition-colors"
          >
            Вернуться на главную
          </button>
        </div>
      </Card>
    </div>
  );
};
