import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import path from 'path';
import { componentTagger } from 'lovable-tagger';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // Use environment variable for backend URL, fallback to localhost for local dev
  const backendUrl = process.env.VITE_BACKEND_URL || 'http://localhost:8080';

  return {
    server: {
      host: '::',
      port: 5173,
      proxy: {
        '/api': {
          target: backendUrl,
          changeOrigin: true,
        },
        '/auth': {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
    plugins: [react(), mode === 'development' && componentTagger()].filter(Boolean),
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  };
});
