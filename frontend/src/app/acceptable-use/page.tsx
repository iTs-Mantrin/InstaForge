import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';

export const metadata = {
  title: 'Acceptable Use Policy',
  description: `Acceptable Use Policy for ${APP_NAME} — rules and restrictions governing proper use of our Instagram downloading service.`,
};

const sections = [
  {
    title: '1. Purpose and Scope',
    content:
      `This Acceptable Use Policy (AUP) sets forth the rules, restrictions, and guidelines that govern your use of ${APP_NAME}. By accessing or using the service, you agree to comply with this policy in full. This policy applies to all users of the service, regardless of how they access or interact with the platform. We reserve the right to update or modify this policy at any time without prior notice. Continued use of the service after changes constitutes acceptance of the updated policy.`,
  },
  {
    title: '2. Permitted Use',
    content:
      'The service is intended solely for lawful, personal, non-commercial use. Permitted uses include downloading publicly accessible Instagram content that: (a) you have created yourself; (b) you have explicit permission from the rights holder to download; (c) is licensed for downloading under applicable Creative Commons, public domain, or open licensing provisions; or (d) is permitted under applicable fair use, fair dealing, or similar copyright exceptions in your jurisdiction. The service may only be accessed through the provided user interface and in accordance with any rate limits or technical restrictions we have implemented.',
  },
  {
    title: '3. Prohibited Activities — Illegal Use',
    content:
      'You are strictly prohibited from using the service for any illegal activities or purposes. This includes, but is not limited to: (a) copyright infringement and the unauthorized downloading of protected content; (b) violating any applicable local, national, or international laws or regulations; (c) distributing, sharing, or selling downloaded content without authorization from the rights holder; (d) using downloaded content for commercial purposes without the necessary licenses or permissions; (e) engaging in any activity that constitutes a criminal offense or gives rise to civil liability.',
  },
  {
    title: '4. Prohibited Activities — Copyright Infringement',
    content:
      'Copyright infringement is strictly prohibited. You must not use the service to download, distribute, or make available any content that infringes upon the intellectual property rights of others. This includes: (a) downloading copyrighted content without authorization from the copyright owner; (b) circumventing technological protection measures implemented by content owners; (c) encouraging or assisting others in infringing activities; (d) using downloaded content in a manner that violates the exclusive rights of the copyright owner. All copyright ownership remains with the original creators and rights holders.',
  },
  {
    title: '5. Prohibited Activities — Automated Abuse and Scraping',
    content:
      'Any form of automated, programmatic, or systematic access to the service is strictly prohibited. This includes: (a) using bots, spiders, crawlers, scrapers, or other automated tools to access the service; (b) sending automated or bulk requests to our servers; (c) systematically downloading content at a rate that exceeds normal human usage patterns; (d) using any automated means to extract, collect, or harvest data from the service; (e) deploying scripts, macros, or other automation tools to interact with the service. We employ automated systems to detect and block such activities, and violations will result in immediate restriction of access.',
  },
  {
    title: '6. Prohibited Activities — Security Violations',
    content:
      'You must not engage in any activity that compromises the security or integrity of the service or its infrastructure. Prohibited activities include: (a) attempting to gain unauthorized access to any part of the service, our servers, or our network; (b) attempting to bypass, disable, or circumvent any security measures, rate limits, access controls, or technical restrictions we have implemented; (c) conducting vulnerability scans, penetration tests, or other security assessments without our express written authorization; (d) interfering with the proper functioning of the service, including denial-of-service attacks, flooding, or overloading our infrastructure; (e) introducing malicious code, viruses, worms, or harmful software through the service.',
  },
  {
    title: '7. User Compliance with Applicable Laws',
    content:
      'You agree to comply with all applicable laws, regulations, and legal requirements when using the service. This includes, but is not limited to: (a) copyright and intellectual property laws; (b) data protection and privacy laws; (c) laws governing electronic communications and online conduct; (d) export control and sanctions laws. It is your responsibility to understand and comply with the laws applicable to your jurisdiction. We make no representations that the service is appropriate or available for use in all locations.',
  },
  {
    title: '8. Monitoring and Enforcement',
    content:
      'We reserve the right to monitor usage of the service to detect and prevent violations of this policy. If we determine, in our sole discretion, that you have violated this Acceptable Use Policy, we may take one or more of the following enforcement actions: (a) issue a warning; (b) implement rate limits or access restrictions; (c) block your IP address, user agent, or device identifiers; (d) temporarily suspend your access to the service; (e) permanently terminate your access to the service; (f) take legal action as appropriate. We reserve the right to take enforcement action without prior notice, particularly in cases of egregious or harmful violations.',
  },
  {
    title: '9. Limitation of Liability',
    content:
      'To the maximum extent permitted by applicable law, InstaForge and its operators shall not be liable for any damages, losses, or expenses arising out of or relating to: (a) your violation of this Acceptable Use Policy; (b) any unauthorized or prohibited use of the service by you or through your access; (c) any enforcement actions we take in response to violations; (d) any suspension or termination of your access to the service. We provide the service "as is" without any warranties regarding the detection or prevention of all prohibited activities.',
  },
  {
    title: '10. Reporting Violations',
    content:
      'If you become aware of any violation of this Acceptable Use Policy, whether by another user or through other means, please report it to us promptly at manish.kumar@modifyly.in. Please include sufficient detail to enable us to investigate the reported violation. We will review all reported violations and take appropriate action as we deem necessary. We appreciate your cooperation in maintaining a safe and lawful environment for all users.',
  },
];

export default function AcceptableUsePage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Acceptable Use Policy
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
