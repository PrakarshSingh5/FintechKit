'use client';

import { Zap, Shield, Globe, Terminal, RefreshCw, FileCheck } from 'lucide-react';
import { motion } from 'framer-motion';

const features = [
  {
    icon: Globe,
    title: "Unified Interface",
    description: "One API for Stripe, Plaid, and TrueLayer. Switch providers by changing a single string."
  },
  {
    icon: Shield,
    title: "Compliance Ready",
    description: "Built-in GDPR and PCI helpers. We handle the regulatory complexity so you don't have to."
  },
  {
    icon: RefreshCw,
    title: "Auto-Retries & Resiliency",
    description: "Automatic exponential backoff, circuit breakers, and idempotency keys handling out of the box."
  },
  {
    icon: Terminal,
    title: "Type-Safe SDK",
    description: "Written in Go with full type definitions. Catch errors at compile time, not in production."
  },
  {
    icon: Zap,
    title: "Zero Latency Overhead",
    description: "FintechKit is a library, not a proxy. No network hops to our servers. Maximum performance."
  },
  {
    icon: FileCheck,
    title: "Webhook Verification",
    description: "Secure webhook handling with automatic signature verification for all supported providers."
  }
];

export function Features() {
  return (
    <section className="py-24 bg-background" id="features">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">Everything you need to ship fintech apps</h2>
          <p className="text-muted max-w-2xl mx-auto">
            Stop building boilerplate. FintechKit provides the critical infrastructure layers you usually forget to build.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ delay: idx * 0.1 }}
              viewport={{ once: true }}
              className="p-6 rounded-2xl border border-border bg-secondary/5 hover:bg-secondary/10 transition-colors"
            >
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4 text-primary">
                <feature.icon className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-bold mb-2">{feature.title}</h3>
              <p className="text-muted leading-relaxed text-sm">
                {feature.description}
              </p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
