import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { gtmPageView } from '@/lib/gtm';

export const usePageTracking = () => {
  const location = useLocation();

  useEffect(() => {
    // Send page view event on route change
    gtmPageView(location.pathname + location.search);
  }, [location]);
};
