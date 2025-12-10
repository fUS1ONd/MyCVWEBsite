export const GTM_CONSENT_KEY = 'cookie_consent';

declare global {
  interface Window {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    dataLayer: any[];
  }
}

export const enableAnalytics = () => {
  if (typeof window !== 'undefined' && window.dataLayer) {
    window.dataLayer.push({
      event: 'analytics_active',
    });
  }
};

export const initializeGTM = () => {
  if (typeof window !== 'undefined') {
    const consent = localStorage.getItem(GTM_CONSENT_KEY);
    if (consent === 'true') {
      enableAnalytics();
    }
  }
};

export const gtmPageView = (url: string) => {
  if (typeof window !== 'undefined' && window.dataLayer) {
    window.dataLayer.push({
      event: 'page_view',
      page_path: url,
    });
  }
};
