'use client';

import { motion } from 'framer-motion';
import { ArrowRight, CheckCircle2, Lock } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useState } from 'react';

const CodeBlock = ({ code, language, label }: { code: string; language: string; label: string }) => (
  <div className="relative rounded-xl overflow-hidden bg-[#0a0a0a] border border-border shadow-2xl">
    <div className="flex items-center justify-between px-4 py-3 bg-secondary/30 border-b border-border">
      <div className="flex items-center gap-2">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-red-500/20 border border-red-500/50" />
          <div className="w-3 h-3 rounded-full bg-yellow-500/20 border border-yellow-500/50" />
          <div className="w-3 h-3 rounded-full bg-green-500/20 border border-green-500/50" />
        </div>
        <span className="text-xs text-muted font-mono ml-2">{label}</span>
      </div>
      <div className="text-xs text-muted font-mono uppercase">{language}</div>
    </div>
    <div className="p-4 overflow-x-auto">
      <pre className="font-mono text-sm leading-relaxed text-gray-300">
        <code>{code}</code>
      </pre>
    </div>
  </div>
);

export function Hero() {
  const [activeTab, setActiveTab] = useState<'stripe' | 'plaid'>('stripe');

  const messyStripeCode = `// The Hard Way (Raw Stripe)
const stripe = require('stripe')(key);

try {
  const session = await stripe.checkout.sessions.create({
    payment_method_types: ['card'],
    line_items: items,
    mode: 'payment',
    success_url: successUrl,
    cancel_url: cancelUrl,
  });
} catch (err) {
  // Handle 10 types of errors...
  if (err.type === 'StripeCardError') ...
}
`;

  const cleanFintechKitCode = `// The FintechKit Way
import { fintech } from 'fintechkit';

// One line, built-in reliability
const payment = await fintech.payment.create({
  provider: 'stripe', // Just switch string to 'plaid'
  amount: 5000,
  currency: 'USD'
});

// Returns unified response
// Errors handled automatically
`;

  return (
    <section className="relative pt-32 pb-20 overflow-hidden">
      {/* Background Gradients */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[400px] bg-primary/10 blur-[100px] rounded-full opacity-30 pointer-events-none" />
      
      <div className="container mx-auto px-4 relative z-10">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          
          {/* Left: Content */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="space-y-8"
          >
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-secondary/50 border border-border text-xs font-medium text-primary">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
              </span>
              v1.0 is now available
            </div>

            <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight leading-[1.1]">
              The Unified Interface for <span className="text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-500">Fintech APIs</span>
            </h1>
            
            <p className="text-lg text-muted max-w-lg leading-relaxed">
              Integrate Stripe, Plaid, and TrueLayer in minutes. A type-safe Go framework that runs in 
              <span className="text-foreground font-semibold"> your infrastructure</span>, keeping your keys and data secure.
            </p>

            <div className="flex flex-wrap gap-4">
              <button className="bg-foreground text-background px-8 py-3 rounded-lg font-bold hover:bg-gray-200 transition-colors flex items-center gap-2">
                Get Started
                <ArrowRight className="w-4 h-4" />
              </button>
              <button className="px-8 py-3 rounded-lg font-bold border border-border hover:bg-secondary/50 transition-colors flex items-center gap-2 text-muted hover:text-foreground">
                <Lock className="w-4 h-4" />
                Audit on GitHub
              </button>
            </div>

            <div className="flex items-center gap-6 pt-4">
              {['Open Source', 'Type Safe', 'Zero Fees'].map((item) => (
                <div key={item} className="flex items-center gap-2 text-sm text-muted">
                  <CheckCircle2 className="w-4 h-4 text-primary" />
                  {item}
                </div>
              ))}
            </div>
          </motion.div>

          {/* Right: Code Visualization */}
          <motion.div 
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="relative group"
          >
            {/* Glow Effect */}
            <div className="absolute -inset-1 bg-gradient-to-r from-primary to-transparent opacity-20 blur-xl group-hover:opacity-30 transition-opacity duration-1000" />
            
            <div className="space-y-4">
              <CodeBlock 
                label="❌ The Old Way (Fragmented)" 
                language="javascript" 
                code={messyStripeCode} 
              />
              <motion.div 
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.4 }}
              >
                 <CodeBlock 
                  label="✅ The FintechKit Way (Unified)" 
                  language="javascript" 
                  code={cleanFintechKitCode} 
                />
              </motion.div>
            </div>
          </motion.div>

        </div>
      </div>
    </section>
  );
}
