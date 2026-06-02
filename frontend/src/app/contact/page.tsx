'use client';

import { motion } from 'framer-motion';
import { MessageSquare, Clock, Shield, Briefcase, Wrench, ArrowRight, Mail } from 'lucide-react';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';

const fadeUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
};

const stagger = {
  animate: { transition: { staggerChildren: 0.1 } },
};

const contactCategories = [
  {
    icon: MessageSquare,
    title: 'General Support',
    description: 'Have a question about using the service, need help with a download, or want to report a bug? Our support team is here to help.',
    email: 'manish.kumar@modifyly.in',
    responseTime: 'within 24 hours on business days',
  },
  {
    icon: Shield,
    title: 'Copyright Inquiries',
    description: 'For copyright infringement notices, DMCA takedown requests, or any intellectual property concerns. Please include relevant details and supporting documentation.',
    email: 'manish.kumar@modifyly.in',
    responseTime: 'within 2-3 business days',
    subjectPrefix: '[Copyright]',
  },
  {
    icon: Briefcase,
    title: 'Business Partnerships',
    description: 'Interested in partnering with us, advertising opportunities, or business collaborations? We would love to hear from you.',
    email: 'manish.kumar@modifyly.in',
    responseTime: 'within 3-5 business days',
    subjectPrefix: '[Business]',
  },
  {
    icon: Wrench,
    title: 'Technical Issues',
    description: 'Experiencing technical difficulties, errors, or performance problems? Our technical team will investigate and resolve your issue.',
    email: 'manish.kumar@modifyly.in',
    responseTime: 'within 24 hours on business days',
    subjectPrefix: '[Technical]',
  },
];

export default function ContactPage() {
  return (
    <div className="min-h-screen flex flex-col bg-surface-950">
      <Navbar />
      <main className="flex-1 pt-24 pb-16">
        <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
            className="text-center mb-12"
          >
            <h1 className="text-3xl sm:text-4xl font-bold text-surface-100 mb-4">
              Contact Us
            </h1>
            <p className="text-surface-400 text-lg max-w-2xl mx-auto">
              We are here to help. Choose the appropriate contact channel below for the fastest response.
            </p>
          </motion.div>

          {/* Contact Categories */}
          <motion.div
            variants={stagger}
            initial="initial"
            animate="animate"
            className="grid gap-6 sm:grid-cols-2 mb-12"
          >
            {contactCategories.map((category) => (
              <motion.div
                key={category.title}
                variants={fadeUp}
                className="glass rounded-2xl p-6 hover:ring-1 hover:ring-brand-500/20 transition-all"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-500/10 text-brand-400 mb-4">
                  <category.icon className="h-5 w-5" />
                </div>
                <h3 className="text-base font-semibold text-surface-100 mb-2">
                  {category.title}
                </h3>
                <p className="text-xs text-surface-400 leading-relaxed mb-4">
                  {category.description}
                </p>
                <div className="space-y-2">
                  <a
                    href={`mailto:${category.email}?subject=${category.subjectPrefix ? encodeURIComponent(category.subjectPrefix + ' ') : ''}${encodeURIComponent(category.title + ' Inquiry')}`}
                    className="flex items-center gap-2 text-sm text-accent-400 hover:text-accent-300 transition-colors group"
                  >
                    <Mail className="h-3.5 w-3.5" />
                    <span>{category.email}</span>
                    <ArrowRight className="h-3 w-3 group-hover:translate-x-0.5 transition-transform" />
                  </a>
                  <div className="flex items-center gap-2 text-xs text-surface-500">
                    <Clock className="h-3 w-3" />
                    <span>Response: {category.responseTime}</span>
                  </div>
                </div>
              </motion.div>
            ))}
          </motion.div>

          {/* Important Note */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.3 }}
            className="glass rounded-2xl p-6"
          >
            <h2 className="text-base font-semibold text-surface-100 mb-2">
              Before Contacting Us
            </h2>
            <p className="text-xs text-surface-400 leading-relaxed">
              For quick answers to common questions, please visit our <a href="/#faq" className="text-accent-400 hover:text-accent-300 transition-colors">FAQ section</a>. 
              For copyright and legal inquiries, please use the dedicated Copyright Inquiries channel above and include all relevant details to help us process your request efficiently. 
              We aim to respond to all legitimate inquiries as promptly as possible.
            </p>
          </motion.div>

          <div className="text-center mt-8">
            <a
              href="/"
              className="text-sm text-brand-400 hover:text-brand-300 transition-colors"
            >
              &larr; Back to Home
            </a>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
