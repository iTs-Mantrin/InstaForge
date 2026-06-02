import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';

export const metadata = {
  title: 'Copyright Policy',
  description: `Copyright Policy for ${APP_NAME} — our policy regarding copyright ownership, infringement, and takedown procedures.`,
};

const sections = [
  {
    title: '1. Copyright Ownership',
    content:
      `${APP_NAME} respects the intellectual property rights of all content creators and rights holders. All Instagram content processed through our service — including posts, images, videos, reels, and stories — remains the exclusive intellectual property of their respective creators, owners, or authorized rights holders. We do not claim ownership over any third-party content that is accessed, processed, or downloaded through our platform. ${APP_NAME} itself, including its name, branding, logo, software code, design, and original content, is the exclusive property of InstaForge and is protected by applicable copyright and intellectual property laws.`,
  },
  {
    title: '2. Service Role as Technical Intermediary',
    content:
      `${APP_NAME} acts solely as a technical intermediary tool that facilitates the downloading of publicly accessible Instagram content. The service functions as a passive conduit — it does not host, store, publish, distribute, or transmit user-provided content beyond the immediate processing required to generate temporary downloadable links. The service does not initiate the transmission of, select the content of, or modify the content being downloaded. Users provide the specific URL, and the service processes that request on an automated basis without human review of the content. Copyright ownership in all processed content remains exclusively with the original creators and rights holders at all times.`,
  },
  {
    title: '3. User Responsibility for Copyright Compliance',
    content:
      'Users of the service are solely and fully responsible for ensuring they have the legal right to download, store, distribute, display, or otherwise use any content accessed through our platform. It is the user\'s obligation to obtain all necessary permissions, licenses, or authorizations from the copyright owner before downloading or using their content. We strongly recommend that users only download content that they have created themselves, have explicit permission to download, or that is licensed for such use under applicable Creative Commons, public domain, or fair use provisions. We do not assume any responsibility for users who infringe upon the copyrights of others through the use of our service.',
  },
  {
    title: '4. Reporting Copyright Infringement',
    content:
      'If you believe in good faith that your copyrighted work has been used or processed through our service in a manner that constitutes copyright infringement, please submit a written copyright infringement notice to us. To be effective, your notice should include the following information: (a) A physical or electronic signature of the copyright owner or an authorized representative acting on their behalf; (b) Identification of the copyrighted work(s) claimed to have been infringed; (c) Identification of the infringing material and information reasonably sufficient to enable us to locate the material within the service; (d) Your full name, mailing address, telephone number, and email address; (e) A statement that you have a good faith belief that the use of the material in the manner complained of is not authorized by the copyright owner, its agent, or the law; (f) A statement, made under penalty of perjury, that the information in the notice is accurate and that you are the copyright owner or authorized to act on behalf of the copyright owner.',
  },
  {
    title: '5. Review and Response Procedures',
    content:
      'Upon receipt of a valid copyright infringement notice that substantially complies with the requirements above, we will: (a) Acknowledge receipt of the notice within a reasonable timeframe; (b) Review the notice for completeness and sufficiency; (c) Investigate the claimed infringement; (d) Take appropriate action, which may include removing or disabling access to the allegedly infringing content; (e) Notify the submitting party of the actions taken. We reserve the right to request additional information before taking action if the initial notice is incomplete or insufficient.',
  },
  {
    title: '6. Counter-Notification Procedure',
    content:
      'If you believe that content you requested was removed or access to it was disabled as a result of a mistaken or misidentified copyright infringement notice, you may submit a counter-notification. Your counter-notification must include: (a) Your physical or electronic signature; (b) Identification of the material that was removed or to which access was disabled and the location where the material appeared before removal; (c) A statement, made under penalty of perjury, that you have a good faith belief that the material was removed or disabled as a result of mistake or misidentification; (d) Your full name, mailing address, telephone number, and email address; and (e) A statement that you consent to the jurisdiction of the federal court in your judicial district. Upon receipt of a valid counter-notification, we may restore the material in accordance with applicable law.',
  },
  {
    title: '7. Repeat Infringement Policy',
    content:
      `${APP_NAME} reserves the right to terminate or restrict access to the service for users who are determined to be repeat infringers of copyright. A repeat infringer is a user who has been the subject of more than one valid copyright infringement notice. We also reserve the right to remove or disable access to allegedly infringing content without prior notice and to terminate access for users who engage in egregious or willful infringement, even on a first occurrence.`,
  },
  {
    title: '8. Contact Information',
    content:
      'All copyright infringement notices, counter-notifications, and related inquiries should be directed to our designated copyright agent at: manish.kumar@modifyly.in. We will respond to valid notices promptly and take appropriate action as required by applicable copyright laws.',
  },
];

export default function CopyrightPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Copyright Policy
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
