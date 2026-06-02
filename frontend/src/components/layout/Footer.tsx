import Link from 'next/link';
import { Download, Globe, MessageSquare, Rss } from 'lucide-react';
import { APP_NAME, APP_DOMAIN } from '@/lib/constants';

const linkGroups = [
  {
    title: 'Quick Links',
    links: [
      { label: 'Home', href: '/' },
      { label: 'History', href: '/history' },
      { label: 'About Us', href: '/about' },
      { label: 'FAQ', href: '/#faq' },
    ],
  },
  {
    title: 'Support',
    links: [
      { label: 'Contact', href: '/contact' },
      { label: 'FAQ', href: '/#faq' },
    ],
  },
  {
    title: 'Policies',
    links: [
      { label: 'Terms of Service', href: '/terms' },
      { label: 'Privacy Policy', href: '/privacy' },
      { label: 'Cookie Policy', href: '/cookies' },
      { label: 'Copyright Policy', href: '/copyright' },
      { label: 'DMCA Policy', href: '/dmca' },
    ],
  },
  {
    title: 'Legal',
    links: [
      { label: 'Disclaimer', href: '/disclaimer' },
      { label: 'Legal Notice', href: '/legal-notice' },
      { label: 'Acceptable Use', href: '/acceptable-use' },
      { label: 'Content Removal', href: '/content-removal' },
    ],
  },
];

export default function Footer() {
  const year = new Date().getFullYear();

  return (
    <footer className="border-t border-(--border-subtle) bg-surface-950">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-16 lg:py-20">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-5">
          {/* Brand */}
          <div className="md:col-span-1">
            <Link href="/" className="inline-flex items-center gap-2.5 mb-5">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-brand-500 to-accent-600 shadow-sm shadow-brand-500/20">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="h-4 w-4 text-white">
                <rect width="18" height="18" x="3" y="3" rx="4" />
                <circle cx="12" cy="12" r="5" />
                <circle cx="17" cy="7" r="1.5" fill="currentColor" />
              </svg>
              </div>
              <span className="text-lg font-bold text-surface-100">{APP_NAME}</span>
            </Link>
            <p className="text-sm text-surface-400 leading-relaxed max-w-sm">
              Download Instagram posts, reels, and stories instantly. No limits, no sign-up, completely free.
            </p>
            <div className="flex items-center gap-3 mt-5">
              <a
                href={`https://${APP_DOMAIN}`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-8 w-8 items-center justify-center rounded-lg border border-(--border-subtle) bg-surface-900 text-surface-500 hover:border-(--border-active) hover:text-surface-100 hover:bg-surface-800 transition-all"
                aria-label="Website"
              >
                <Globe className="h-3.5 w-3.5" />
              </a>
              <a
                href="mailto:manish.kumar@modifyly.in"
                className="flex h-8 w-8 items-center justify-center rounded-lg border border-(--border-subtle) bg-surface-900 text-surface-500 hover:border-(--border-active) hover:text-surface-100 hover:bg-surface-800 transition-all"
                aria-label="Email"
              >
                <MessageSquare className="h-3.5 w-3.5" />
              </a>
              <a
                href="#"
                className="flex h-8 w-8 items-center justify-center rounded-lg border border-(--border-subtle) bg-surface-900 text-surface-500 hover:border-(--border-active) hover:text-surface-100 hover:bg-surface-800 transition-all"
                aria-label="RSS"
              >
                <Rss className="h-3.5 w-3.5" />
              </a>
            </div>
          </div>

          {/* Link Groups */}
          {linkGroups.map((group) => (
            <div key={group.title}>
              <h3 className="text-xs font-semibold tracking-wider text-surface-500 uppercase mb-4">
                {group.title}
              </h3>
              <ul className="space-y-3">
                {group.links.map((link) => (
                  <li key={link.label}>
                    <Link
                      href={link.href}
                      className="text-sm text-surface-400 hover:text-surface-100 transition-colors"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom */}
        <div className="mt-14 pt-8 border-t border-(--border-subtle)">
          <p className="text-xs text-surface-500 text-center">
            &copy; {year} {APP_NAME}. All rights reserved. — Not affiliated with Instagram or Meta.
          </p>
        </div>
      </div>
    </footer>
  );
}
