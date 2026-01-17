import Link from 'next/link';
import { Github, Terminal } from 'lucide-react';
import { Button } from './ui/button'; // Placeholder, will create simple button for now if not exists

export function Navbar() {
  return (
    <nav className="fixed top-0 w-full z-50 border-b border-border/40 bg-background/80 backdrop-blur-md">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="bg-primary/10 p-2 rounded-lg">
            <Terminal className="w-5 h-5 text-primary" />
          </div>
          <span className="font-bold text-lg tracking-tight">FintechKit</span>
        </div>
        
        <div className="flex items-center gap-6 text-sm font-medium text-muted hover:text-foreground transition-colors">
          <Link href="#features" className="hover:text-primary transition-colors">Features</Link>
          <Link href="#security" className="hover:text-primary transition-colors">Security</Link>
          <Link href="https://github.com/PrakarshSingh5/FintechKit" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2 hover:text-primary transition-colors">
            <Github className="w-4 h-4" />
            <span>GitHub</span>
          </Link>
        </div>
      </div>
    </nav>
  );
}
