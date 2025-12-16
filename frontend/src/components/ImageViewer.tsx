import { useEffect, useState, useRef } from 'react';
import { useImageViewer } from '@/contexts/ImageViewerContext';
import { Button } from '@/components/ui/button';
import { X, Plus, Minus } from 'lucide-react';
import { cn } from '@/lib/utils';

export function ImageViewer() {
  const { src, isOpen, closeImage } = useImageViewer();
  const [scale, setScale] = useState(1);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const dragStart = useRef({ x: 0, y: 0 });
  const lastPosition = useRef({ x: 0, y: 0 });
  const initialPinchDistance = useRef<number | null>(null);
  const lastScale = useRef(1);
  const hasMoved = useRef(false);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      setScale(1);
      setPosition({ x: 0, y: 0 });
      lastPosition.current = { x: 0, y: 0 };
      lastScale.current = 1;
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') closeImage();
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [closeImage]);

  const updateScale = (newScale: number) => {
    const s = Math.min(Math.max(0.5, newScale), 5);
    setScale(s);
    lastScale.current = s;
  };

  const handleZoomIn = () => updateScale(scale + 0.5);
  const handleZoomOut = () => updateScale(scale - 0.5);

  const handleImageClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (isDragging || hasMoved.current) return;
    if (e.ctrlKey) {
      handleZoomOut();
    } else {
      handleZoomIn();
    }
  };

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    hasMoved.current = false;
    dragStart.current = { x: e.clientX, y: e.clientY };
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging) return;
    const dx = e.clientX - dragStart.current.x;
    const dy = e.clientY - dragStart.current.y;

    if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
      hasMoved.current = true;
    }

    setPosition({
      x: lastPosition.current.x + dx,
      y: lastPosition.current.y + dy,
    });
  };

  const handleMouseUp = () => {
    setIsDragging(false);
    lastPosition.current = position;
  };

  // Touch handlers for pinch zoom and pan
  const handleTouchStart = (e: React.TouchEvent) => {
    if (e.touches.length === 2) {
      const dist = Math.hypot(
        e.touches[0].clientX - e.touches[1].clientX,
        e.touches[0].clientY - e.touches[1].clientY
      );
      initialPinchDistance.current = dist;
    } else if (e.touches.length === 1) {
      setIsDragging(true);
      hasMoved.current = false;
      dragStart.current = { x: e.touches[0].clientX, y: e.touches[0].clientY };
    }
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (e.touches.length === 2 && initialPinchDistance.current) {
      const dist = Math.hypot(
        e.touches[0].clientX - e.touches[1].clientX,
        e.touches[0].clientY - e.touches[1].clientY
      );
      const delta = dist - initialPinchDistance.current;
      // Sensitivity factor
      const scaleDiff = delta * 0.01;
      updateScale(lastScale.current + scaleDiff);
    } else if (e.touches.length === 1 && isDragging) {
      const dx = e.touches[0].clientX - dragStart.current.x;
      const dy = e.touches[0].clientY - dragStart.current.y;

      if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
        hasMoved.current = true;
      }

      setPosition({
        x: lastPosition.current.x + dx,
        y: lastPosition.current.y + dy,
      });
    }
  };

  const handleTouchEnd = () => {
    setIsDragging(false);
    lastPosition.current = position;
    lastScale.current = scale;
    initialPinchDistance.current = null;
  };

  if (!isOpen || !src) return null;

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/90 backdrop-blur-sm"
      onClick={closeImage}
    >
      <div className="absolute top-4 right-4 z-50">
        <Button
          variant="ghost"
          size="icon"
          onClick={closeImage}
          className="text-white hover:bg-white/20 rounded-full h-12 w-12"
        >
          <X className="h-8 w-8" />
        </Button>
      </div>

      <div
        className="relative w-full h-full flex items-center justify-center overflow-hidden touch-none"
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
      >
        <img
          src={src}
          alt="Zoomed"
          className={cn(
            'max-w-full max-h-full transition-transform duration-75 ease-linear select-none cursor-move',
            isDragging ? 'cursor-grabbing' : 'cursor-grab'
          )}
          style={{
            transform: `translate(${position.x}px, ${position.y}px) scale(${scale})`,
          }}
          onClick={handleImageClick}
          draggable={false}
        />
      </div>

      <div
        className="absolute bottom-8 left-1/2 -translate-x-1/2 z-50 flex items-center gap-4 bg-black/60 px-6 py-3 rounded-full backdrop-blur-md"
        onClick={(e) => e.stopPropagation()}
      >
        <Button
          variant="ghost"
          size="icon"
          onClick={handleZoomOut}
          className="text-white hover:bg-white/20 rounded-full h-8 w-8"
        >
          <Minus className="h-5 w-5" />
        </Button>

        <span className="text-white font-medium min-w-[3rem] text-center">
          {Math.round(scale * 100)}%
        </span>

        <Button
          variant="ghost"
          size="icon"
          onClick={handleZoomIn}
          className="text-white hover:bg-white/20 rounded-full h-8 w-8"
        >
          <Plus className="h-5 w-5" />
        </Button>
      </div>
    </div>
  );
}
