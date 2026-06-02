import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME, APP_DOMAIN } from '@/lib/constants';

export const metadata = {
  title: 'Disclaimer',
  description: `Disclaimer for ${APP_NAME} — important legal disclaimers regarding the use of our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. General Information Only',
    content:
      `The information, content, and services provided by ${APP_NAME} (${APP_DOMAIN}) are for general informational and educational purposes only. All information on the website is provided in good faith. However, we make no representation or warranty of any kind, express or implied, regarding the accuracy, adequacy, validity, reliability, availability, or completeness of any information on the site.`,
  },
  {
    title: '2. No Professional Advice',
    content:
      'The content and services available on this website do not constitute and should not be construed as professional legal, financial, technical, or any other form of professional advice. You should not rely on any information provided by the service as a substitute for consultation with qualified professionals who are familiar with your specific circumstances. We recommend seeking independent professional advice before taking any action based on the information provided through the service.',
  },
  {
    title: '3. Service Availability and Reliability',
    content:
      `${APP_NAME} is provided on an "as is" and "as available" basis. We do not guarantee that the service will be uninterrupted, timely, secure, or error-free. We reserve the right to modify, suspend, or discontinue any aspect of the service at any time without prior notice. We shall not be liable for any interruption, delay, or failure of the service, including but not limited to downtime, data loss, or technical malfunctions.`,
  },
  {
    title: '4. Third-Party Links and Content',
    content:
      'Our website and service may contain links to third-party websites, platforms, or services that are not owned or controlled by us. We have no control over, and assume no responsibility for, the content, privacy policies, security practices, or terms of use of any third-party websites or services. Your use of third-party websites and services is at your own risk, and you should review their applicable terms and policies. The inclusion of any link does not imply endorsement, affiliation, or sponsorship by us.',
  },
  {
    title: '5. Affiliation Disclaimer',
    content:
      `${APP_NAME} is an independent service and is not affiliated with, endorsed by, sponsored by, or otherwise connected to Instagram, Meta Platforms Inc., or any of their respective subsidiaries, affiliates, or parent companies. All product names, logos, brands, trademarks, and service marks appearing on the website are the property of their respective owners. Any reference to third-party trademarks is for identification purposes only and does not indicate any relationship, sponsorship, or endorsement.`,
  },
  {
    title: '6. Fair Use and Copyright',
    content:
      'Content processed through our service may be subject to fair use, fair dealing, or similar provisions under applicable copyright laws. We do not encourage, condone, or facilitate copyright infringement. Users are solely responsible for ensuring their use of the service and any downloaded content complies with applicable copyright laws and fair use provisions. We respect the intellectual property rights of others and expect our users to do the same.',
  },
  {
    title: '7. Limitation of Liability',
    content:
      'To the fullest extent permitted by applicable law, InstaForge, its operators, affiliates, employees, agents, and contributors shall not be liable for any direct, indirect, incidental, special, consequential, exemplary, or punitive damages, including without limitation loss of profits, data, use, goodwill, or other intangible losses, arising out of or in connection with: (a) your use of or inability to use the service; (b) any content downloaded or otherwise obtained through the service; (c) unauthorized access to or alteration of your data; (d) statements or conduct of any third party on the service; (e) any other matter relating to the service.',
  },
  {
    title: '8. Changes and Updates',
    content:
      'We reserve the right to update, modify, or change this Disclaimer at any time without prior notice. Changes will be effective immediately upon posting the updated version on this page. We encourage you to review this Disclaimer periodically to stay informed of any updates. Your continued use of the service after any changes constitutes acceptance of the revised Disclaimer.',
  },
];

export default function DisclaimerPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Disclaimer
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
