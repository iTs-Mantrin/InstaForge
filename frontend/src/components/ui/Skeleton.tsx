import { cn } from '@/lib/utils';

interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className }: SkeletonProps) {
  return <div className={cn('skeleton', className)} />;
}

export function VideoInfoSkeleton() {
  return (
    <div className="glass rounded-2xl overflow-hidden">
      <div className="grid grid-cols-1 md:grid-cols-5 gap-0">
        <div className="md:col-span-2 aspect-video md:min-h-[240px] skeleton" />
        <div className="md:col-span-3 p-6 space-y-4">
          <Skeleton className="h-6 w-3/4" />
          <Skeleton className="h-4 w-1/2" />
          <div className="flex gap-4">
            <Skeleton className="h-3 w-20" />
            <Skeleton className="h-3 w-20" />
          </div>
        </div>
      </div>
    </div>
  );
}

export function QualitySelectorSkeleton() {
  return (
    <div className="glass rounded-2xl p-6 space-y-6">
      <div className="flex gap-2">
        <Skeleton className="h-10 flex-1 rounded-lg" />
        <Skeleton className="h-10 flex-1 rounded-lg" />
      </div>
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-16 rounded-xl" />
        ))}
      </div>
      <Skeleton className="h-12 w-full rounded-xl" />
    </div>
  );
}

export function MediaInfoSkeleton() {
  return (
    <div className="glass rounded-2xl overflow-hidden">
      <div className="aspect-square sm:aspect-auto sm:min-h-[300px] skeleton" />
      <div className="p-5 space-y-4">
        <div className="flex items-center gap-3">
          <Skeleton className="h-9 w-9 rounded-full" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-28" />
            <Skeleton className="h-3 w-20" />
          </div>
        </div>
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
        <div className="flex gap-4">
          <Skeleton className="h-3 w-16" />
          <Skeleton className="h-3 w-16" />
        </div>
        <div className="flex gap-2 pt-2">
          <Skeleton className="h-10 flex-1 rounded-xl" />
          <Skeleton className="h-10 flex-1 rounded-xl" />
        </div>
      </div>
    </div>
  );
}
