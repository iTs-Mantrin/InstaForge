'use client';

import { useRouter } from 'next/navigation';
import { useState, useCallback } from 'react';
import { motion } from 'framer-motion';
import { ArrowDown, Sparkles, Shield, Zap, Link as LinkIcon, Search as SearchIcon } from 'lucide-react';
import { APP_NAME, APP_TAGLINE, APP_DESCRIPTION } from '@/lib/constants';
import { useDownloadStore } from '@/store/download';
import { cn } from '@/lib/utils';
import SearchBar from './SearchBar';

export default function Hero() {
  const router = useRouter();
  const fetchMediaPreview = useDownloadStore((s) => s.fetchMediaPreview);
  const [loading, setLoading] = useState(false);
  const [mode, setMode] = useState<'link' | 'user'>('link');

  const handleSearch = useCallback(
    async (query: string) => {
      setLoading(true);
      try {
        if (mode === 'link') {
          // Navigate to download page with the URL
          await fetchMediaPreview(query);
          router.push(`/download/${encodeURIComponent(btoa(query))}`);
        } else {
          // Navigate to user feed page
          router.push(`/user/${encodeURIComponent(query)}`);
        }
      } catch {
        if (mode === 'link') {
          router.push(`/download/${encodeURIComponent(btoa(query))}`);
        } else {
          router.push(`/user/${encodeURIComponent(query)}`);
        }
      } finally {
        setLoading(false);
      }
    },
    [mode, fetchMediaPreview, router]
  );

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.1 },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.5 },
    },
  };

  return (
    <section className="relative min-h-[90vh] flex items-center justify-center overflow-hidden">
      {/* Animated Background */}
      <div className="absolute inset-0 -z-10 overflow-hidden">
        {/* Gradient orbs — responsive sizes */}
        <div className="absolute top-1/4 left-1/4 w-48 h-48 sm:w-72 sm:h-72 lg:w-96 lg:h-96 bg-brand-500/20 rounded-full blur-3xl animate-float" />
        <div className="absolute bottom-1/4 right-1/4 w-40 h-40 sm:w-60 sm:h-60 lg:w-80 lg:h-80 bg-accent-500/15 rounded-full blur-3xl animate-float" style={{ animationDelay: '-3s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 sm:w-[400px] sm:h-[400px] lg:w-[600px] lg:h-[600px] bg-brand-600/5 rounded-full blur-3xl" />

        {/* Grid overlay */}
        <div
          className="absolute inset-0 opacity-[0.015]"
          style={{
            backgroundImage:
              'linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)',
            backgroundSize: '60px 60px',
          }}
        />
      </div>

      <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-20 text-center">
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="space-y-8"
        >
          {/* Badge */}
          <motion.div variants={itemVariants} className="flex justify-center">
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full glass text-sm text-surface-300">
              <Sparkles className="h-3.5 w-3.5 text-brand-400" />
              <span>Free & Unlimited Downloads</span>
            </div>
          </motion.div>

          {/* Heading */}
          <motion.h1
            variants={itemVariants}
            className="text-[clamp(1.75rem,6vw,4.5rem)] sm:text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight leading-[1.1]"
          >
            <span className="text-surface-100">Download <span className="text-[#E4405F]">Instagram</span></span>
            <br />
            <span className="gradient-text">Content Instantly</span>
          </motion.h1>

          {/* Description */}
          <motion.p
            variants={itemVariants}
            className="mx-auto max-w-2xl text-base sm:text-lg text-surface-400 leading-relaxed"
          >
            {APP_DESCRIPTION}
          </motion.p>

          {/* Mode Toggle */}
          <motion.div variants={itemVariants} className="flex justify-center">
            <div className="inline-flex rounded-xl bg-surface-800/50 p-1 gap-1">
              <button
                onClick={() => setMode('link')}
                className={cn(
                  'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                  mode === 'link'
                    ? 'bg-gradient-to-r from-brand-600 to-accent-600 text-white shadow-lg shadow-brand-500/20'
                    : 'text-surface-400 hover:text-surface-100'
                )}
              >
                <LinkIcon className="h-4 w-4" />
                Link
              </button>
              <button
                onClick={() => setMode('user')}
                className={cn(
                  'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                  mode === 'user'
                    ? 'bg-gradient-to-r from-brand-600 to-accent-600 text-white shadow-lg shadow-brand-500/20'
                    : 'text-surface-400 hover:text-surface-100'
                )}
              >
                <SearchIcon className="h-4 w-4" />
                User
              </button>
            </div>
          </motion.div>

          {/* Search Bar */}
          <motion.div variants={itemVariants} className="pt-2">
            <SearchBar
              mode={mode}
              onSearch={handleSearch}
              isLoading={loading}
              autoFocus
            />
          </motion.div>

          {/* Trust badges */}
          <motion.div
            variants={itemVariants}
            className="flex flex-wrap items-center justify-center gap-x-6 gap-y-2 pt-8 text-xs text-surface-500"
          >
            <span className="flex items-center gap-1.5 whitespace-nowrap">
              <Shield className="h-3.5 w-3.5 flex-shrink-0" /> No sign-up required
            </span>
            <span className="flex items-center gap-1.5 whitespace-nowrap">
              <Zap className="h-3.5 w-3.5 flex-shrink-0" /> Unlimited downloads
            </span>
            <span className="flex items-center gap-1.5 whitespace-nowrap">
              <span className="text-surface-600">Posts · Reels · Stories</span>
            </span>
          </motion.div>

          {/* Scroll indicator */}
          <motion.div
            variants={itemVariants}
            className="pt-12 flex justify-center"
            animate={{ y: [0, 8, 0] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            <ArrowDown className="h-5 w-5 text-surface-500" />
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
}
