import { useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Profile, ApiResponse } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Github, Mail } from 'lucide-react';
import { FaTelegram } from 'react-icons/fa';
import { Skeleton } from '@/components/ui/skeleton';
import { useToast } from '@/hooks/use-toast';

export default function Home() {
  const { toast } = useToast();
  const {
    data: profile,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<Profile>>('/api/v1/profile');
      return response.data.data;
    },
  });

  useEffect(() => {
    if (error) {
      toast({
        title: 'Ошибка загрузки профиля',
        description: error instanceof Error ? error.message : 'Не удалось подключиться к серверу',
        variant: 'destructive',
      });
    }
  }, [error, toast]);

  if (isLoading) {
    return (
      <div className="container max-w-4xl py-16 space-y-8">
        <div className="flex flex-col items-center text-center space-y-4">
          <Skeleton className="h-32 w-32 rounded-full" />
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-4 w-48" />
        </div>
      </div>
    );
  }

  if (!profile) {
    return null;
  }

  return (
    <div className="container max-w-4xl py-16 space-y-12">
      {/* Hero Section */}
      <div className="flex flex-col items-center text-center space-y-6">
        <Avatar className="h-32 w-32 border-2 border-border">
          <AvatarImage src={profile.photo_url} alt={profile.name} />
          <AvatarFallback className="text-3xl">{profile.name.charAt(0)}</AvatarFallback>
        </Avatar>

        <div className="space-y-2">
          <h1 className="text-4xl font-bold tracking-tight">{profile.name}</h1>
        </div>

        <p className="text-lg text-muted-foreground max-w-2xl">{profile.description}</p>

        {/* Social Links */}
        <div className="flex flex-wrap justify-center gap-2">
          {profile.contacts.github && (
            <Button variant="outline" size="sm" asChild>
              <a href={profile.contacts.github} target="_blank" rel="noopener noreferrer">
                <Github className="h-4 w-4 mr-2" />
                GitHub
              </a>
            </Button>
          )}
          {profile.contacts.telegram && (
            <Button variant="outline" size="sm" asChild>
              <a href={profile.contacts.telegram} target="_blank" rel="noopener noreferrer">
                <FaTelegram className="h-4 w-4 mr-2" />
                Telegram
              </a>
            </Button>
          )}
          {profile.contacts.email && (
            <Button variant="outline" size="sm" asChild>
              <a href={`mailto:${profile.contacts.email}`}>
                <Mail className="h-4 w-4 mr-2" />
                Email
              </a>
            </Button>
          )}
        </div>
      </div>

      {/* Activity Section */}
      <Card>
        <CardContent className="pt-6">
          <h2 className="text-2xl font-semibold mb-4">Деятельность</h2>
          <div className="prose prose-slate dark:prose-invert max-w-none">
            <p className="whitespace-pre-wrap break-words [overflow-wrap:anywhere] text-muted-foreground">
              {profile.activity}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
