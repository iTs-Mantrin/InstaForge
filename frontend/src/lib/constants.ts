import type { PlatformInfo, FAQItem } from '@/types';

// ============================================================
// API Configuration
// ============================================================

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || 'http://localhost:4000';

export const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL || 'http://localhost:4000';

// ============================================================
// Supported Platforms
// ============================================================

export const SUPPORTED_PLATFORMS: PlatformInfo[] = [
  { name: 'Instagram Posts', icon: 'image', description: 'Single & carousel posts' },
  { name: 'Instagram Reels', icon: 'video', description: 'Vertical video reels' },
  { name: 'Instagram Stories', icon: 'story', description: '24-hour stories' },
];

// ============================================================
// App Info
// ============================================================

export const APP_NAME = 'InstaForge';
export const APP_DOMAIN = 'InstaForge.in';
export const APP_TAGLINE = 'Download Instagram content instantly — free, fast, and unlimited.';
export const APP_DESCRIPTION =
  'InstaForge lets you download Instagram posts, reels, and stories in original quality. No limits, no sign-up needed.';

// ============================================================
// Features
// ============================================================

export const FEATURES = [
  {
    title: 'Posts & Carousels',
    description: 'Download single images, videos, or entire carousel posts with one click.',
    icon: 'image',
  },
  {
    title: 'Reels Support',
    description: 'Save Instagram Reels in original quality. Perfect for offline viewing.',
    icon: 'video',
  },
  {
    title: 'Story Saver',
    description: 'Download stories before they disappear. Keep them forever.',
    icon: 'story',
  },
  {
    title: 'User Feed',
    description: 'Browse and download any public user\'s posts, reels, and stories by username.',
    icon: 'users',
  },
  {
    title: 'Original Quality',
    description: 'Download media in the original resolution — no compression or quality loss.',
    icon: 'sparkles',
  },
  {
    title: 'No Limits',
    description: 'Download unlimited content. No daily caps or restrictions.',
    icon: 'infinity',
  },
] as const;

// ============================================================
// FAQ
// ============================================================

export const FAQ_ITEMS: FAQItem[] = [
  {
    question: 'Is this service free?',
    answer:
      'Yes, InstaForge is completely free to use. There are no hidden charges, subscription fees, or usage limits. You can download as much content as you need without paying a single cent.',
  },
  {
    question: 'Do I need an Instagram account?',
    answer:
      'No, you do not need an Instagram account to use InstaForge. Simply paste the public Instagram URL and download. We do not require or request Instagram login credentials at any point.',
  },
  {
    question: 'Do I need to provide my password?',
    answer:
      'Absolutely not. InstaForge never asks for your Instagram password, username, or any login credentials. Our service only processes publicly accessible content through pasted URLs.',
  },
  {
    question: 'Can I download public posts?',
    answer:
      'Yes, you can download any publicly accessible Instagram post by pasting its URL. Simply copy the link from Instagram, paste it into our search bar, and click Search to generate downloadable media files.',
  },
  {
    question: 'Can I download public reels?',
    answer:
      'Yes, InstaForge supports downloading public Instagram Reels. Paste the reel URL and you will be able to download it in its original video quality.',
  },
  {
    question: 'Can I download public stories?',
    answer:
      'Yes, you can download publicly accessible Instagram stories. Paste the story URL and the service will generate a downloadable link for the story media before it expires.',
  },
  {
    question: 'Does the website store downloaded content?',
    answer:
      'No, InstaForge does not permanently store downloaded content. Media files are temporarily processed to generate downloadable links and are not hosted or retained on our servers beyond the immediate processing session.',
  },
  {
    question: 'Does the website access private accounts?',
    answer:
      'No. InstaForge only processes publicly accessible Instagram content. We do not have the ability to access private accounts, private posts, or any content behind Instagram\'s privacy settings.',
  },
  {
    question: 'Is the website affiliated with Instagram?',
    answer:
      'No, InstaForge is an independent service and is not affiliated with, endorsed by, or connected to Instagram, Meta Platforms Inc., or any of their subsidiaries.',
  },
  {
    question: 'Why is a link not working?',
    answer:
      'A link may not work if the content is from a private account, the URL is invalid or malformed, the content has been deleted, or if Instagram has restricted access. Please ensure you are using a valid public URL and try again.',
  },
  {
    question: 'How do I report copyright concerns?',
    answer:
      'If you believe your copyrighted content is being distributed without authorization through our service, please contact us at manish.kumar@modifyly.in with a detailed description and we will review your request promptly.',
  },
  {
    question: 'How do I contact support?',
    answer:
      'You can reach our support team by visiting the Contact page on our website or by emailing manish.kumar@modifyly.in. We aim to respond to all inquiries within 24 hours on business days.',
  },
];

// ============================================================
// Navigation
// ============================================================

export const NAV_LINKS = [
  { label: 'Home', href: '/' },
  { label: 'About', href: '/about' },
  { label: 'History', href: '/history' },
  { label: 'Contact', href: '/contact' },
] as const;

// ============================================================
// Feed Pagination
// ============================================================

export const FEED_PAGE_SIZE = 12;
export const FEED_SCROLL_THRESHOLD = 200;


