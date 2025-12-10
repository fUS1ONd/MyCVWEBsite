import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { enableAnalytics, GTM_CONSENT_KEY } from '@/lib/gtm';

export function CookieConsent() {
  const [show, setShow] = useState(false);

  useEffect(() => {
    const consent = localStorage.getItem(GTM_CONSENT_KEY);
    if (consent === null) {
      const timer = setTimeout(() => setShow(true), 1000);
      return () => clearTimeout(timer);
    }
  }, []);

  const handleAccept = () => {
    localStorage.setItem(GTM_CONSENT_KEY, 'true');
    enableAnalytics();
    setShow(false);
  };

  const handleDecline = () => {
    localStorage.setItem(GTM_CONSENT_KEY, 'false');
    setShow(false);
  };

  if (!show) return null;

  return (
    <div className="fixed bottom-0 left-0 right-0 p-4 z-50 flex justify-center pointer-events-none">
      <Card className="w-full max-w-2xl pointer-events-auto shadow-xl border-primary/20 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <CardContent className="p-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="space-y-2 text-center sm:text-left">
            <h3 className="font-semibold tracking-tight">Мы используем файлы cookie</h3>
            <p className="text-sm text-muted-foreground">
              Этот сайт использует cookie для сбора обезличенной статистики и улучшения
              пользовательского опыта.
            </p>
          </div>
          <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
            <Button variant="outline" size="sm" onClick={handleDecline}>
              Отказаться
            </Button>
            <Button size="sm" onClick={handleAccept}>
              Принять
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
