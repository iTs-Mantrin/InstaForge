import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';

export const metadata = {
  title: 'DMCA Takedown Policy',
  description: `DMCA Takedown Policy for ${APP_NAME} — procedures for submitting copyright takedown notices under the Digital Millennium Copyright Act.`,
};

const sections = [
  {
    title: '1. DMCA Compliance Overview',
    content:
      `${APP_NAME} complies with the Digital Millennium Copyright Act (DMCA) and responds to valid notices of alleged copyright infringement in accordance with 17 U.S.C. § 512. This policy outlines the procedures for submitting a DMCA takedown notice and explains how we process and respond to such notices. Our service acts as a technical intermediary that processes publicly accessible Instagram content solely to generate temporary downloadable links. We do not permanently host, store, or distribute media files. Copyright ownership in all content processed through our platform remains with the original creators and rights holders.`,
  },
  {
    title: '2. Designated Copyright Agent',
    content:
      'All DMCA takedown notices and counter-notifications should be submitted in writing to our Designated Copyright Agent. Please submit all DMCA-related correspondence to: manish.kumar@modifyly.in. We acknowledge receipt of valid notices within 2 business days and will take appropriate action as described in this policy.',
  },
  {
    title: '3. Requirements for a Valid DMCA Takedown Notice',
    content:
      'To submit a valid DMCA takedown notice, you must provide a written communication that includes substantially the following information: (a) A physical or electronic signature of the copyright owner or a person authorized to act on their behalf; (b) Identification of the copyrighted work or works claimed to have been infringed; (c) Identification of the infringing material that is to be removed or disabled, and information reasonably sufficient to permit us to locate the material (such as the specific URL or other identifying details); (d) Your full name, mailing address, telephone number, and email address; (e) A statement that you have a good faith belief that the use of the material in the manner complained of is not authorized by the copyright owner, its agent, or the law; (f) A statement that the information in the notification is accurate, and under penalty of perjury, that you are authorized to act on behalf of the owner of an exclusive right that is allegedly infringed. Notices that do not substantially comply with these requirements may be considered deficient and we may request additional information before proceeding.',
  },
  {
    title: '4. Review and Response Procedures',
    content:
      'Upon receiving a DMCA takedown notice that substantially complies with the above requirements, we will take the following steps: (a) Review the notice for completeness and legal sufficiency; (b) Remove or disable access to the allegedly infringing material promptly; (c) Notify the user who submitted or requested the content, if identifiable, that the material has been removed or access disabled; (d) Forward a copy of the takedown notice to the affected user; (e) Maintain a record of the notice as required by applicable law. We reserve the right to reject or request revision of any notice that does not substantially comply with DMCA requirements or that appears to be fraudulent, abusive, or made in bad faith.',
  },
  {
    title: '5. Counter-Notification Procedures',
    content:
      'If you believe that material was removed or access was disabled as a result of a mistake or misidentification, you may submit a counter-notification. A valid counter-notification must include substantially the following: (a) Your physical or electronic signature; (b) Identification of the material that was removed or to which access was disabled and the location where the material appeared before it was removed or disabled; (c) A statement, under penalty of perjury, that you have a good faith belief that the material was removed or disabled as a result of mistake or misidentification; (d) Your full name, mailing address, telephone number, and email address; (e) A statement that you consent to the jurisdiction of the federal district court in your judicial district. Upon receipt of a valid counter-notification, we will forward it to the original complaining party and may restore the removed material within 10-14 business days unless the complaining party files a court action seeking to restrain the activity.',
  },
  {
    title: '6. Repeat Infringer Policy',
    content:
      `${APP_NAME} maintains a policy of terminating or restricting access to the service for users who are determined to be repeat infringers. A user may be considered a repeat infringer if they have been the subject of multiple valid DMCA takedown notices. We also reserve the right to take appropriate action, including immediate termination of access, in cases of egregious or willful infringement. We maintain records of all DMCA notices associated with user activity as part of this policy.`,
  },
  {
    title: '7. Misrepresentations and Abuse',
    content:
      'Pursuant to Section 512(f) of the DMCA, any person who knowingly and materially misrepresents that material or activity is infringing, or that material was removed or disabled by mistake or misidentification, may be held liable for damages, including costs and attorneys\' fees, incurred by the alleged infringer, the copyright owner, or the service provider. Please ensure that all information in your notice or counter-notification is accurate and submitted in good faith.',
  },
  {
    title: '8. Contact Information',
    content:
      'All DMCA-related correspondence should be directed to: manish.kumar@modifyly.in. Please include "DMCA Notice" in the subject line to ensure proper handling. We will respond to valid notices promptly and take appropriate action in accordance with applicable law.',
  },
];

export default function DMCAPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              DMCA Policy
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
