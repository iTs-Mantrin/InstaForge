import Link from 'next/link';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import { APP_NAME } from '@/lib/constants';
import { Download, Shield, Users, Globe, Sparkles, Lock } from 'lucide-react';

export const metadata = {
  title: 'About Us',
  description: `About ${APP_NAME} — learn about our mission to provide fast, simple, and privacy-respecting Instagram content downloading tools.`,
};

const values = [
  {
    icon: Download,
    title: 'Simplicity & Speed',
    description: 'We believe downloading Instagram content should be effortless. Paste a link, get your files. No sign-ups, no complicated steps, no unnecessary friction.',
  },
  {
    icon: Shield,
    title: 'Privacy First',
    description: 'We do not require Instagram login credentials, we do not access private accounts, and we do not permanently store downloaded content. Your privacy is built into every aspect of our service.',
  },
  {
    icon: Lock,
    title: 'Copyright Respect',
    description: 'We respect the intellectual property rights of content creators. Our service only processes publicly accessible content, and intellectual property rights remain exclusively with the original creators and rights holders.',
  },
  {
    icon: Users,
    title: 'User-Focused',
    description: 'Everything we build is designed with the user in mind. From intuitive interfaces to reliable performance, we strive to deliver a seamless experience for our global community of users.',
  },
  {
    icon: Globe,
    title: 'Globally Accessible',
    description: 'Our service is available to users around the world. We are committed to providing a reliable, fast, and accessible platform regardless of your location.',
  },
  {
    icon: Sparkles,
    title: 'Continuous Improvement',
    description: 'We are constantly working to improve our service, add new features, and enhance performance. We listen to user feedback and evolve our platform to better serve our community.',
  },
];

export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8">
          {/* Hero */}
          <div className="text-center mb-16">
            <h1 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-surface-100 mb-6">
              About {APP_NAME}
            </h1>
            <p className="text-lg text-surface-400 max-w-2xl mx-auto leading-relaxed">
              We provide a fast, simple, and privacy-respecting tool for downloading publicly accessible Instagram content. Our mission is to make content downloading as effortless as possible while respecting intellectual property rights and user privacy.
            </p>
          </div>

          {/* Mission */}
          <div className="glass rounded-2xl p-8 sm:p-10 mb-12">
            <h2 className="text-2xl font-bold text-surface-100 mb-4">
              Our Mission
            </h2>
            <p className="text-sm text-surface-400 leading-relaxed">
              {APP_NAME} was created with a clear mission: to provide a free, unlimited, and reliable tool that enables users to download publicly accessible Instagram content. We believe that accessing content you have the right to view should be simple and straightforward. Our platform processes only publicly accessible content — we do not access private accounts, we do not require Instagram login credentials, and we do not permanently store downloaded content. We are committed to operating transparently, respecting copyright laws, and protecting user privacy at every level of our service.
            </p>
          </div>

          {/* How It Works */}
          <div className="glass rounded-2xl p-8 sm:p-10 mb-12">
            <h2 className="text-2xl font-bold text-surface-100 mb-4">
              How It Works
            </h2>
            <p className="text-sm text-surface-400 leading-relaxed mb-4">
              Our service is straightforward: you paste a public Instagram URL — whether it is a post, reel, or story — and our platform processes that link to generate temporary downloadable media files. We act as a technical intermediary, fetching publicly accessible content at your request and providing you with a downloadable link. We do not host, store, or distribute content beyond the immediate processing required to fulfill your request. The entire process takes seconds and requires no account creation or personal information.
            </p>
            <p className="text-sm text-surface-400 leading-relaxed">
              {APP_NAME} is fully independent and is not affiliated with, endorsed by, or connected to Instagram, Meta Platforms Inc., or any of their subsidiaries or affiliates.
            </p>
          </div>

          {/* Values */}
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-surface-100 text-center mb-10">
              Our Values
            </h2>
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {values.map((value) => (
                <div key={value.title} className="glass rounded-2xl p-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-500/10 text-brand-400 mb-4">
                    <value.icon className="h-5 w-5" />
                  </div>
                  <h3 className="text-sm font-semibold text-surface-100 mb-2">
                    {value.title}
                  </h3>
                  <p className="text-xs text-surface-400 leading-relaxed">
                    {value.description}
                  </p>
                </div>
              ))}
            </div>
          </div>

          {/* Commitment */}
          <div className="glass rounded-2xl p-8 sm:p-10">
            <h2 className="text-2xl font-bold text-surface-100 mb-4">
              Our Commitment
            </h2>
            <p className="text-sm text-surface-400 leading-relaxed mb-4">
              We are committed to maintaining a platform that respects both user privacy and creator rights. We will continue to: (a) Only process publicly accessible content; (b) Never request or store Instagram login credentials; (c) Never permanently store downloaded content; (d) Respect intellectual property rights and respond to valid copyright concerns promptly; and (e) Operate transparently with clear policies and terms.
            </p>
            <p className="text-sm text-surface-400 leading-relaxed">
              If you have any questions, feedback, or concerns, please do not hesitate to contact us. We value the trust our users place in us and are dedicated to earning it every day.
            </p>
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
