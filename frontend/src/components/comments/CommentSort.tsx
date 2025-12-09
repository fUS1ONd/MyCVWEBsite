import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

export type SortOption = 'newest' | 'oldest' | 'most_liked';

interface CommentSortProps {
  value: SortOption;
  onChange: (value: SortOption) => void;
}

export function CommentSort({ value, onChange }: CommentSortProps) {
  return (
    <Select value={value} onValueChange={onChange}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Sort comments" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="newest">Newest first</SelectItem>
        <SelectItem value="oldest">Oldest first</SelectItem>
        <SelectItem value="most_liked">Most liked</SelectItem>
      </SelectContent>
    </Select>
  );
}
