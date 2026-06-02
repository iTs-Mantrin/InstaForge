import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';

export const metadata = {
  title: 'Content Removal Request',
  description: `Content Removal Request for ${APP_NAME} — process for requesting removal of content from our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. Overview',
    content:
      `${APP_NAME} respects the intellectual property rights and privacy of all individuals and organizations. This page outlines the procedure for submitting a content removal request if you believe that your copyrighted content or personal information has been processed through our service without appropriate authorization. Our service acts solely as a technical intermediary that processes publicly accessible Instagram content to generate temporary downloadable links. We do not permanently host, store, or distribute content. Copyright in all content processed through our platform remains exclusively with the original creators and rights holders.`,
  },
  {
    title: '2. Who Can Submit a Request',
    content:
      'The following parties may submit a content removal request: (a) Copyright owners or their authorized representatives who believe their copyrighted content has been improperly accessed or processed through our service; (b) Individuals whose personal data or privacy rights may have been violated through the processing of content containing their personal information; (c) Content creators who have created the content and wish to have it removed. All requestors must provide verifiable proof of their identity and ownership or authorization before we can process the request.',
  },
  {
    title: '3. Information Required for a Removal Request',
    content:
      'To submit a content removal request, please provide the following information: (a) Your full legal name, email address, and mailing address; (b) Proof of your identity (such as a government-issued ID) and proof of your ownership or authorization to act on behalf of the rights holder; (c) The specific Instagram URL or other identifying details of the content you wish to have removed; (d) A detailed description of your concern and the basis for the removal request (such as copyright infringement, privacy violation, or other legal grounds); (e) A statement that you have a good faith belief that the content should be removed; (f) A statement, made under penalty of perjury, that the information provided in your request is accurate and complete. Incomplete requests may be delayed until we receive the necessary information.',
  },
  {
    title: '4. How to Submit a Request',
    content:
      'Content removal requests should be submitted in writing via email to: manish.kumar@modifyly.in. Please include "Content Removal Request" in the subject line to ensure proper handling. We recommend using a clear and descriptive subject line to facilitate prompt processing. You may also submit requests through alternative channels if notified otherwise. We do not currently accept removal requests through social media, phone calls, or informal communication channels.',
  },
  {
    title: '5. Review Process and Expected Response Times',
    content:
      'Upon receiving a content removal request, we will follow this process: (a) Acknowledge receipt of your request within 2 business days; (b) Review your request for completeness and verify the information provided; (c) Evaluate the legal and factual basis for removal; (d) Make a determination and take appropriate action within 5-7 business days of receiving a complete request; (e) Notify you of our decision and any actions taken. If additional information is required, we will contact you and the timeline will be adjusted accordingly. Complex cases may require additional time for review.',
  },
  {
    title: '6. Actions We May Take',
    content:
      'Upon determining that a removal request is valid and substantiated, we may take one or more of the following actions: (a) Remove or disable access to the identified content; (b) Block the specific URL from being processed through our service in the future; (c) Restrict access to the service by the user who submitted the content; (d) Take any other action we deem appropriate under the circumstances. We reserve the right to reject requests that are frivolous, abusive, not legally substantiated, or that do not comply with applicable laws.',
  },
  {
    title: '7. Counter-Request Process',
    content:
      'If you believe that content was removed in error or as a result of a misidentification, you may submit a counter-request. Your counter-request must include: (a) Your full name and contact information; (b) Identification of the content that was removed; (c) A detailed explanation of why you believe the removal was in error; (d) Your consent to the jurisdiction of relevant courts; (e) Your signature (physical or electronic). We will review counter-requests and respond within 7-10 business days.',
  },
  {
    title: '8. Contact Information',
    content:
      'For all content removal requests, counter-requests, and related inquiries, please contact us at: manish.kumar@modifyly.in. We are committed to handling all requests fairly, transparently, and promptly in accordance with applicable laws and our internal policies.',
  },
];

export default function ContentRemovalPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Content Removal Request
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
