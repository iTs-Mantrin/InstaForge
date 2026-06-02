'use client';

import { useEffect, useRef, useCallback } from 'react';
import { motion } from 'framer-motion';
import { Loader2, AlertTriangle } from 'lucide-react';
import type { InstagramPost } from '@/types';
import { useFeedStore } from '@/store/feed';
import { FEED_SCROLL_THRESHOLD } from '@/lib/constants';
import { cn } from '@/lib/utils';
import UserProfileHeader from './UserProfileHeader';
import MediaGrid from './MediaGrid';

interface UserFeedProps {
  username: string;
  onDownloadPost: (post: InstagramPost, mediaId?: string) => void;
  isDownloading?: boolean;
}

export default function UserFeed({ username, onDownloadPost, isDownloading }: UserFeedProps) {
  const {
    user,
    userLoading,
    userError,
    items,
    feedLoading,
    feedError,
    loadingMore,
    hasMore,
    fetchUser,
    fetchFeed,
    loadMore,
  } = useFeedStore();

  const sentinelRef = useRef<HTMLDivElement>(null);

  // Fetch user info + initial feed on mount
  useEffect(() => {
    useFeedStore.getState().reset();
    fetchUser(username);
    fetchFeed(username, true);
  }, [username, fetchUser, fetchFeed]);

  // Infinite scroll via IntersectionObserver
  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const [entry] = entries;
      if (entry.isIntersecting && hasMore && !loadingMore && !feedLoading) {
        loadMore(username);
      }
    },
    [hasMore, loadingMore, feedLoading, loadMore, username]
  );

  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel) return;

    const observer = new IntersectionObserver(handleObserver, {
      rootMargin: `${FEED_SCROLL_THRESHOLD}px`,
    });
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [handleObserver]);

  // Loading state
  if (userLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <Loader2 className="h-8 w-8 text-brand-400 animate-spin mb-3" />
        <p className="text-sm text-surface-400">Loading user profile...</p>
      </div>
    );
  }

  // Error state
  if (userError) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-col items-center justify-center py-20"
      >
        <div className="flex items-start gap-3 p-4 rounded-xl bg-red-500/10 border border-red-500/20 max-w-md">
          <AlertTriangle className="h-5 w-5 text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-medium text-red-300">Failed to load user</p>
            <p className="text-xs text-red-400/80 mt-1">{userError}</p>
          </div>
        </div>
      </motion.div>
    );
  }

  // No user found
  if (!user) return null;

  return (
    <div className="space-y-6">
      {/* Profile Header */}
      <UserProfileHeader user={user} />

      {/* Feed */}
      <div>
        <h3 className="text-sm font-semibold text-surface-100 mb-4">
          Posts ({user.postCount})
        </h3>

        {feedLoading && items.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16">
            <Loader2 className="h-8 w-8 text-brand-400 animate-spin mb-3" />
            <p className="text-sm text-surface-400">Loading posts...</p>
          </div>
        ) : feedError && items.length === 0 ? (
          <div className="flex items-start gap-3 p-4 rounded-xl bg-red-500/10 border border-red-500/20 max-w-md mx-auto">
            <AlertTriangle className="h-5 w-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium text-red-300">Failed to load posts</p>
              <p className="text-xs text-red-400/80 mt-1">{feedError}</p>
            </div>
          </div>
        ) : (
          <>
            <MediaGrid
              items={items}
              onDownloadPost={onDownloadPost}
              isDownloading={isDownloading}
            />

            {/* Sentinel for infinite scroll */}
            <div ref={sentinelRef} className="h-4" />

            {/* Loading more indicator */}
            {loadingMore && (
              <div className="flex items-center justify-center py-6">
                <Loader2 className="h-5 w-5 text-brand-400 animate-spin" />
              </div>
            )}

            {/* No more items */}
            {!hasMore && items.length > 0 && (
              <p className="text-center text-xs text-surface-500 py-6">
                All posts loaded
              </p>
            )}
          </>
        )}
      </div>
    </div>
  );
}
