'use client';

import { create } from 'zustand';
import type { InstagramPost, InstagramUser, InstagramStory, InstagramFeedResponse } from '@/types';
import { getUserInfo, getUserFeed, getUserStories } from '@/services/api';

interface FeedState {
  // User info
  user: InstagramUser | null;
  userLoading: boolean;
  userError: string | null;

  // Feed items
  items: InstagramPost[];
  feedLoading: boolean;
  feedError: string | null;
  cursor: string | undefined;
  hasMore: boolean;
  loadingMore: boolean;

  // Stories
  stories: InstagramStory[];
  storiesLoading: boolean;

  // Actions
  fetchUser: (username: string) => Promise<void>;
  fetchFeed: (username: string, refresh?: boolean) => Promise<void>;
  loadMore: (username: string) => Promise<void>;
  fetchStories: (username: string) => Promise<void>;
  reset: () => void;
}

export const useFeedStore = create<FeedState>((set, get) => ({
  // User info
  user: null,
  userLoading: false,
  userError: null,

  // Feed items
  items: [],
  feedLoading: false,
  feedError: null,
  cursor: undefined,
  hasMore: true,
  loadingMore: false,

  // Stories
  stories: [],
  storiesLoading: false,

  // Actions
  fetchUser: async (username: string) => {
    set({ userLoading: true, userError: null });
    try {
      const user = await getUserInfo(username);
      set({ user, userLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch user info';
      set({ userError: message, userLoading: false, user: null });
    }
  },

  fetchFeed: async (username: string, refresh = true) => {
    if (refresh) {
      set({ feedLoading: true, feedError: null, items: [], cursor: undefined, hasMore: true });
    }

    try {
      const response: InstagramFeedResponse = await getUserFeed(
        username,
        refresh ? undefined : get().cursor
      );
      set((state) => ({
        items: refresh ? response.items : [...state.items, ...response.items],
        cursor: response.cursor,
        hasMore: response.hasMore,
        feedLoading: false,
        loadingMore: false,
      }));
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch feed';
      set({ feedError: message, feedLoading: false, loadingMore: false });
    }
  },

  loadMore: async (username: string) => {
    const { loadingMore, hasMore, feedLoading } = get();
    if (loadingMore || !hasMore || feedLoading) return;

    set({ loadingMore: true });
    try {
      const response: InstagramFeedResponse = await getUserFeed(username, get().cursor);
      set((state) => ({
        items: [...state.items, ...response.items],
        cursor: response.cursor,
        hasMore: response.hasMore,
        loadingMore: false,
      }));
    } catch {
      set({ loadingMore: false });
    }
  },

  fetchStories: async (username: string) => {
    set({ storiesLoading: true });
    try {
      const response = await getUserStories(username);
      set({ stories: response.stories, storiesLoading: false });
    } catch {
      set({ storiesLoading: false });
    }
  },

  reset: () => {
    set({
      user: null,
      userLoading: false,
      userError: null,
      items: [],
      feedLoading: false,
      feedError: null,
      cursor: undefined,
      hasMore: true,
      loadingMore: false,
      stories: [],
      storiesLoading: false,
    });
  },
}));
