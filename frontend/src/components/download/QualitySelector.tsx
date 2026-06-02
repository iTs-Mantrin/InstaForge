'use client';

import { motion } from 'framer-motion';
import { Download } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { InstagramMediaItem } from '@/types';

interface DownloadControlsProps {
  mediaItems: InstagramMediaItem[];
  onDownloadMedia: (mediaId: string) => void;
  onDownloadAll?: () => void;
  disabled?: boolean;
}

export default function DownloadControls({
  mediaItems,
  onDownloadMedia,
  onDownloadAll,
  disabled,
}: DownloadControlsProps) {
  const isCarousel = mediaItems.length > 1;

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay: 0.1 }}
      className="glass rounded-2xl p-6 space-y-4"
    >
      <h3 className="text-sm font-semibold text-surface-100">Download Options</h3>

      <div className="space-y-2">
        {mediaItems.map((media, index) => (
          <button
            key={media.id}
            onClick={() => onDownloadMedia(media.id)}
            disabled={disabled}
            className={cn(
              'w-full flex items-center justify-between px-4 py-3 rounded-xl text-sm transition-all',
              'bg-surface-800/50 text-surface-300 hover:bg-surface-700/50 hover:text-surface-100',
              'border border-(--border-subtle)',
              disabled && 'opacity-50 cursor-not-allowed'
            )}
          >
            <span>
              {media.type === 'video' ? 'Video' : 'Image'}
              {isCarousel && ` #${index + 1}`}
            </span>
            <Download className="h-4 w-4" />
          </button>
        ))}
      </div>

      {isCarousel && onDownloadAll && (
        <button
          onClick={onDownloadAll}
          disabled={disabled}
          className={cn(
            'w-full flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl text-sm font-semibold transition-all duration-200',
            'bg-gradient-to-r from-brand-600 to-accent-600 text-white',
            'shadow-lg shadow-brand-500/25 hover:shadow-brand-500/40 hover:-translate-y-0.5',
            'active:translate-y-0',
            disabled && 'opacity-50 cursor-not-allowed hover:translate-y-0 shadow-none'
          )}
        >
          <Download className="h-4 w-4" />
          Download All ({mediaItems.length} items)
        </button>
      )}

      {!isCarousel && (
        <button
          onClick={() => onDownloadMedia(mediaItems[0].id)}
          disabled={disabled}
          className={cn(
            'w-full flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl text-sm font-semibold transition-all duration-200',
            'bg-gradient-to-r from-brand-600 to-accent-600 text-white',
            'shadow-lg shadow-brand-500/25 hover:shadow-brand-500/40 hover:-translate-y-0.5',
            'active:translate-y-0',
            disabled && 'opacity-50 cursor-not-allowed hover:translate-y-0 shadow-none'
          )}
        >
          <Download className="h-4 w-4" />
          Download {mediaItems[0]?.type === 'video' ? 'Video' : 'Image'}
        </button>
      )}
    </motion.div>
  );
}
