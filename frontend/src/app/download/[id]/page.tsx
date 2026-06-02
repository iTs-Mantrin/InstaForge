'use client';

import { useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, AlertTriangle, RefreshCw } from 'lucide-react';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { getInstagramPreview } from '@/services/api';
import { useDownloadStore } from '@/store/download';
import { useToastStore } from '@/store/toast';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import MediaPreview from '@/components/download/VideoMetadata';
import ProgressModal from '@/components/download/ProgressModal';
import ToastContainer from '@/components/ui/Toast';
import { MediaInfoSkeleton } from '@/components/ui/Skeleton';
import { isValidInstagramUrl } from '@/lib/utils';
import { APP_NAME } from '@/lib/constants';

export default function DownloadPage() {
  const params = useParams();
  const router = useRouter();
  const addToast = useToastStore((s) => s.addToast);

  const {
    mediaInfo,
    isDownloading,
    fetchMediaPreview,
    downloadMedia,
    setShowProgressModal,
  } = useDownloadStore();

  // Decode URL from base64 param
  const encodedId = params.id as string;
  let mediaUrl = '';
  try {
    mediaUrl = atob(decodeURIComponent(encodedId));
  } catch {
    mediaUrl = decodeURIComponent(encodedId);
  }

  // Fetch media info on mount
  const { isLoading, isError, error, refetch } = useQuery({
    queryKey: ['mediaInfo', mediaUrl],
    queryFn: () => fetchMediaPreview(mediaUrl).then(() => useDownloadStore.getState().mediaInfo!),
    enabled: !!mediaUrl && isValidInstagramUrl(mediaUrl),
    retry: 2,
    staleTime: 5 * 60 * 1000,
  });

  const handleDownloadMedia = (mediaId?: string) => {
    if (!mediaInfo || isDownloading) return;
    downloadMedia(mediaUrl, mediaId);
    setShowProgressModal(true);
  };

  const handleDownloadAll = () => {
    if (!mediaInfo || isDownloading) return;
    // Download the first media item to trigger the job,
    // the backend should handle downloading all items
    downloadMedia(mediaUrl);
    setShowProgressModal(true);
  };

  // Show error toast on failure
  useEffect(() => {
    if (isError) {
      addToast({
        type: 'error',
        title: 'Failed to fetch media',
        message: error instanceof Error ? error.message : 'Could not retrieve Instagram media information',
      });
    }
  }, [isError, error, addToast]);

  // If not a valid URL, redirect home
  if (!mediaUrl || !isValidInstagramUrl(mediaUrl)) {
    return (
      <div className="min-h-screen flex flex-col bg-surface-950">
        <Navbar />
        <main className="flex-1 flex items-center justify-center p-4">
          <div className="text-center max-w-md">
            <AlertTriangle className="h-12 w-12 text-yellow-400 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-surface-100 mb-2">Invalid URL</h1>
            <p className="text-surface-400 text-sm mb-6">
              The Instagram URL you provided is not valid. Please check and try again.
            </p>
            <Link
              href="/"
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-accent-600 text-white text-sm font-semibold hover:shadow-lg hover:shadow-brand-500/25 transition-all"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Home
            </Link>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-4xl px-3 sm:px-6 lg:px-8">
          {/* Back link */}
          <motion.div
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            className="mb-6"
          >
            <Link
              href="/"
              className="inline-flex items-center gap-1.5 text-sm text-surface-400 hover:text-surface-100 transition-colors"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to {APP_NAME}
            </Link>
          </motion.div>

          {/* Error state */}
          {isError && !isLoading && !mediaInfo && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="glass rounded-2xl p-8 text-center max-w-lg mx-auto"
            >
              <AlertTriangle className="h-10 w-10 text-red-400 mx-auto mb-4" />
              <h2 className="text-lg font-semibold text-surface-100 mb-2">
                Could not load media
              </h2>
              <p className="text-sm text-surface-400 mb-6">
                {error instanceof Error ? error.message : 'An unexpected error occurred.'}
              </p>
              <button
                onClick={() => refetch()}
                className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-accent-600 text-white text-sm font-semibold hover:shadow-lg hover:shadow-brand-500/25 transition-all"
              >
                <RefreshCw className="h-4 w-4" />
                Try Again
              </button>
            </motion.div>
          )}

          {/* Content */}
          {mediaInfo && !isError && (
            <MediaPreview
              post={mediaInfo}
              onDownloadMedia={handleDownloadMedia}
              onDownloadAll={handleDownloadAll}
              isDownloading={isDownloading}
            />
          )}

          {/* Loading skeleton */}
          {isLoading && !mediaInfo && (
            <div className="max-w-2xl mx-auto">
              <MediaInfoSkeleton />
            </div>
          )}
        </div>
      </main>
      <Footer />
      <ProgressModal />
      <ToastContainer />
    </div>
  );
}
