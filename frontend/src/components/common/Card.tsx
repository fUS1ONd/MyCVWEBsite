import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  padding?: 'none' | 'sm' | 'md' | 'lg';
  onClick?: () => void;
}

const paddingClasses = {
  none: '',
  sm: 'p-4',
  md: 'p-6',
  lg: 'p-8',
};

export const Card: React.FC<CardProps> = ({
  children,
  className = '',
  padding = 'md',
  onClick,
}) => {
  const baseClasses = 'bg-white rounded-lg shadow-md border border-gray-200 transition-shadow';
  const hoverClass = onClick ? 'cursor-pointer hover:shadow-lg' : '';
  const paddingClass = paddingClasses[padding];

  return (
    <div className={`${baseClasses} ${paddingClass} ${hoverClass} ${className}`} onClick={onClick}>
      {children}
    </div>
  );
};

export const CardHeader: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <div className="mb-4 pb-4 border-b border-gray-200">{children}</div>;
};

export const CardTitle: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <h3 className="text-lg font-semibold text-gray-900">{children}</h3>;
};

export const CardContent: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <div className="text-gray-700">{children}</div>;
};

export const CardFooter: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <div className="mt-4 pt-4 border-t border-gray-200">{children}</div>;
};
