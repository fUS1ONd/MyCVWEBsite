import { createContext, useContext, useState, ReactNode } from 'react';

interface ImageViewerContextType {
  src: string | null;
  isOpen: boolean;
  openImage: (src: string) => void;
  closeImage: () => void;
}

const ImageViewerContext = createContext<ImageViewerContextType | undefined>(undefined);

export function ImageViewerProvider({ children }: { children: ReactNode }) {
  const [src, setSrc] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);

  const openImage = (imageSrc: string) => {
    setSrc(imageSrc);
    setIsOpen(true);
  };

  const closeImage = () => {
    setIsOpen(false);
    setSrc(null);
  };

  return (
    <ImageViewerContext.Provider value={{ src, isOpen, openImage, closeImage }}>
      {children}
    </ImageViewerContext.Provider>
  );
}

export function useImageViewer() {
  const context = useContext(ImageViewerContext);
  if (context === undefined) {
    throw new Error('useImageViewer must be used within an ImageViewerProvider');
  }
  return context;
}
