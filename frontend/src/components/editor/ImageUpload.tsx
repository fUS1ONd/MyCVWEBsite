import { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { useMutation } from '@tanstack/react-query';
import axiosInstance from '@/lib/axios';
import { Button } from '@/components/ui/button';
import { toast } from '@/hooks/use-toast';
import { Upload, X, Loader2, Image as ImageIcon } from 'lucide-react';
import { MediaFile } from '@/lib/types';

interface ImageUploadProps {
  onUpload: (url: string, mediaFile: MediaFile) => void;
  accept?: Record<string, string[]>;
  maxSize?: number;
  label?: string;
}

export function ImageUpload({
  onUpload,
  accept = {
    'image/jpeg': ['.jpg', '.jpeg'],
    'image/png': ['.png'],
    'image/webp': ['.webp'],
    'image/gif': ['.gif'],
  },
  maxSize = 10 * 1024 * 1024, // 10MB
  label = 'Upload Image',
}: ImageUploadProps) {
  const [preview, setPreview] = useState<string | null>(null);

  const uploadMutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);

      const response = await axiosInstance.post<MediaFile>('/api/v1/admin/media', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data;
    },
    onSuccess: (data) => {
      toast({
        title: 'Image uploaded successfully',
        description: data.filename,
      });
      onUpload(data.url, data);
      setPreview(null);
    },
    onError: (error: unknown) => {
      const errorMessage =
        error instanceof Error
          ? error.message
          : typeof error === 'object' && error !== null && 'response' in error
            ? String(
                (error as { response?: { data?: { error?: { message?: string } } } }).response?.data
                  ?.error?.message || ''
              )
            : '';
      toast({
        title: 'Failed to upload image',
        description: errorMessage || 'An error occurred',
        variant: 'destructive',
      });
      setPreview(null);
    },
  });

  const onDrop = useCallback(
    (acceptedFiles: File[]) => {
      const file = acceptedFiles[0];
      if (!file) return;

      // Show preview
      const reader = new FileReader();
      reader.onload = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);

      // Upload
      uploadMutation.mutate(file);
    },
    [uploadMutation]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept,
    maxSize,
    multiple: false,
    disabled: uploadMutation.isPending,
  });

  const clearPreview = () => {
    setPreview(null);
  };

  return (
    <div className="space-y-2">
      <div
        {...getRootProps()}
        className={`
          border-2 border-dashed rounded-lg p-6 text-center cursor-pointer
          transition-colors duration-200
          ${isDragActive ? 'border-primary bg-primary/5' : 'border-muted-foreground/25'}
          ${uploadMutation.isPending ? 'opacity-50 cursor-not-allowed' : 'hover:border-primary'}
        `}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center gap-2">
          {uploadMutation.isPending ? (
            <>
              <Loader2 className="h-10 w-10 text-muted-foreground animate-spin" />
              <p className="text-sm text-muted-foreground">Uploading...</p>
            </>
          ) : preview ? (
            <div className="relative">
              <img src={preview} alt="Preview" className="max-h-48 rounded-lg" />
              <Button
                type="button"
                variant="destructive"
                size="icon"
                className="absolute top-2 right-2"
                onClick={(e) => {
                  e.stopPropagation();
                  clearPreview();
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ) : (
            <>
              {isDragActive ? (
                <>
                  <Upload className="h-10 w-10 text-primary" />
                  <p className="text-sm text-primary font-medium">Drop the image here</p>
                </>
              ) : (
                <>
                  <ImageIcon className="h-10 w-10 text-muted-foreground" />
                  <div>
                    <p className="text-sm font-medium">{label}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Drag & drop or click to select
                    </p>
                    <p className="text-xs text-muted-foreground">
                      Max size: {(maxSize / (1024 * 1024)).toFixed(0)}MB
                    </p>
                  </div>
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
