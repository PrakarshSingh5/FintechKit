'use client';

import { Check, Copy } from 'lucide-react';
import { useState } from 'react';

export function GettingStarted() {
  const [copied, setCopied] = useState(false);
  const installCmd = "go get github.com/fintechkit/fintechkit";

  const handleCopy = () => {
    navigator.clipboard.writeText(installCmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="py-24 bg-[#0a0a0a] border-t border-border" id="docs">
      <div className="container mx-auto px-4 text-center">
        <h2 className="text-3xl md:text-4xl font-bold mb-8">Ready to build?</h2>
        <p className="text-muted mb-12">Get started in seconds. No API keys to register with us.</p>
        
        <div className="max-w-xl mx-auto">
          <div className="flex items-center justify-between p-4 rounded-xl bg-secondary border border-border">
            <code className="font-mono text-sm md:text-base text-foreground">
              {installCmd}
            </code>
            <button 
              onClick={handleCopy}
              className="p-2 hover:bg-background rounded-lg transition-colors text-muted hover:text-foreground"
            >
              {copied ? <Check className="w-4 h-4 text-green-500" /> : <Copy className="w-4 h-4" />}
            </button>
          </div>
          
          <div className="mt-12">
            <a 
              href="https://github.com/PrakarshSingh5/FintechKit#readme" 
              target="_blank"
              rel="noopener noreferrer"
              className="inline-block bg-foreground text-background px-8 py-4 rounded-full font-bold hover:bg-gray-200 transition-colors text-lg"
            >
              Read the Full Documentation
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
