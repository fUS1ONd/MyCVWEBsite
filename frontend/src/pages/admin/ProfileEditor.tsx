import { useState, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import axiosInstance from '@/lib/axios';
import { Profile, ApiResponse } from '@/lib/types';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { toast } from '@/hooks/use-toast';
import { Skeleton } from '@/components/ui/skeleton';
import { Upload, X } from 'lucide-react';

export default function ProfileEditor() {
  const queryClient = useQueryClient();
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { data: profile, isLoading } = useQuery({
    queryKey: ['admin', 'profile'],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<Profile>>('/api/v1/profile');
      return response.data.data;
    },
  });

  const {
    register,
    setValue,
    watch,
    handleSubmit,
    formState: { isDirty },
  } = useForm<Profile>({
    values: profile,
  });

  const photoUrl = watch('photo_url');

  const mutation = useMutation({
    mutationFn: async (data: Profile) => {
      await axiosInstance.put('/api/v1/admin/profile', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'profile'] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      toast({ title: 'Profile updated successfully' });
    },
    onError: () => {
      toast({ title: 'Failed to update profile', variant: 'destructive' });
    },
  });

  const onSubmit = (data: Profile) => {
    mutation.mutate(data);
  };

  const handleUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    setIsUploading(true);
    try {
      const response = await axiosInstance.post<{ success: boolean; data: { url: string } }>(
        '/api/v1/admin/upload',
        formData,
        { headers: { 'Content-Type': 'multipart/form-data' } }
      );
      setValue('photo_url', response.data.data.url, { shouldDirty: true });
      toast({ title: 'Photo uploaded' });
    } catch (error) {
      toast({ title: 'Failed to upload photo', variant: 'destructive' });
    } finally {
      setIsUploading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-3xl space-y-6">
        <Skeleton className="h-8 w-64" />
        <Card>
          <CardContent className="pt-6 space-y-4">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-32 w-full" />
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-3xl space-y-6">
      <div>
        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight">Profile Settings</h1>
        <p className="text-muted-foreground mt-2">Manage your public profile information</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>Your public profile details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Имя *</Label>
              <Input id="name" {...register('name')} required placeholder="Ваше имя" />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Описание *</Label>
              <Textarea
                id="description"
                {...register('description')}
                placeholder="Краткое описание ваших навыков и специализации"
                className="min-h-[100px]"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="activity">Деятельность *</Label>
              <Textarea
                id="activity"
                {...register('activity')}
                placeholder="Расскажите о своей текущей деятельности и опыте"
                className="min-h-[120px]"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="photo_url">URL фотографии</Label>
              <div className="flex gap-2">
                <Input id="photo_url" {...register('photo_url')} placeholder="https://..." />
                <input
                  type="file"
                  className="hidden"
                  ref={fileInputRef}
                  onChange={handleUpload}
                  accept="image/png, image/jpeg, image/webp"
                />
                <Button
                  type="button"
                  variant="outline"
                  disabled={isUploading}
                  onClick={() => fileInputRef.current?.click()}
                >
                  <Upload className="h-4 w-4 mr-2" />
                  {isUploading ? 'Uploading...' : 'Upload Photo'}
                </Button>
                {photoUrl && (
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    onClick={() => setValue('photo_url', '', { shouldDirty: true })}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>

            <div className="border-t pt-6">
              <h3 className="text-sm font-medium mb-4">Контакты</h3>
              <div className="grid gap-4">
                <div className="space-y-2">
                  <Label htmlFor="contacts.email">Email *</Label>
                  <Input
                    id="contacts.email"
                    {...register('contacts.email')}
                    type="email"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="contacts.github">GitHub</Label>
                  <Input
                    id="contacts.github"
                    {...register('contacts.github')}
                    type="url"
                    placeholder="https://github.com/username"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="contacts.telegram">Telegram</Label>
                  <Input
                    id="contacts.telegram"
                    {...register('contacts.telegram')}
                    type="url"
                    placeholder="https://t.me/username"
                  />
                </div>
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <Button type="submit" disabled={!isDirty || mutation.isPending}>
                {mutation.isPending ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  );
}
