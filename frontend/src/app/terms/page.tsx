import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME, APP_DOMAIN } from '@/lib/constants';

export const metadata = {
  title: 'Terms of Service',
  description: `Terms of Service for ${APP_NAME} — the rules and guidelines governing the use of our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. Acceptance of Terms',
    content:
      `By accessing or using ${APP_NAME} (${APP_DOMAIN}), you agree to be legally bound by these Terms of Service. If you do not agree with any part of these terms, you must immediately discontinue use of the service. These terms apply to all visitors, users, and anyone who accesses or uses the service in any capacity.`,
  },
  {
    title: '2. Description of Service',
    content:
      `${APP_NAME} provides a technical tool that allows users to download publicly accessible Instagram content — including posts, reels, and stories — by pasting a public URL. The service acts solely as a technical intermediary that fetches publicly accessible media and generates temporary downloadable links. We do not permanently host, store, or distribute copyrighted content. The service does not require Instagram login credentials and does not access private accounts or private content.`,
  },
  {
    title: '3. User Responsibilities',
    content:
      'You agree to use the service only for lawful purposes and in full compliance with all applicable local, national, and international laws and regulations. You are solely responsible for ensuring that you have the legal right to download, store, distribute, or use any content accessed through our service. You must not: (a) download copyrighted content without authorization from the rights holder; (b) use the service to infringe upon the intellectual property rights of others; (c) engage in any automated, systematic, or programmatic access to the service including scraping, crawling, or bot activity; (d) attempt to bypass rate limits, access controls, security measures, or any technical restrictions we have implemented; (e) use the service in any manner that could damage, disable, overburden, or impair our servers, network, or infrastructure; (f) attempt to reverse engineer, decompile, disassemble, or derive the source code of the service; (g) impersonate any person or entity, or falsely state or misrepresent your affiliation with any person or entity.',
  },
  {
    title: '4. Intellectual Property Rights',
    content:
      'All content downloaded through our service remains the intellectual property of its original creators, owners, or authorized rights holders. We do not claim ownership over any third-party media processed through our platform. The service name, branding, logo, design, and software code are the exclusive property of InstaForge. All trademarks, service marks, trade names, and logos appearing on the service are the property of their respective owners. You may not use any trademark or trade name appearing on our service without the prior written consent of the respective owner.',
  },
  {
    title: '5. Copyright and Fair Use',
    content:
      `${APP_NAME} respects the intellectual property rights of others. We do not encourage, condone, or facilitate copyright infringement. The service is intended for downloading content you have the legal right to access and download, such as your own uploads, content with explicit permission from the rights holder, or content where downloading is authorized under applicable fair use or fair dealing provisions. Copyright ownership in all downloaded content remains exclusively with the original creator or rights holder.`,
  },
  {
    title: '6. Prohibited Activities',
    content:
      'The following activities are strictly prohibited: (a) copyright infringement or the unauthorized downloading of protected content; (b) automated or programmatic access including the use of bots, scrapers, spiders, crawlers, or any other automated tools; (c) spam or unsolicited bulk communications; (d) distributing malware, viruses, or any other harmful code through the service; (e) attempting to gain unauthorized access to our systems or user data; (f) interfering with or disrupting the integrity or performance of the service; (g) any activity that violates applicable law or regulations; (h) using the service to harass, abuse, or harm others.',
  },
  {
    title: '7. Limitation of Liability',
    content:
      `To the maximum extent permitted by applicable law, ${APP_NAME}, its operators, affiliates, employees, and agents shall not be liable for any direct, indirect, incidental, special, consequential, exemplary, or punitive damages arising out of or relating to your use of or inability to use the service. This includes, without limitation, damages for loss of profits, goodwill, use, data, or other intangible losses, even if we have been advised of the possibility of such damages. The service is provided on an "as is" and "as available" basis without any warranties, express or implied.`,
  },
  {
    title: '8. No Warranty',
    content:
      `${APP_NAME} is provided on an "as is" and "as available" basis without any representations or warranties of any kind, whether express, implied, statutory, or otherwise. We expressly disclaim all warranties, including but not limited to: (a) warranties of merchantability, fitness for a particular purpose, and non-infringement; (b) warranties regarding the availability, reliability, accuracy, completeness, or timeliness of the service; (c) warranties that the service will be uninterrupted, error-free, secure, or free from viruses or other harmful components. No advice or information obtained from us or through the service shall create any warranty not expressly stated in these terms.`,
  },
  {
    title: '9. Account Suspension and Termination',
    content:
      'We reserve the right, in our sole discretion, to suspend, restrict, or permanently terminate your access to the service without prior notice for any reason, including but not limited to: (a) violation of these Terms of Service; (b) engaging in prohibited activities such as automated scraping or abuse; (c) infringement of intellectual property rights; (d) conduct that we believe is harmful to the service, other users, or third parties; (e) illegal or unauthorized use of the service. We also reserve the right to block IP addresses, user agents, or device identifiers associated with abusive or prohibited activity. Termination of access may occur without liability to you.',
  },
  {
    title: '10. Anti-Spam and Anti-Bot Protections',
    content:
      `${APP_NAME} employs automated systems and manual review to detect and prevent spam, automated abuse, bot activity, scraping, and other unauthorized programmatic access. We reserve the right to implement rate limits, CAPTCHAs, IP blocks, and other technical measures to protect the integrity and availability of the service. Circumvention of these protections is strictly prohibited and may result in immediate termination of access.`,
  },
  {
    title: '11. Governing Law',
    content:
      'These Terms of Service and any disputes arising out of or relating to them shall be governed by and construed in accordance with the laws of India, without regard to its conflict of law provisions. You agree to submit to the exclusive personal jurisdiction of the courts located in India for the resolution of any disputes. We reserve the right to bring proceedings in the jurisdiction of your residence or any other relevant jurisdiction.',
  },
  {
    title: '12. Changes to Terms',
    content:
      'We reserve the right to modify, update, or replace these Terms of Service at any time at our sole discretion. Changes will be effective immediately upon posting the revised terms on this page. We encourage you to review these terms periodically. Your continued use of the service after any changes take effect constitutes your acceptance of the modified terms. If you do not agree to the modified terms, you must discontinue use of the service.',
  },
  {
    title: '13. Contact Information',
    content:
      'If you have any questions, concerns, or requests regarding these Terms of Service, please contact us at manish.kumar@modifyly.in. We will make reasonable efforts to address your inquiries promptly.',
  },
];

export default function TermsPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Terms of Service
            </h1>
            <p className="text-surface-400 text-sm">
              Last updated: June 1, 2026
            </p>
          </div>
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
