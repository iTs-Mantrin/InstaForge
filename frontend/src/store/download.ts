'use client';

import { create } from 'zustand';
import type {
  InstagramPost,
  DownloadJob,
  ProgressData,
  DownloadStatus,
} from '@/types';
import { getInstagramPreview, startDownload as startDownloadApi, getProgress, getDownloadResult, cancelDownload as cancelDownloadApi } from '@/services/api';
import { subscribeToProgress } from '@/services/websocket';

interface DownloadState {
  // Input
  url: string;
  setUrl: (url: string) => void;

  // Media info
  mediaInfo: InstagramPost | null;
  mediaInfoLoading: boolean;
  mediaInfoError: string | null;

  // Download job
  activeJob: DownloadJob | null;
  progress: ProgressData | null;
  isDownloading: boolean;
  downloadError: string | null;

  // Modal
  showProgressModal: boolean;
  setShowProgressModal: (show: boolean) => void;

  // Actions
  fetchMediaPreview: (url: string) => Promise<void>;
  downloadMedia: (url: string, mediaId?: string) => Promise<void>;
  pollProgress: () => Promise<void>;
  cancelDownload: () => void;
  reset: () => void;
}

export const useDownloadStore = create<DownloadState>((set, get) => ({
  // Input
  url: '',
  setUrl: (url) => set({ url }),

  // Media info
  mediaInfo: null,
  mediaInfoLoading: false,
  mediaInfoError: null,

  // Download job
  activeJob: null,
  progress: null,
  isDownloading: false,
  downloadError: null,

  // Modal
  showProgressModal: false,
  setShowProgressModal: (show) => set({ showProgressModal: show }),

  // Actions
  fetchMediaPreview: async (url: string) => {
    set({ mediaInfoLoading: true, mediaInfoError: null, url });
    try {
      const info = await getInstagramPreview(url);
      set({ mediaInfo: info, mediaInfoLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch media info';
      set({ mediaInfoError: message, mediaInfoLoading: false, mediaInfo: null });
    }
  },

  downloadMedia: async (url: string, mediaId?: string) => {
    set({ isDownloading: true, downloadError: null, showProgressModal: true });

    try {
      const job = await startDownloadApi(url, mediaId);
      set({ activeJob: job });

      // Subscribe to WebSocket progress (async, may be a no-op if backend unreachable)
      const unsubscribe = await subscribeToProgress(
        job.taskId,
        (wsProgress) => {
          set((state) => ({
            progress: {
              ...(state.progress || {
                taskId: job.taskId,
                percent: 0,
                speed: '',
                eta: '',
                filename: '',
                status: 'queued',
                errorMsg: '',
                downloadUrl: '',
              }),
              ...wsProgress,
              status: wsProgress.status as DownloadStatus,
              taskId: job.taskId,
            },
          }));
        },
        (newStatus) => {
          set((state) => ({
            progress: state.progress
              ? { ...state.progress, status: newStatus as DownloadStatus }
              : null,
          }));
        }
      );

      // Fallback: poll progress as well
      const pollInterval = setInterval(async () => {
        try {
          const progress = await getProgress(job.taskId);
          set({ progress });

          if (['done', 'error', 'failed', 'cancelled'].includes(progress.status)) {
            clearInterval(pollInterval);

            if (progress.status === 'done' && !progress.downloadUrl) {
              try {
                const result = await getDownloadResult(job.taskId);
                set((state) => ({
                  progress: state.progress
                    ? { ...state.progress, downloadUrl: result.downloadUrl }
                    : null,
                }));
              } catch {
                // URL fetch failed
              }
            }

            set({ isDownloading: false });
            unsubscribe();
          }
        } catch {
          clearInterval(pollInterval);
        }
      }, 2000);

      setTimeout(() => clearInterval(pollInterval), 3600000);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to start download';
      set({ downloadError: message, isDownloading: false });
    }
  },

  pollProgress: async () => {
    const { activeJob } = get();
    if (!activeJob) return;
    try {
      const progress = await getProgress(activeJob.taskId);
      set({ progress });
    } catch {
      // silently fail
    }
  },

  cancelDownload: async () => {
    const { activeJob } = get();
    if (activeJob?.taskId) {
      try {
        await cancelDownloadApi(activeJob.taskId);
      } catch {
        // API call is best-effort
      }
    }
    set({
      activeJob: null,
      progress: null,
      isDownloading: false,
      showProgressModal: false,
    });
  },

  reset: () => {
    set({
      mediaInfo: null,
      mediaInfoLoading: false,
      mediaInfoError: null,
      activeJob: null,
      progress: null,
      isDownloading: false,
      downloadError: null,
      showProgressModal: false,
    });
  },
}));
