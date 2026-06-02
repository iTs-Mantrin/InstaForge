import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME, APP_DOMAIN } from '@/lib/constants';

export const metadata = {
  title: 'Legal Notice',
  description: `Legal Notice for ${APP_NAME} — legal entity information and contact details.`,
};

const sections = [
  {
    title: '1. Legal Entity Information',
    content:
      `${APP_NAME} is operated and maintained by an individual developer. The service is provided on an "as is" basis for personal and non-commercial use. If you need to contact us for legal purposes, please use the contact information provided below.`,
  },
  {
    title: '2. Service Address',
    content:
      'The primary service address associated with this website is: India. All legal correspondence should be directed to our email address in the first instance.',
  },
  {
    title: '3. Contact Information',
    content:
      'For all legal inquiries, notices, and correspondence, please contact us at: manish.kumar@modifyly.in. We will acknowledge receipt of legal notices within a reasonable timeframe and respond as appropriate.',
  },
  {
    title: '4. Governing Law',
    content:
      'These terms and conditions, and any disputes arising out of or in connection with them, shall be governed by and construed in accordance with the laws of India. You agree to submit to the exclusive jurisdiction of the courts located in India for the resolution of any disputes.',
  },
  {
    title: '5. Dispute Resolution',
    content:
      'Any dispute arising from or relating to the use of this service shall first be attempted to be resolved through informal negotiation. If the dispute cannot be resolved within 30 days, either party may seek legal remedies in accordance with applicable laws. Both parties agree to attempt to resolve disputes amicably before initiating formal legal proceedings.',
  },
  {
    title: '6. Contact for Legal Inquiries',
    content:
      'For any legal questions, requests, or formal notices, please reach out via email: manish.kumar@modifyly.in. Please include "Legal Notice" in the subject line for proper handling. We aim to respond to all legal inquiries within 5-7 business days.',
  },
];

export default function LegalNoticePage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Legal Notice
            </h1>
            <p className="text-surface-400 text-sm">
              Last updated: June 1, 2026
            </p>
          </div>

          {/* Content */}
          <div className="glass rounded-2xl p-6 sm:p-10 divide-y divide-(--border-subtle)">
            {sections.map((section) => (
              <div key={section.title} className="py-6 first:pt-0 last:pb-0">
                <h2 className="text-lg font-semibold text-surface-100 mb-3">
                  {section.title}
                </h2>
                <p className="text-sm text-surface-400 leading-relaxed">
                  {section.content}
                </p>
              </div>
            ))}
          </div>

          {/* Back link */}
          <div className="text-center mt-8">
            <Link
              href="/"
              className="text-sm text-brand-400 hover:text-brand-300 transition-colors"
            >
              &larr; Back to Home
            </Link>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
