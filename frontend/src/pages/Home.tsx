import { useQuery } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Profile } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Github, Linkedin, Twitter, Globe, Mail } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';

const mockProfile: Profile = {
  name: 'Александр Петров',
  title: 'Full-Stack разработчик',
  description: 'Создаю современные веб-приложения с фокусом на производительность и UX',
  bio: 'Опытный разработчик с более чем 5 годами опыта в создании масштабируемых веб-приложений. Специализируюсь на React, TypeScript, Go и облачных технологиях.\n\nУвлечен созданием интуитивных пользовательских интерфейсов и надежных backend-систем. Постоянно изучаю новые технологии и лучшие практики разработки.',
  avatar_url: 'https://api.dicebear.com/7.x/avataaars/svg?seed=alex',
  github_url: 'https://github.com',
  linkedin_url: 'https://linkedin.com',
  twitter_url: 'https://twitter.com',
  website_url: 'https://example.com',
  email: 'alex@example.com',
};

export default function Home() {
  const { data: profile, isLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      try {
        const response = await axiosInstance.get<Profile>('/api/v1/profile');
        return response.data;
      } catch (error) {
        return mockProfile;
      }
    },
  });

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

  const displayProfile = profile || mockProfile;

  return (
    <div className="container max-w-4xl py-16 space-y-12">
      {/* Hero Section */}
      <div className="flex flex-col items-center text-center space-y-6">
        <Avatar className="h-32 w-32 border-2 border-border">
          <AvatarImage src={displayProfile.avatar_url} alt={displayProfile.name} />
          <AvatarFallback className="text-3xl">{displayProfile.name.charAt(0)}</AvatarFallback>
        </Avatar>

        <div className="space-y-2">
          <h1 className="text-4xl font-bold tracking-tight">{displayProfile.name}</h1>
          {displayProfile.title && (
            <p className="text-xl text-muted-foreground">{displayProfile.title}</p>
          )}
        </div>

        {displayProfile.description && (
          <p className="text-lg text-muted-foreground max-w-2xl">{displayProfile.description}</p>
        )}

        {/* Social Links */}
        <div className="flex flex-wrap justify-center gap-2">
          {displayProfile.github_url && (
            <Button variant="outline" size="sm" asChild>
              <a href={displayProfile.github_url} target="_blank" rel="noopener noreferrer">
                <Github className="h-4 w-4 mr-2" />
                GitHub
              </a>
            </Button>
          )}
          {displayProfile.linkedin_url && (
            <Button variant="outline" size="sm" asChild>
              <a href={displayProfile.linkedin_url} target="_blank" rel="noopener noreferrer">
                <Linkedin className="h-4 w-4 mr-2" />
                LinkedIn
              </a>
            </Button>
          )}
          {displayProfile.twitter_url && (
            <Button variant="outline" size="sm" asChild>
              <a href={displayProfile.twitter_url} target="_blank" rel="noopener noreferrer">
                <Twitter className="h-4 w-4 mr-2" />
                Twitter
              </a>
            </Button>
          )}
          {displayProfile.website_url && (
            <Button variant="outline" size="sm" asChild>
              <a href={displayProfile.website_url} target="_blank" rel="noopener noreferrer">
                <Globe className="h-4 w-4 mr-2" />
                Website
              </a>
            </Button>
          )}
          {displayProfile.email && (
            <Button variant="outline" size="sm" asChild>
              <a href={`mailto:${displayProfile.email}`}>
                <Mail className="h-4 w-4 mr-2" />
                Email
              </a>
            </Button>
          )}
        </div>
      </div>

      {/* Bio Section */}
      {displayProfile.bio && (
        <Card>
          <CardContent className="pt-6">
            <h2 className="text-2xl font-semibold mb-4">О себе</h2>
            <div className="prose prose-slate dark:prose-invert max-w-none">
              <p className="whitespace-pre-wrap text-muted-foreground">{displayProfile.bio}</p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
