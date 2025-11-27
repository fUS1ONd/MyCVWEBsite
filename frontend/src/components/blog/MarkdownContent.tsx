import React from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeHighlight from 'rehype-highlight';
import rehypeRaw from 'rehype-raw';
import remarkGfm from 'remark-gfm';
import 'highlight.js/styles/github-dark.css';

interface MarkdownContentProps {
  content: string;
}

export const MarkdownContent: React.FC<MarkdownContentProps> = ({ content }) => {
  return (
    <div className="prose prose-lg max-w-none prose-headings:font-bold prose-a:text-blue-600 prose-a:no-underline hover:prose-a:underline prose-img:rounded-lg prose-pre:bg-gray-900 prose-pre:text-gray-100">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight, rehypeRaw]}
        components={{
          // Custom image rendering
          img: ({ ...props }) => (
            <img
              {...props}
              className="w-full h-auto rounded-lg shadow-md"
              loading="lazy"
              alt={props.alt || ''}
            />
          ),
          // Custom link rendering
          a: ({ ...props }) => {
            const href = props.href || '';
            const isExternal = href.startsWith('http');
            return (
              <a
                {...props}
                target={isExternal ? '_blank' : undefined}
                rel={isExternal ? 'noopener noreferrer' : undefined}
              />
            );
          },
          // Custom video rendering for embedded videos
          iframe: ({ ...props }) => (
            <div className="relative w-full" style={{ paddingBottom: '56.25%' }}>
              <iframe
                {...props}
                className="absolute top-0 left-0 w-full h-full rounded-lg"
                allowFullScreen
              />
            </div>
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
};
