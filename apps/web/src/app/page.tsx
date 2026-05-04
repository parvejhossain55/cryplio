import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import Hero from "@/components/sections/Hero";
import MarketOverview from "@/components/sections/MarketOverview";
import Features from "@/components/sections/Features";
import PriceTicker from "@/components/sections/PriceTicker";
import HowItWorks from "@/components/sections/HowItWorks";
import FAQ from "@/components/sections/FAQ";
import MobileApp from "@/components/sections/MobileApp";
import BackToTop from "@/components/ui/BackToTop";

export default function Home() {
  return (
    <main className="min-h-screen">
      <Navbar />
      <div className="pt-[72px]">
        <PriceTicker />
      </div>
      <Hero />
      <MarketOverview />
      <div id="features">
        <Features />
      </div>
      <HowItWorks />
      <FAQ />
      <MobileApp />

      {/* Social Proof Section */}
      <section className="py-20 bg-background text-center px-4">
        <div className="container mx-auto">
          <div className="glass p-12 rounded-[40px] border-border">
            <h2 className="text-3xl md:text-5xl font-bold mb-12">Trusted by 2M+ Traders</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-10">
              <div>
                <p className="text-4xl md:text-5xl font-black text-primary mb-2">$4.2B</p>
                <p className="text-sm font-bold text-text-dim uppercase tracking-widest">Total Volume</p>
              </div>
              <div>
                <p className="text-4xl md:text-5xl font-black text-accent mb-2">12M+</p>
                <p className="text-sm font-bold text-text-dim uppercase tracking-widest">Transactions</p>
              </div>
              <div>
                <p className="text-4xl md:text-5xl font-black text-primary mb-2">50+</p>
                <p className="text-sm font-bold text-text-dim uppercase tracking-widest">Countries</p>
              </div>
              <div>
                <p className="text-4xl md:text-5xl font-black text-accent mb-2">24/7</p>
                <p className="text-sm font-bold text-text-dim uppercase tracking-widest">Live Support</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-32 px-4 md:px-6">
        <div className="container mx-auto max-w-5xl">
          <div className="relative overflow-hidden bg-surface border border-white/5 rounded-[3rem] p-12 md:p-24 text-center">
            {/* Ambient Glows */}
            <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-primary/10 rounded-full blur-[120px] -translate-y-1/2 translate-x-1/2" />
            <div className="absolute bottom-0 left-0 w-[500px] h-[500px] bg-secondary/5 rounded-full blur-[120px] translate-y-1/2 -translate-x-1/2" />

            <div className="relative z-10 space-y-8">
              <h2 className="text-5xl md:text-8xl font-black italic uppercase tracking-tighter leading-[0.8] text-white">
                READY TO <br />
                <span className="gradient-text">OPERATE?</span>
              </h2>
              <p className="text-text-dim text-xl font-medium max-w-xl mx-auto leading-relaxed uppercase tracking-tight italic">
                Join the global pool of institutional traders. Initialize your clearing account in less than 120 seconds.
              </p>
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
                <button className="w-full sm:w-auto bg-white text-background px-12 py-5 rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-2xl shadow-white/5">
                  INITIALIZE ACCOUNT
                </button>
                <button className="w-full sm:w-auto bg-surface-light text-white border border-white/10 px-12 py-5 rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/5 transition-all">
                  CONTACT DEPLOYMENT
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <Footer />
      <BackToTop />
    </main>
  );
}

