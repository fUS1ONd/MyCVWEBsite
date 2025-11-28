import { Link, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Navigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { FileText, User, Home } from 'lucide-react';

export default function Dashboard() {
  const { user, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return <div className="container py-16">Loading...</div>;
  }

  if (!user || user.role !== 'admin') {
    return <Navigate to="/login" replace />;
  }

  const isActive = (path: string) => location.pathname === path;

  return (
    <div className="min-h-screen bg-muted/30">
      <div className="border-b bg-background">
        <div className="container flex h-16 items-center gap-6">
          <Link to="/" className="flex items-center gap-2">
            <Home className="h-5 w-5" />
            <span className="font-semibold">Back to Site</span>
          </Link>
          <div className="flex-1" />
          <nav className="flex items-center gap-2">
            <Button variant={isActive('/admin') ? 'secondary' : 'ghost'} size="sm" asChild>
              <Link to="/admin">
                <User className="h-4 w-4 mr-2" />
                Profile
              </Link>
            </Button>
            <Button variant={isActive('/admin/posts') ? 'secondary' : 'ghost'} size="sm" asChild>
              <Link to="/admin/posts">
                <FileText className="h-4 w-4 mr-2" />
                Posts
              </Link>
            </Button>
          </nav>
        </div>
      </div>
      <div className="container py-8">
        <Outlet />
      </div>
    </div>
  );
}
