import { Comment } from '@/lib/types';
import { CommentItem } from './CommentItem';

interface CommentThreadProps {
  comment: Comment;
  postSlug: string;
  allComments: Comment[];
  replyingTo: number | null;
  setReplyingTo: (id: number | null) => void;
  depth?: number;
}

export function CommentThread({
  comment,
  postSlug,
  allComments,
  replyingTo,
  setReplyingTo,
  depth = 0,
}: CommentThreadProps) {
  const replies = comment.replies || [];

  const borderColors = [
    'border-primary/30',
    'border-blue-500/30',
    'border-purple-500/30',
    'border-green-500/30',
  ];

  const borderColor = borderColors[depth % borderColors.length];

  return (
    <div className="space-y-4">
      <CommentItem
        comment={comment}
        postSlug={postSlug}
        replyingTo={replyingTo}
        setReplyingTo={setReplyingTo}
        depth={depth}
      />

      {replies.length > 0 && (
        <div className={`ml-3 sm:ml-6 space-y-4 border-l-2 ${borderColor} pl-3 sm:pl-4`}>
          {replies.map((reply) => (
            <CommentThread
              key={reply.id}
              comment={reply}
              postSlug={postSlug}
              allComments={allComments}
              replyingTo={replyingTo}
              setReplyingTo={setReplyingTo}
              depth={depth + 1}
            />
          ))}
        </div>
      )}
    </div>
  );
}
