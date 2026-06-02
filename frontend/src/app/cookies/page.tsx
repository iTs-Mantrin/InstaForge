import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';

export const metadata = {
  title: 'Cookie Policy',
  description: `Cookie Policy for ${APP_NAME} — how we use cookies and similar technologies on our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. What Are Cookies',
    content:
      'Cookies are small text files that are placed on your computer, tablet, or mobile device when you visit a website. They are widely used to make websites work more efficiently, enhance user experience, and provide information to website owners. Cookies may be set by the website you are visiting (first-party cookies) or by third-party services integrated into the website (third-party cookies). Cookies can be temporary (session cookies) that are deleted when you close your browser, or persistent (permanent cookies) that remain on your device for a set period or until you manually delete them.',
  },
  {
    title: '2. How We Use Cookies',
    content:
      `${APP_NAME} uses cookies and similar tracking technologies for several purposes, including: ensuring the proper functioning of the website, analyzing usage patterns to improve our service, providing security against abuse and automated attacks, and delivering relevant advertising. We do not use cookies to collect your Instagram login credentials, access your private data, or track your browsing activity beyond what is described in this policy. You have the ability to control and manage cookies through your browser settings.`,
  },
  {
    title: '3. Essential/Functional Cookies',
    content:
      'Essential cookies are necessary for the proper functioning of our website. These cookies enable core functionality such as maintaining your session state, remembering your download history locally during your browsing session, and ensuring the technical stability of the platform. Without these cookies, certain features of the service may not function correctly. Essential cookies do not track your activity outside of our website and do not require your prior consent under most privacy regulations.',
  },
  {
    title: '4. Analytics Cookies',
    content:
      'We use analytics cookies to collect information about how visitors interact with our website. These cookies help us understand which pages are most frequently visited, how users navigate the site, how long they spend on different sections, and what technical configurations they use (such as browser type and screen resolution). This information is aggregated and anonymized, and is used solely to improve the performance, usability, and content of our service. Analytics cookies may be set by third-party analytics providers on our behalf.',
  },
  {
    title: '5. Advertising Cookies',
    content:
      'We may use advertising cookies to deliver relevant advertisements to you and to measure the effectiveness of our advertising campaigns. These cookies may be set by our advertising partners and may track your browsing activity across different websites to build a profile of your interests. This enables us to show you advertisements that are more relevant to you. Advertising cookies are typically persistent cookies and may remain on your device for an extended period unless you delete them. You can opt out of targeted advertising through your browser settings or through industry opt-out mechanisms.',
  },
  {
    title: '6. Security Cookies',
    content:
      'Security cookies are used to protect our website and users from malicious activity, unauthorized access, and abusive behavior. These cookies help us: detect and prevent automated bot traffic and scraping attempts; identify and block suspicious requests; protect against denial-of-service attacks; and maintain the overall security and integrity of the platform. Security cookies are essential for the safe operation of the service and do not track user behavior for marketing or analytics purposes.',
  },
  {
    title: '7. Third-Party Services and Cookies',
    content:
      `${APP_NAME} may integrate with third-party services that set their own cookies on your device. These third parties may include: analytics providers (such as Google Analytics) that help us understand usage patterns; advertising networks that serve relevant advertisements; and content delivery networks (CDNs) that ensure fast and reliable delivery of website resources. These third-party services have their own privacy and cookie policies governing their use of cookies and data. We encourage you to review the cookie policies of any third-party services integrated into our website for more detailed information.`,
  },
  {
    title: '8. Managing and Disabling Cookies',
    content:
      'You have the right to control and manage cookies in several ways. Most web browsers allow you to: view the cookies stored on your device; delete individual or all cookies; block cookies from specific websites; block all cookies from being set; and set preferences for cookies before they are set. The method for managing cookies varies by browser. Please consult your browser\'s help documentation or settings menu for specific instructions. Please note that if you choose to disable or reject certain cookies, particularly essential and functional cookies, some features and functionality of our website may be affected or may not work as intended. Disabling analytics or advertising cookies does not affect the core functionality of the service.',
  },
  {
    title: '9. Cookie Retention Period',
    content:
      'The duration that cookies remain on your device varies depending on the type of cookie. Session cookies are temporary and are deleted automatically when you close your browser. Persistent cookies remain on your device for a specified period, which can range from a few hours to several months or years, depending on their purpose. We regularly review our use of cookies and ensure that retention periods are set to the minimum necessary to achieve the relevant purpose. You can delete any cookie at any time through your browser settings.',
  },
  {
    title: '10. Updates to This Policy',
    content:
      'We may update this Cookie Policy from time to time to reflect changes in our use of cookies, technological developments, or applicable legal requirements. Changes will be effective immediately upon posting the updated policy on this page with a revised effective date. We encourage you to review this Cookie Policy periodically to stay informed about how we use cookies and similar technologies.',
  },
  {
    title: '11. Contact Information',
    content:
      'If you have any questions, concerns, or requests regarding this Cookie Policy or our use of cookies and tracking technologies, please contact us at manish.kumar@modifyly.in. We will make reasonable efforts to address your inquiries promptly.',
  },
];

export default function CookiePolicyPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Cookie Policy
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
