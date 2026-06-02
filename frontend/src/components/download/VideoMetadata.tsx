'use client';

import Image from 'next/image';
import { useState } from 'react';
import { motion } from 'framer-motion';
import { Heart, MessageCircle, User, ChevronLeft, ChevronRight, Image as ImageIcon, Video, Download, DownloadCloud } from 'lucide-react';
import type { InstagramPost, InstagramMediaItem } from '@/types';
import { formatCount, cn } from '@/lib/utils';

interface MediaPreviewProps {
  post: InstagramPost;
  onDownloadMedia?: (mediaId: string) => void;
  onDownloadAll?: () => void;
  isDownloading?: boolean;
}

export default function MediaPreview({ post, onDownloadMedia, onDownloadAll, isDownloading }: MediaPreviewProps) {
  const [carouselIndex, setCarouselIndex] = useState(0);
  const [imgError, setImgError] = useState<Record<string, boolean>>({});

  const isCarousel = post.mediaType === 'carousel';
  const currentMedia = isCarousel ? post.mediaItems[carouselIndex] : post.mediaItems[0];
  const isVideo = currentMedia?.type === 'video';

  const goNext = () => {
    if (carouselIndex < post.mediaItems.length - 1) setCarouselIndex((i) => i + 1);
  };
  const goPrev = () => {
    if (carouselIndex > 0) setCarouselIndex((i) => i - 1);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="glass rounded-2xl overflow-hidden"
    >
      {/* Media Display */}
      <div className="relative bg-surface-900 flex items-center justify-center min-h-[300px] max-h-[500px] overflow-hidden">
        {isVideo ? (
          <video
            src={currentMedia?.url}
            className="w-full h-full object-contain max-h-[500px]"
            controls
            poster={currentMedia?.thumbnail}
          />
        ) : (
          <div className="relative w-full h-full flex items-center justify-center">
            {!imgError[currentMedia?.id || ''] ? (
              <Image
                src={currentMedia?.url || post.displayUrl}
                alt={post.caption || 'Instagram media'}
                width={currentMedia?.width || 600}
                height={currentMedia?.height || 600}
                className="object-contain max-h-[500px] w-auto"
                sizes="(max-width: 768px) 100vw, 50vw"
                priority
                onError={() => {
                  if (currentMedia?.id) setImgError((prev) => ({ ...prev, [currentMedia.id]: true }));
                }}
              />
            ) : (
              <div className="flex flex-col items-center gap-2 text-surface-500">
                <ImageIcon className="h-12 w-12" />
                <span className="text-sm">Image unavailable</span>
              </div>
            )}
          </div>
        )}

        {/* Media type badge */}
        <div className="absolute top-3 left-3">
          {isVideo ? (
            <span className="flex items-center gap-1 px-2 py-1 rounded-lg bg-black/70 backdrop-blur-sm text-[10px] text-white">
              <Video className="h-3 w-3" /> Video
            </span>
          ) : (
            <span className="flex items-center gap-1 px-2 py-1 rounded-lg bg-black/70 backdrop-blur-sm text-[10px] text-white">
              <ImageIcon className="h-3 w-3" /> Photo
            </span>
          )}
        </div>

        {/* Carousel Navigation */}
        {isCarousel && post.mediaItems.length > 1 && (
          <>
            {carouselIndex > 0 && (
              <button
                onClick={goPrev}
                className="absolute left-2 top-1/2 -translate-y-1/2 p-1.5 rounded-full bg-black/60 backdrop-blur-sm text-white hover:bg-black/80 transition-colors"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
            )}
            {carouselIndex < post.mediaItems.length - 1 && (
              <button
                onClick={goNext}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-full bg-black/60 backdrop-blur-sm text-white hover:bg-black/80 transition-colors"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            )}
            {/* Dots */}
            <div className="absolute bottom-3 left-1/2 -translate-x-1/2 flex items-center gap-1.5">
              {post.mediaItems.map((_, i) => (
                <button
                  key={i}
                  onClick={() => setCarouselIndex(i)}
                  className={cn(
                    'w-1.5 h-1.5 rounded-full transition-all',
                    i === carouselIndex ? 'bg-white w-3' : 'bg-white/50'
                  )}
                />
              ))}
            </div>
          </>
        )}

        {/* Carousel badge */}
        {isCarousel && (
          <div className="absolute top-3 right-3 px-2 py-1 rounded-lg bg-black/70 backdrop-blur-sm text-[10px] text-white">
            {carouselIndex + 1}/{post.mediaItems.length}
          </div>
        )}
      </div>

      {/* Details */}
      <div className="p-5 space-y-4">
        {/* Owner */}
        <div className="flex items-center gap-3">
          {post.ownerProfilePic ? (
            <Image
              src={post.ownerProfilePic}
              alt={post.ownerUsername}
              width={36}
              height={36}
              className="rounded-full object-cover"
            />
          ) : (
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-surface-700">
              <User className="h-4 w-4 text-surface-400" />
            </div>
          )}
          <div>
            <p className="text-sm font-semibold text-surface-100">@{post.ownerUsername}</p>
            {post.ownerFullName && (
              <p className="text-xs text-surface-400">{post.ownerFullName}</p>
            )}
          </div>
        </div>

        {/* Caption */}
        {post.caption && (
          <p className="text-sm text-surface-300 leading-relaxed line-clamp-3">
            {post.caption}
          </p>
        )}

        {/* Stats */}
        <div className="flex items-center gap-4 text-xs text-surface-500">
          <span className="flex items-center gap-1">
            <Heart className="h-3.5 w-3.5" />
            {formatCount(post.likesCount)}
          </span>
          <span className="flex items-center gap-1">
            <MessageCircle className="h-3.5 w-3.5" />
            {formatCount(post.commentsCount)}
          </span>
          <span className="text-surface-600">
            {post.mediaItems.length} {post.mediaItems.length === 1 ? 'item' : 'items'}
          </span>
        </div>

        {/* Download Buttons */}
        <div className="space-y-2 pt-2 border-t border-(--border-subtle)">
          <p className="text-xs font-medium text-surface-400 uppercase tracking-wider">Download</p>

          {/* Individual media download */}
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
            {post.mediaItems.map((media, index) => (
              <button
                key={media.id}
                onClick={() => onDownloadMedia?.(media.id)}
                disabled={isDownloading}
                className={cn(
                  'flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-medium transition-all',
                  'bg-surface-800/50 text-surface-300 hover:bg-surface-700/50 hover:text-surface-100',
                  'border border-(--border-subtle)',
                  isDownloading && 'opacity-50 cursor-not-allowed'
                )}
              >
                {media.type === 'video' ? (
                  <Video className="h-3.5 w-3.5 flex-shrink-0" />
                ) : (
                  <ImageIcon className="h-3.5 w-3.5 flex-shrink-0" />
                )}
                <span className="truncate">
                  {media.type === 'video' ? 'Video' : 'Image'} {post.mediaItems.length > 1 ? `#${index + 1}` : ''}
                </span>
                <Download className="h-3 w-3 ml-auto flex-shrink-0" />
              </button>
            ))}
          </div>

          {/* Download All (for carousel) */}
          {isCarousel && post.mediaItems.length > 1 && (
            <button
              onClick={onDownloadAll}
              disabled={isDownloading}
              className={cn(
                'w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-sm font-semibold transition-all',
                'bg-gradient-to-r from-brand-600 to-accent-600 text-white',
                'shadow-lg shadow-brand-500/25 hover:shadow-brand-500/40 hover:-translate-y-0.5',
                'active:translate-y-0',
                isDownloading && 'opacity-50 cursor-not-allowed hover:translate-y-0 shadow-none'
              )}
            >
              <DownloadCloud className="h-4 w-4" />
              Download All ({post.mediaItems.length} items)
            </button>
          )}

          {/* Single download button for non-carousel */}
          {!isCarousel && post.mediaItems.length === 1 && (
            <button
              onClick={() => onDownloadMedia?.(post.mediaItems[0].id)}
              disabled={isDownloading}
              className={cn(
                'w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-sm font-semibold transition-all',
                'bg-gradient-to-r from-brand-600 to-accent-600 text-white',
                'shadow-lg shadow-brand-500/25 hover:shadow-brand-500/40 hover:-translate-y-0.5',
                'active:translate-y-0',
                isDownloading && 'opacity-50 cursor-not-allowed hover:translate-y-0 shadow-none'
              )}
            >
              <Download className="h-4 w-4" />
              Download {mediaItemLabel(post.mediaItems[0])}
            </button>
          )}
        </div>
      </div>
    </motion.div>
  );
}

function mediaItemLabel(media: InstagramMediaItem): string {
  if (media.type === 'video') return 'Video';
  return 'Image';
}
