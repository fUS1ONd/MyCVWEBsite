import { useState } from 'react';
import { cn } from '@/lib/utils';
import { Skeleton } from '@/components/ui/skeleton';

interface SmartImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  aspectRatio?: string;
  containerClassName?: string;
}

export function SmartImage({
  src,
  alt,
  className,
  containerClassName,
  aspectRatio = 'aspect-video',
  ...props
}: SmartImageProps) {
  const [isLoaded, setIsLoaded] = useState(false);

  if (!src) return null;

  return (
    <div
      className={cn('relative overflow-hidden w-full bg-muted', aspectRatio, containerClassName)}
    >
      {/* Background Blurred Layer */}
      <img
        src={src}
        alt=""
        aria-hidden="true"
        className="absolute inset-0 w-full h-full object-cover blur-xl scale-110 opacity-50 dark:opacity-30"
      />

      {/* Loading Skeleton */}
      {!isLoaded && <Skeleton className="absolute inset-0 w-full h-full z-20" />}

      {/* Main Image Layer */}
      <img
        src={src}
        alt={alt || ''}
        onLoad={() => setIsLoaded(true)}
        className={cn(
          'relative w-full h-full object-contain z-10 transition-opacity duration-300',
          isLoaded ? 'opacity-100' : 'opacity-0',
          className
        )}
        {...props}
      />
    </div>
  );
}
