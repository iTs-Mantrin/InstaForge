// ============================================================
// Instagram Media Types
// ============================================================

export type InstagramMediaType = 'image' | 'video' | 'carousel';

export interface InstagramMediaItem {
  id: string;
  type: 'image' | 'video';
  url: string;
  thumbnail?: string;
  width?: number;
  height?: number;
  hasAudio?: boolean;
  downloadUrl?: string;
}

export interface InstagramPost {
  id: string;
  shortcode: string;
  mediaType: InstagramMediaType;
  mediaItems: InstagramMediaItem[];
  caption?: string;
  likesCount: number;
  commentsCount: number;
  timestamp: number;
  ownerUsername: string;
  ownerFullName?: string;
  ownerProfilePic?: string;
  displayUrl: string;
  isSponsored?: boolean;
}

export interface InstagramUser {
  id: string;
  username: string;
  fullName: string;
  profilePicUrl: string;
  followerCount: number;
  followingCount: number;
  postCount: number;
  isVerified: boolean;
  biography?: string;
  externalUrl?: string;
  isPrivate?: boolean;
}

export interface InstagramStory {
  id: string;
  type: 'image' | 'video';
  url: string;
  thumbnail?: string;
  timestamp: number;
  expiresAt: number;
  ownerUsername: string;
  ownerProfilePic?: string;
  viewed: boolean;
}

export interface InstagramFeedResponse {
  items: InstagramPost[];
  hasMore: boolean;
  cursor?: string;
}

export interface StoriesResponse {
  stories: InstagramStory[];
  owner: {
    username: string;
    profilePicUrl: string;
  };
}

// ============================================================
// API Response Types
// ============================================================

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  timestamp: string;
  message?: string;
  error?: {
    code: string;
    message: string;
    details?: string;
  };
}

export interface DownloadJob {
  taskId: string;
  status: DownloadStatus;
  source: string;
  message: string;
}

export type DownloadStatus =
  | 'queued'
  | 'downloading'
  | 'processing'
  | 'uploading'
  | 'done'
  | 'error'
  | 'failed'
  | 'cancelled';

export interface ProgressData {
  taskId: string;
  percent: number;
  speed: string;
  eta: string;
  filename: string;
  status: DownloadStatus;
  errorMsg: string;
  downloadUrl: string;
}

export interface DownloadResult {
  downloadUrl: string;
}

// ============================================================
// Frontend State Types
// ============================================================

export type DownloadType = 'media';

export interface HistoryItem {
  taskId: string;
  title: string;
  thumbnail: string;
  url: string;
  mediaType: InstagramMediaType;
  status: DownloadStatus;
  downloadUrl: string;
  timestamp: number;
  ownerUsername: string;
}

// ============================================================
// Form Types
// ============================================================

export interface DownloadFormData {
  url: string;
  mediaId?: string;
}

export interface SearchMode {
  mode: 'link' | 'user';
}

// ============================================================
// Toast Types
// ============================================================

export interface Toast {
  id: string;
  type: 'success' | 'error' | 'info' | 'warning';
  title: string;
  message?: string;
  duration?: number;
}

// ============================================================
// Supported Platform
// ============================================================

export interface PlatformInfo {
  name: string;
  icon: string;
  description: string;
}

// ============================================================
// FAQ
// ============================================================

export interface FAQItem {
  question: string;
  answer: string;
}
