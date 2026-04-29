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
      <section className="py-24 px-4">
        <div className="container mx-auto max-w-4xl">
          <div className="relative overflow-hidden bg-primary rounded-[48px] p-10 md:p-20 text-center">
            {/* Decoration */}
            <div className="absolute top-0 right-0 w-96 h-96 bg-white/10 rounded-full -translate-y-1/2 translate-x-1/2 blur-3xl" />

            <h2 className="text-4xl md:text-6xl font-black text-white mb-8 relative z-10">Ready to start trading?</h2>
            <p className="text-white/80 text-lg mb-10 max-w-xl mx-auto relative z-10">
              Join thousands of traders on Cryplio today. Create an account in seconds and start swapping crypto for cash.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 relative z-10">
              <button className="w-full sm:w-auto bg-white text-primary px-10 py-5 rounded-2xl text-xl font-bold hover:scale-105 transition-transform shadow-xl">
                Create Account
              </button>
              <button className="w-full sm:w-auto bg-primary-dark text-white px-10 py-5 rounded-2xl text-xl font-bold border border-white/20 hover:bg-white/10 transition-all">
                Contact Sales
              </button>
            </div>
          </div>
        </div>
      </section>

      <Footer />
      <BackToTop />
    </main>
  );
}

