'use client';

import Image from 'next/image';
import { motion } from 'framer-motion';
import { Heart, MessageCircle, Film, ImageIcon, Download } from 'lucide-react';
import type { InstagramPost } from '@/types';
import { formatCount, cn } from '@/lib/utils';

interface MediaGridProps {
  items: InstagramPost[];
  onDownloadPost: (post: InstagramPost, mediaId?: string) => void;
  isDownloading?: boolean;
}

export default function MediaGrid({ items, onDownloadPost, isDownloading }: MediaGridProps) {
  if (items.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-surface-500">
        <ImageIcon className="h-12 w-12 mb-3" />
        <p className="text-sm">No posts found</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      {items.map((post, index) => (
        <motion.div
          key={post.id}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: (index % 12) * 0.05 }}
          className="group relative aspect-square rounded-xl overflow-hidden bg-surface-900"
        >
          {/* Thumbnail */}
          <Image
            src={post.displayUrl}
            alt={post.caption || 'Instagram post'}
            fill
            className="object-cover transition-transform duration-300 group-hover:scale-105"
            sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 25vw"
          />

          {/* Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200">
            <div className="absolute bottom-3 left-3 right-3 space-y-2">
              {/* Stats */}
              <div className="flex items-center gap-3 text-white/80">
                <span className="flex items-center gap-1 text-xs">
                  <Heart className="h-3 w-3" />
                  {formatCount(post.likesCount)}
                </span>
                <span className="flex items-center gap-1 text-xs">
                  <MessageCircle className="h-3 w-3" />
                  {formatCount(post.commentsCount)}
                </span>
              </div>

              {/* Download button */}
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onDownloadPost(post);
                }}
                disabled={isDownloading}
                className={cn(
                  'flex items-center justify-center gap-1.5 w-full py-1.5 rounded-lg text-xs font-medium transition-all',
                  'bg-white/20 backdrop-blur-sm text-white hover:bg-white/30',
                  isDownloading && 'opacity-50 cursor-not-allowed'
                )}
              >
                <Download className="h-3 w-3" />
                Download
              </button>
            </div>
          </div>

          {/* Type badge */}
          <div className="absolute top-2 right-2">
            {post.mediaType === 'video' || post.mediaItems.some((m) => m.type === 'video') ? (
              <span className="flex items-center gap-1 px-1.5 py-0.5 rounded-md bg-black/60 backdrop-blur-sm text-[10px] text-white">
                <Film className="h-2.5 w-2.5" />
              </span>
            ) : post.mediaType === 'carousel' ? (
              <span className="flex items-center gap-1 px-1.5 py-0.5 rounded-md bg-black/60 backdrop-blur-sm text-[10px] text-white">
                <ImageIcon className="h-2.5 w-2.5" />
              </span>
            ) : null}
          </div>

          {/* Carousel indicator */}
          {post.mediaType === 'carousel' && (
            <div className="absolute top-2 right-2 px-1.5 py-0.5 rounded-md bg-black/60 backdrop-blur-sm text-[10px] text-white">
              {post.mediaItems.length}
            </div>
          )}
        </motion.div>
      ))}
    </div>
  );
}
