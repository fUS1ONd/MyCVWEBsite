import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';

// Error boundary for catching load errors
window.addEventListener('error', (event) => {
  console.error('Global error:', event.error);
  console.error('Error loading:', event.filename, event.lineno, event.colno);
  console.error('Full event:', event);
});

// Unhandled promise rejections
window.addEventListener('unhandledrejection', (event) => {
  console.error('Unhandled promise rejection:', event.reason);
});

createRoot(document.getElementById('root')!).render(<App />);
