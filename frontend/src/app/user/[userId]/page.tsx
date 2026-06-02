'use client';

import { useParams, useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import { ArrowLeft, AlertTriangle } from 'lucide-react';
import Link from 'next/link';
import { useDownloadStore } from '@/store/download';
import { useToastStore } from '@/store/toast';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import UserFeed from '@/components/user/UserFeed';
import ProgressModal from '@/components/download/ProgressModal';
import ToastContainer from '@/components/ui/Toast';
import { APP_NAME } from '@/lib/constants';
import { isValidInstagramUsername } from '@/lib/utils';
import type { InstagramPost } from '@/types';

export default function UserPage() {
  const params = useParams();
  const router = useRouter();
  const { downloadMedia, isDownloading, setShowProgressModal } = useDownloadStore();
  const addToast = useToastStore((s) => s.addToast);

  const userId = decodeURIComponent(params.userId as string);
  const isValid = isValidInstagramUsername(userId);

  const handleDownloadPost = (post: InstagramPost, mediaId?: string) => {
    if (isDownloading) return;

    const postUrl = `https://instagram.com/p/${post.shortcode}`;
    downloadMedia(postUrl, mediaId);
    setShowProgressModal(true);
  };

  if (!isValid) {
    return (
      <div className="min-h-screen flex flex-col bg-surface-950">
        <Navbar />
        <main className="flex-1 flex items-center justify-center p-4">
          <div className="text-center max-w-md">
            <AlertTriangle className="h-12 w-12 text-yellow-400 mx-auto mb-4" />
            <h1 className="text-xl font-bold text-surface-100 mb-2">Invalid Username</h1>
            <p className="text-surface-400 text-sm mb-6">
              The Instagram username you entered is not valid. Please check and try again.
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
        <div className="mx-auto max-w-5xl px-3 sm:px-6 lg:px-8">
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

          {/* User Feed */}
          <UserFeed
            username={userId}
            onDownloadPost={handleDownloadPost}
            isDownloading={isDownloading}
          />
        </div>
      </main>
      <Footer />
      <ProgressModal />
      <ToastContainer />
    </div>
  );
}
