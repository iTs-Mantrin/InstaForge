import axios from 'axios';
import type {
  ApiResponse,
  InstagramPost,
  InstagramUser,
  InstagramFeedResponse,
  StoriesResponse,
  DownloadJob,
  ProgressData,
  DownloadResult,
} from '@/types';
import { API_BASE_URL } from '@/lib/constants';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 120000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// ============================================================
// Instagram Preview (post/reel by URL)
// ============================================================

export async function getInstagramPreview(url: string): Promise<InstagramPost> {
  const { data } = await apiClient.post<ApiResponse<InstagramPost>>(
    '/instagram/preview',
    { url }
  );
  if (!data.success) {
    throw new Error(data.error?.message || data.message || 'Failed to fetch Instagram media info');
  }
  return data.data;
}

// ============================================================
// User Info
// ============================================================

export async function getUserInfo(username: string): Promise<InstagramUser> {
  const { data } = await apiClient.get<ApiResponse<InstagramUser>>(
    `/instagram/user/${encodeURIComponent(username)}/info`
  );
  if (!data.success) {
    throw new Error(data.error?.message || data.message || 'Failed to fetch user info');
  }
  return data.data;
}

// ============================================================
// User Feed (with cursor-based pagination)
// ============================================================

export async function getUserFeed(
  username: string,
  cursor?: string,
  limit: number = 12
): Promise<InstagramFeedResponse> {
  const params: Record<string, string | number> = { limit };
  if (cursor) params.cursor = cursor;
  const { data } = await apiClient.get<ApiResponse<InstagramFeedResponse>>(
    `/instagram/user/${encodeURIComponent(username)}/feed`,
    { params }
  );
  if (!data.success) {
    throw new Error(data.error?.message || data.message || 'Failed to fetch user feed');
  }
  return data.data;
}

// ============================================================
// User Stories
// ============================================================

export async function getUserStories(username: string): Promise<StoriesResponse> {
  const { data } = await apiClient.get<ApiResponse<StoriesResponse>>(
    `/instagram/user/${encodeURIComponent(username)}/stories`
  );
  if (!data.success) {
    throw new Error(data.error?.message || data.message || 'Failed to fetch user stories');
  }
  return data.data;
}

// ============================================================
// Download
// ============================================================

export async function startDownload(
  url: string,
  mediaId?: string
): Promise<DownloadJob> {
  const { data } = await apiClient.post<ApiResponse<DownloadJob>>(
    '/instagram/download',
    { url, mediaId }
  );
  if (!data.success) {
    throw new Error(data.error?.message || data.message || 'Failed to start download');
  }
  return data.data;
}

// ============================================================
// Progress
// ============================================================

export async function getProgress(taskId: string): Promise<ProgressData> {
  const { data } = await apiClient.get<ApiResponse<ProgressData>>(
    `/instagram/progress/${taskId}`
  );
  if (!data.success) {
    throw new Error(data.error?.message || 'Failed to fetch progress');
  }
  return data.data;
}

// ============================================================
// Download Result (file URL)
// ============================================================

export async function getDownloadResult(
  taskId: string
): Promise<DownloadResult> {
  const { data } = await apiClient.get<ApiResponse<DownloadResult>>(
    `/instagram/file/${taskId}`
  );
  if (!data.success) {
    throw new Error(data.error?.message || 'Failed to fetch download result');
  }
  return data.data;
}

// ============================================================
// Cancel Download
// ============================================================

export async function cancelDownload(taskId: string): Promise<void> {
  const { data } = await apiClient.delete<ApiResponse<{ status: string }>>(
    `/instagram/${taskId}`
  );
  if (!data.success) {
    throw new Error(data.error?.message || 'Failed to cancel download');
  }
}

export default apiClient;
