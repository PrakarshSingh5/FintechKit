'use client';

import { motion } from 'framer-motion';
import { Check, X, Code2, ArrowRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useState } from 'react';

const CodeWindow = ({ title, code, className, isGood }: { title: string, code: string, className?: string, isGood?: boolean }) => (
  <div className={cn("rounded-xl overflow-hidden border shadow-2xl bg-[#09090b]", className)}>
    {/* Window Header */}
    <div className="flex items-center justify-between px-4 py-3 bg-[#18181b] border-b border-[#27272a]">
      <div className="flex items-center gap-2">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-[#ef4444]/20 border border-[#ef4444]/50" />
          <div className="w-3 h-3 rounded-full bg-[#eab308]/20 border border-[#eab308]/50" />
          <div className="w-3 h-3 rounded-full bg-[#22c55e]/20 border border-[#22c55e]/50" />
        </div>
        <span className="text-xs text-zinc-400 font-mono ml-2">{title}</span>
      </div>
      <div className={cn(
        "text-[10px] font-bold px-2 py-0.5 rounded uppercase tracking-wider",
        isGood 
          ? "bg-green-500/10 text-green-500 border border-green-500/20" 
          : "bg-red-500/10 text-red-500 border border-red-500/20"
      )}>
        {isGood ? 'Optimized' : 'Legacy'}
      </div>
    </div>
    
    {/* Code Content */}
    <div className="p-4 overflow-x-auto custom-scrollbar">
      <pre className="font-mono text-sm leading-relaxed">
        <code dangerouslySetInnerHTML={{ __html: code }} />
      </pre>
    </div>
  </div>
);

const traditionalCode = `
<span class="text-purple-400">import</span> Stripe <span class="text-purple-400">from</span> <span class="text-green-400">'stripe'</span>;
<span class="text-purple-400">const</span> stripe = <span class="text-purple-400">new</span> Stripe(<span class="text-green-400">'sk_test_...'</span>);

<span class="text-gray-500">// 1. Handle complexity yourself</span>
<span class="text-purple-400">try</span> {
  <span class="text-purple-400">const</span> payment = <span class="text-purple-400">await</span> stripe.paymentIntents.create({
    amount: <span class="text-blue-400">2000</span>,
    currency: <span class="text-green-400">'usd'</span>,
    payment_method_types: [<span class="text-green-400">'card'</span>],
  });

  <span class="text-gray-500">// 2. Manually handle webhooks</span>
  <span class="text-purple-400">const</span> sig = req.headers[<span class="text-green-400">'stripe-signature'</span>];
  
  <span class="text-gray-500">// 3. Maintain separate logic for new providers</span>
  <span class="text-gray-500">// Repeat this for Razorpay, Plaid...</span>

} <span class="text-purple-400">catch</span> (err) {
  <span class="text-gray-500">// 4. Fragmented error handling</span>
  <span class="text-purple-400">if</span> (err.type === <span class="text-green-400">'StripeCardError'</span>) {
    <span class="text-gray-500">// handle card error</span>
  }
}
`.trim();

const modernCode = `
<span class="text-purple-400">import</span> { fintech } <span class="text-purple-400">from</span> <span class="text-green-400">'fintechkit'</span>;

<span class="text-gray-500">// ONE unified interface for ANY provider</span>
<span class="text-purple-400">const</span> payment = <span class="text-purple-400">await</span> fintech.payment.create({
  provider: <span class="text-green-400">'stripe'</span>, <span class="text-gray-500">// or 'razorpay'</span>
  amount: <span class="text-blue-400">2000</span>,
  currency: <span class="text-green-400">'USD'</span>
});

<span class="text-gray-500">// • Automatic retry logic</span>
<span class="text-gray-500">// • Unified error handling</span>
<span class="text-gray-500">// • Built-in logging</span>
`.trim();

export function IntegrationShowcase() {
  return (
    <section className="py-24 bg-zinc-950/[0.3] relative overflow-hidden">
        {/* Background blobs */}
        <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-red-500/5 blur-[100px] rounded-full pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-[500px] h-[500px] bg-green-500/5 blur-[100px] rounded-full pointer-events-none" />

      <div className="container mx-auto px-4 relative z-10">
        <div className="flex flex-col lg:flex-row gap-16 items-center">
          
          {/* Text Side */}
          <motion.div 
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="flex-1 space-y-8"
          >
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-secondary/50 border border-border text-sm font-medium">
              <Code2 className="w-4 h-4 text-primary" />
              Dev Experience
            </div>
            
            <h2 className="text-3xl md:text-5xl font-bold leading-tight">
              Stop Writing <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-red-400 to-orange-400 line-through decoration-red-500/50 decoration-4">Spaghetti Code</span>
            </h2>
            
            <p className="text-lg text-muted leading-relaxed">
              Every provider has a different API, different error codes, and different webhook formats. 
              <strong className="text-foreground"> FintechKit</strong> creates a standard abstraction layer so you never have to read API references again.
            </p>

            <div className="space-y-4">
              {[
                { text: "70% less code to maintain", good: true },
                { text: "No vendor lock-in", good: true },
                { text: "Standardized error handling", good: true },
                { text: "Fragmented implementations", good: false },
              ].map((item, i) => (
                <div key={i} className="flex items-center gap-3">
                  <div className={cn(
                    "w-6 h-6 rounded-full flex items-center justify-center border",
                    item.good 
                      ? "bg-green-500/10 border-green-500/20 text-green-500" 
                      : "bg-red-500/10 border-red-500/20 text-red-500 opacity-50"
                  )}>
                    {item.good ? <Check className="w-3.5 h-3.5" /> : <X className="w-3.5 h-3.5" />}
                  </div>
                  <span className={cn("text-base", item.good ? "text-foreground" : "text-muted line-through")}>
                    {item.text}
                  </span>
                </div>
              ))}
            </div>

            <a 
              href="https://github.com/PrakarshSingh5/FintechKit"
              target="_blank"
              rel="noopener noreferrer"
              className="mt-4 px-6 py-3 rounded-lg bg-foreground text-background font-bold hover:opacity-90 transition-opacity flex items-center gap-2 w-fit"
            >
              View Documentation <ArrowRight className="w-4 h-4" />
            </a>
          </motion.div>

          {/* Code Side */}
          <motion.div 
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="flex-1 w-full relative"
          >
            <div className="relative space-y-6">
              {/* Bad Code */}
              <motion.div
                initial={{ opacity: 0.5, scale: 0.95, y: 20 }}
                whileInView={{ opacity: 1, scale: 0.95, y: 0 }}
                whileHover={{ scale: 0.98, opacity: 1 }}
                viewport={{ once: true }}
                className="transform origin-bottom lg:translate-x-8 opacity-60 blur-[1px] hover:blur-none hover:opacity-100 transition-all duration-300 z-0"
              >
                <CodeWindow 
                  title="legacy-integration.ts" 
                  code={traditionalCode} 
                  className="border-red-500/20 bg-red-950/[0.05]"
                  isGood={false}
                />
              </motion.div>

              {/* Good Code */}
              <motion.div
                initial={{ opacity: 0, scale: 1, y: 40 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: 0.2 }}
                className="transform lg:-translate-x-4 z-10 shadow-2xl shadow-green-900/20"
              >
                <CodeWindow 
                  title="fintechkit-integration.ts" 
                  code={modernCode} 
                  className="border-green-500/30 bg-black"
                  isGood={true}
                />
              </motion.div>
            </div>
            
            {/* Connecting Arrow (Decorative) */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 opacity-20 pointer-events-none hidden lg:block">
              <svg width="200" height="100" viewBox="0 0 200 100">
                <path d="M100,0 L100,100" stroke="currentColor" strokeWidth="2" strokeDasharray="4 4" />
              </svg>
            </div>

          </motion.div>
        </div>
      </div>
    </section>
  );
}
