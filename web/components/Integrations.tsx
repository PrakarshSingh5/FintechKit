"use client";

import { motion } from "framer-motion";
import { CreditCard, Landmark, Shield, Coins, Wallet } from "lucide-react";
import { cn } from "@/lib/utils";

const integrations = [
  {
    name: "Stripe",
    description: "Global payment processing and subscription management.",
    icon: CreditCard,
    color: "from-indigo-500 to-purple-500",
    text: "text-indigo-500",
    bg: "bg-indigo-500/10",
    border: "border-indigo-500/20",
    category: "Payments",
  },
  {
    name: "Razorpay",
    description: "Comprehensive payment solutions for Indian market.",
    icon: Wallet,
    color: "from-blue-500 to-cyan-500",
    text: "text-blue-500",
    bg: "bg-blue-500/10",
    border: "border-blue-500/20",
    category: "Payments",
  },
  {
    name: "Plaid",
    description: "Securely connect with thousands of financial institutions.",
    icon: Landmark,
    color: "from-red-500 to-orange-500",
    text: "text-red-500",
    bg: "bg-red-500/10",
    border: "border-red-500/20",
    category: "Banking",
  },
  {
    name: "TrueLayer",
    description: "Open banking APIs for payments and financial data.",
    icon: Shield,
    color: "from-emerald-500 to-teal-500",
    text: "text-emerald-500",
    bg: "bg-emerald-500/10",
    border: "border-emerald-500/20",
    category: "Banking",
  },
  {
    name: "CoinGecko",
    description: "Real-time cryptocurrency data and market analysis.",
    icon: Coins,
    color: "from-yellow-500 to-lime-500",
    text: "text-yellow-500",
    bg: "bg-yellow-500/10",
    border: "border-yellow-500/20",
    category: "Crypto",
  },
];

export function Integrations() {
  return (
    <section className="py-24 relative overflow-hidden">
      <div className="container mx-auto px-4 z-10 relative">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
          >
            <h2 className="text-3xl md:text-5xl font-bold mb-6">
              Supported <span className="text-primary">Integrations</span>
            </h2>
            <p className="text-muted text-lg">
              Unified integration for the world's leading financial platforms.
              Switch providers instantly without rewriting code.
            </p>
          </motion.div>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {integrations.map((integration, index) => (
            <motion.div
              key={integration.name}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              className={cn(
                "group relative p-8 rounded-2xl border bg-secondary/20 backdrop-blur-sm overflow-hidden hover:bg-secondary/30 transition-all duration-300",
                integration.border,
              )}
            >
              {/* Hover Glow Effect */}
              <div
                className={cn(
                  "absolute opacity-0 group-hover:opacity-20 transition-opacity duration-500 -inset-1 blur-2xl bg-gradient-to-r",
                  integration.color,
                )}
              />

              <div className="relative z-10">
                <div className="flex justify-between items-start mb-6">
                  <div
                    className={cn("p-3 rounded-xl inline-flex", integration.bg)}
                  >
                    <integration.icon
                      className={cn(
                        "w-6 h-6",
                        integration.text,
                      )}
                    />
                  </div>
                  <span className="text-xs font-mono py-1 px-3 rounded-full border border-border bg-background/50 text-muted-foreground">
                    {integration.category}
                  </span>
                </div>

                <h3 className="text-xl font-bold mb-3 group-hover:text-primary transition-colors">
                  {integration.name}
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  {integration.description}
                </p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
