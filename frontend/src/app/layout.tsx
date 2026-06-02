import type { Metadata } from "next";
import { Sora } from "next/font/google";
import "./globals.css";
import Providers from "./providers";

const sora = Sora({
  variable: "--font-sora",
  subsets: ["latin"],
  display: "swap",
});

export const viewport = {
  width: 'device-width',
  initialScale: 1,
};

export const metadata: Metadata = {
  title: {
    default: "InstaForge — Download Instagram Content Instantly",
    template: "%s | InstaForge",
  },
  description:
    "Download Instagram posts, reels, and stories for free — original quality. No limits, no sign-up.",
  keywords: [
    "instagram downloader",
    "download instagram posts",
    "instagram reel downloader",
    "instagram story saver",
    "free instagram downloader",
    "InstaForge",
  ],
  icons: {
    icon: "/icon.svg",
  },
  openGraph: {
    title: "InstaForge — Download Instagram Content Instantly",
    description:
      "Download Instagram posts, reels, and stories for free — original quality. No limits, no sign-up.",
    url: "https://InstaForge.in",
    siteName: "InstaForge",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "InstaForge — Download Instagram Content Instantly",
    description:
      "Download Instagram posts, reels, and stories for free — original quality. No limits, no sign-up.",
  },
  robots: {
    index: true,
    follow: true,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${sora.variable} h-full`} suppressHydrationWarning>
      <body className="min-h-full flex flex-col bg-surface-950 text-surface-100 antialiased">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
