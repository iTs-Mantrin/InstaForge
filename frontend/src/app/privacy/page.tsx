import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME, APP_DOMAIN } from '@/lib/constants';

export const metadata = {
  title: 'Privacy Policy',
  description: `Privacy Policy for ${APP_NAME} — how we collect, use, and protect your data when using our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. Information We Collect',
    content:
      `When you use ${APP_NAME}, we collect minimal technical information necessary to provide and improve the service. This includes your IP address, browser type and version, operating system, referring URLs, and device information. We also collect the Instagram URL you submit for processing — this is required to generate downloadable media links. We use cookies and similar tracking technologies for analytics, security, and functional purposes. Our servers automatically log standard access data including page requests, timestamps, and response statuses. We do not collect your Instagram username, password, or any personal identification required for account creation.`,
  },
  {
    title: '2. How We Use Your Information',
    content:
      `The information we collect is used exclusively to: (a) Process your download requests by fetching publicly accessible Instagram media; (b) Generate temporary downloadable media links; (c) Improve and optimize the performance and reliability of the service; (d) Monitor for abuse, security threats, unauthorized automated access, and violations of our Terms of Service; (e) Compile anonymous, aggregated usage statistics to understand how the service is used. We do not sell, rent, lease, or share your personal information with third parties for their marketing or advertising purposes.`,
  },
  {
    title: '3. Data Storage and Retention',
    content:
      `${APP_NAME} does not permanently store Instagram media, user files, or downloaded content. Media files are temporarily cached on our servers (Cloudflare R2) solely to facilitate immediate download delivery. These temporary files are automatically and permanently deleted shortly after the download is completed or the processing session ends. Server log data, including IP addresses and access timestamps, is retained for a limited period (typically up to 30 days) for security analysis and troubleshooting before being anonymized or permanently deleted. We do not maintain archives of downloaded content.`,
  },
  {
    title: '4. Cookies and Tracking Technologies',
    content:
      `${APP_NAME} uses cookies and similar technologies for the following purposes: (a) Essential/Functional Cookies — required for basic site functionality, such as remembering your download history locally in your browser during your session; (b) Analytics Cookies — we use analytics services to understand how visitors interact with the website, which pages are most frequently visited, and to identify technical issues; (c) Security Cookies — used to detect and prevent abusive activity, automated requests, and security threats. (d) Advertising Cookies — we may use third-party advertising partners to serve relevant advertisements. These partners may set cookies to track your browsing activity across websites to build a profile of your interests. You can control cookie preferences through your browser settings. Disabling certain cookies may affect the functionality and performance of the service.`,
  },
  {
    title: '5. Third-Party Services',
    content:
      `${APP_NAME} utilizes the following third-party services: (a) Cloudflare R2 — for temporary file storage and content delivery; (b) Analytics providers — for understanding usage patterns and improving the service; (c) Advertising networks — for serving relevant advertisements. Each third-party service has its own privacy policy governing data handling practices. We ensure that only the minimum necessary data is shared with these services to fulfill their functions. We encourage you to review the privacy policies of these third-party providers for more information on their data practices.`,
  },
  {
    title: '6. User Responsibility for Downloaded Content',
    content:
      'Users are solely responsible for ensuring they have the legal right to download, store, distribute, or use any content accessed through our service. Downloaded Instagram content remains the intellectual property of its original creator or rights holder. We do not claim ownership over any third-party content processed through our platform. Users should respect applicable copyright laws and obtain proper authorization from content owners before downloading or using their content.',
  },
  {
    title: '7. Data Security',
    content:
      'We implement reasonable and appropriate technical and organizational security measures to protect your data against unauthorized access, alteration, disclosure, or destruction. These measures include encryption in transit using HTTPS protocol, secure server configurations, access controls on our infrastructure, and regular security reviews. However, please be aware that no method of electronic transmission or storage is 100% secure, and we cannot guarantee absolute security of your data.',
  },
  {
    title: '8. Your Rights',
    content:
      'Depending on your jurisdiction, you may have the following rights regarding your personal data: (a) The right to access the personal information we hold about you; (b) The right to request correction of inaccurate information; (c) The right to request deletion of your data; (d) The right to object to or restrict processing of your data; (e) The right to data portability; (f) The right to withdraw consent where processing is based on consent. To exercise any of these rights, please contact us using the information provided below. We will respond to your request within the timeframe required by applicable law.',
  },
  {
    title: '9. Changes to This Policy',
    content:
      'We reserve the right to update or modify this Privacy Policy at any time. Changes will be effective immediately upon posting the revised policy on this page, with the updated revision date noted at the top. We encourage you to review this Privacy Policy periodically to stay informed about how we are protecting your information. Your continued use of the service after any changes constitutes acceptance of the updated policy.',
  },
  {
    title: '10. Contact Information',
    content:
      'If you have any questions, concerns, or requests regarding this Privacy Policy or our data practices, please contact us at manish.kumar@modifyly.in. We will make every effort to address your concerns in a timely and thorough manner.',
  },
  {
    title: '11. No Affiliation',
    content:
      `${APP_NAME} is an independent service and is not affiliated with, endorsed by, sponsored by, or otherwise connected to Instagram, Meta Platforms Inc., or any of their respective subsidiaries, affiliates, or parent companies. All trademarks, service marks, and trade names appearing on our service are the property of their respective owners.`,
  },
];

export default function PrivacyPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Privacy Policy
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
