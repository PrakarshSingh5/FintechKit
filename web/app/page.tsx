import { Hero } from "@/components/Hero";
import { TrustArchitecture } from "@/components/TrustArchitecture";
import { Features } from "@/components/Features";
import { GettingStarted } from "@/components/GettingStarted";

export default function Home() {
  return (
    <main className="min-h-screen bg-background text-foreground">
      <Hero />
      <TrustArchitecture />
      <Features />
      <GettingStarted />
    </main>
  );
}
