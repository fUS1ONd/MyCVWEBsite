import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Github } from 'lucide-react';
import { FcGoogle } from 'react-icons/fc';

export default function Login() {
  const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080';

  const handleOAuthLogin = (provider: 'google' | 'github') => {
    window.location.href = `${backendUrl}/auth/${provider}`;
  };

  return (
    <div className="container flex items-center justify-center min-h-[calc(100vh-4rem)] py-16">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <CardTitle className="text-3xl font-bold">Welcome back</CardTitle>
          <CardDescription>Sign in to your account to continue</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button
            variant="outline"
            className="w-full h-12"
            onClick={() => handleOAuthLogin('google')}
          >
            <FcGoogle className="mr-2 h-5 w-5" />
            Continue with Google
          </Button>

          <Button
            variant="outline"
            className="w-full h-12"
            onClick={() => handleOAuthLogin('github')}
          >
            <Github className="mr-2 h-5 w-5" />
            Continue with GitHub
          </Button>

          <p className="text-xs text-center text-muted-foreground mt-6">
            By continuing, you agree to our Terms of Service and Privacy Policy
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
