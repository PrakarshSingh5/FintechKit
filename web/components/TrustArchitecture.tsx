'use client';

import { motion } from 'framer-motion';
import { Server, ShieldCheck, Wifi, LockKeyhole, ArrowRight, XCircle } from 'lucide-react';

export function TrustArchitecture() {
  return (
    <section className="py-24 bg-secondary/5 border-y border-border/40" id="security">
      <div className="container mx-auto px-4">
        
        <div className="text-center max-w-2xl mx-auto mb-16">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-green-500/10 border border-green-500/20 text-xs font-medium text-green-500 mb-4">
            <LockKeyhole className="w-3 h-3" />
            Security First Architecture
          </div>
          <h2 className="text-3xl md:text-4xl font-bold mb-4 tracking-tight">Your Infrastructure. Your Keys. Your Data.</h2>
          <p className="text-muted leading-relaxed">
            Unlike SaaS wrappers that route your sensitive financial data through their servers, 
            FintechKit runs entirely within your application. We never see your API keys or user data.
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-12 max-w-5xl mx-auto">
          
          {/* Bad Model (SaaS) */}
          <div className="relative p-8 rounded-2xl bg-secondary/20 border border-border/50 opacity-50 grayscale hover:grayscale-0 transition-all duration-500">
            <div className="absolute top-4 right-4 text-xs font-mono text-red-500 border border-red-500/30 px-2 py-1 rounded bg-red-500/10">TRADITIONAL SAAS</div>
            <h3 className="text-xl font-bold mb-8 flex items-center gap-2">
              <XCircle className="w-5 h-5 text-red-500" />
              Middleman Risk
            </h3>
            
            <div className="flex items-center justify-between relative">
              <div className="flex flex-col items-center gap-2 z-10">
                <div className="w-16 h-16 rounded-xl bg-background border border-border flex items-center justify-center">
                  <Server className="w-8 h-8 text-muted" />
                </div>
                <span className="text-xs font-mono text-muted">Your App</span>
              </div>

              {/* Connection Line */}
              <div className="absolute top-8 left-16 right-16 h-0.5 bg-dashed border-t-2 border-dashed border-red-500/30 w-full" />
              
              <div className="flex flex-col items-center gap-2 z-10">
                <div className="w-16 h-16 rounded-xl bg-red-500/10 border border-red-500/30 flex items-center justify-center relative">
                  <div className="absolute -top-2 -right-2 w-4 h-4 rounded-full bg-red-500 animate-pulse" />
                  <Wifi className="w-8 h-8 text-red-500" />
                </div>
                <span className="text-xs font-bold text-red-500">3rd Party</span>
              </div>

               <div className="flex flex-col items-center gap-2 z-10">
                 <div className="w-16 h-16 rounded-xl bg-slate-100 flex items-center justify-center">
                   <img src="https://upload.wikimedia.org/wikipedia/commons/b/ba/Stripe_Logo%2C_revised_2016.svg" alt="Stripe" className="w-10 opacity-50" />
                 </div>
                 <span className="text-xs font-mono text-muted">Provider</span>
              </div>
            </div>
            
            <div className="mt-8 p-4 bg-red-500/5 rounded-lg text-sm text-red-400">
              ⚠️ User data leaves your infrastructure. Additional point of failure. API keys shared.
            </div>
          </div>

          {/* Good Model (FintechKit) */}
          <div className="relative p-8 rounded-2xl bg-[#0a0a0a] border border-green-500/30 shadow-2xl shadow-green-900/10">
             <div className="absolute top-4 right-4 text-xs font-mono text-green-500 border border-green-500/30 px-2 py-1 rounded bg-green-500/10">FINTECHKIT</div>
             <h3 className="text-xl font-bold mb-8 flex items-center gap-2 text-foreground">
              <ShieldCheck className="w-5 h-5 text-green-500" />
              Direct Connection
            </h3>
            
            <div className="flex items-center justify-center gap-4 relative">
              <div className="flex flex-col items-center gap-2 z-10">
                <div className="relative">
                  <div className="w-20 h-20 rounded-xl bg-secondary/50 border border-border flex items-center justify-center">
                    <Server className="w-8 h-8 text-foreground" />
                  </div>
                  {/* Embedded Library */}
                  <motion.div 
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ repeat: Infinity, duration: 2, repeatType: "reverse" }}
                    className="absolute -bottom-2 -right-2 px-2 py-1 bg-green-500 text-black text-[10px] font-bold rounded shadow-lg"
                  >
                    FintechKit
                  </motion.div>
                </div>
                <span className="text-xs font-mono text-muted">Your App</span>
              </div>

               {/* Direct Arrow */}
              <div className="flex-1 flex items-center justify-center px-4">
                 <div className="h-0.5 w-full bg-gradient-to-r from-green-500 to-transparent relative">
                    <motion.div 
                      animate={{ x: [0, 100], opacity: [1, 0] }} 
                      transition={{ duration: 1.5, repeat: Infinity, ease: "linear" }}
                      className="absolute top-1/2 -translate-y-1/2 w-full h-full bg-green-400 blur-[2px]"
                   />
                 </div>
                 <ArrowRight className="text-green-500 w-5 h-5 ml-[-10px]" />
              </div>

               <div className="flex flex-col items-center gap-2 z-10">
                 <div className="w-16 h-16 rounded-xl bg-white flex items-center justify-center overflow-hidden p-2">
                   <img src="https://upload.wikimedia.org/wikipedia/commons/b/ba/Stripe_Logo%2C_revised_2016.svg" alt="Stripe" className="" />
                 </div>
                 <span className="text-xs font-mono text-muted">Provider</span>
              </div>
            </div>
            
             <div className="mt-8 p-4 bg-green-500/5 rounded-lg text-sm text-green-400 border border-green-500/10">
              ✅ Zero middlemen. End-to-end encryption. Your keys never leave your server.
            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
