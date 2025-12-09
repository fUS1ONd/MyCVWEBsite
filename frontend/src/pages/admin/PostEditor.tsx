import { useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import axiosInstance from '@/lib/axios';
import { Post, MediaFile } from '@/lib/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from '@/hooks/use-toast';
import { ChevronLeft, Save, X } from 'lucide-react';
import { TipTapEditor } from '@/components/editor/TipTapEditor';
import { ImageUpload } from '@/components/editor/ImageUpload';

type PostFormData = {
  title: string;
  slug: string;
  content: string;
  preview: string;
  published: boolean;
  cover_image_id?: number;
};

export default function PostEditor() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isNew = id === 'new';

  const { data: post, isLoading } = useQuery({
    queryKey: ['admin', 'post', id],
    queryFn: async () => {
      const response = await axiosInstance.get<Post>(`/api/v1/admin/posts/${id}`);
      return response.data;
    },
    enabled: !isNew,
  });

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { isDirty },
  } = useForm<PostFormData>({
    defaultValues: {
      title: '',
      slug: '',
      content: '',
      preview: '',
      published: false,
      cover_image_id: undefined,
    },
  });

  useEffect(() => {
    if (post) {
      setValue('title', post.title);
      setValue('slug', post.slug);
      setValue('content', post.content);
      setValue('preview', post.preview);
      setValue('published', post.published);
      setValue('cover_image_id', post.cover_image_id);
    }
  }, [post, setValue]);

  const mutation = useMutation({
    mutationFn: async (data: PostFormData) => {
      if (isNew) {
        await axiosInstance.post('/api/v1/admin/posts', data);
      } else {
        await axiosInstance.put(`/api/v1/admin/posts/${id}`, data);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      toast({ title: isNew ? 'Post created' : 'Post updated' });
      navigate('/admin/posts');
    },
    onError: () => {
      toast({ title: 'Failed to save post', variant: 'destructive' });
    },
  });

  const onSubmit = (data: PostFormData) => {
    if (!data.slug) {
      data.slug = data.title.toLowerCase().replace(/[^a-z0-9]+/g, '-');
    }
    mutation.mutate(data);
  };

  const handleCoverImageUpload = (url: string, mediaFile: MediaFile) => {
    setValue('cover_image_id', mediaFile.id, { shouldDirty: true });
  };

  const removeCoverImage = () => {
    setValue('cover_image_id', undefined, { shouldDirty: true });
  };

  const coverImageId = watch('cover_image_id');
  const coverImage = post?.cover_image;

  if (isLoading && !isNew) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-64" />
        <Card>
          <CardContent className="pt-6 space-y-4">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-64 w-full" />
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="sm" asChild>
            <Link to="/admin/posts">
              <ChevronLeft className="h-4 w-4 mr-2" />
              Back
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {isNew ? 'New Post' : 'Edit Post'}
            </h1>
          </div>
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Post Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="title">Title *</Label>
                <Input id="title" {...register('title')} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="slug">Slug</Label>
                <Input id="slug" {...register('slug')} placeholder="auto-generated-from-title" />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="preview">Preview Text *</Label>
              <Textarea
                id="preview"
                {...register('preview')}
                placeholder="A short description of the post"
                className="min-h-[80px]"
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Cover Image</Label>
              {coverImageId && coverImage ? (
                <div className="relative">
                  <img
                    src={coverImage.url}
                    alt="Cover"
                    className="w-full h-48 object-cover rounded-lg"
                  />
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    className="absolute top-2 right-2"
                    onClick={removeCoverImage}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              ) : (
                <ImageUpload onUpload={handleCoverImageUpload} label="Upload Cover Image" />
              )}
            </div>

            <div className="flex items-center space-x-2">
              <Switch
                id="published"
                checked={watch('published')}
                onCheckedChange={(checked) => setValue('published', checked, { shouldDirty: true })}
              />
              <Label htmlFor="published">Published</Label>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Content</CardTitle>
          </CardHeader>
          <CardContent>
            <TipTapEditor
              content={watch('content')}
              onChange={(content) => setValue('content', content, { shouldDirty: true })}
              placeholder="Write your post content..."
            />
          </CardContent>
        </Card>

        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" asChild>
            <Link to="/admin/posts">Cancel</Link>
          </Button>
          <Button type="submit" disabled={!isDirty || mutation.isPending}>
            <Save className="h-4 w-4 mr-2" />
            {mutation.isPending ? 'Saving...' : 'Save Post'}
          </Button>
        </div>
      </form>
    </div>
  );
}
