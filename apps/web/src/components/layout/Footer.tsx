import React from "react";
import Link from "next/link";
import { Wallet, Send, Code2, Briefcase, MessageCircle, Terminal, Shield, Activity } from "lucide-react";

const Footer = () => {
    const footerLinks = [
        {
            title: "EXCHANGE",
            links: [
                { name: "P2P Marketplace", href: "/marketplace" },
                { name: "Institutional OTC", href: "#" },
                { name: "Mesh Swap", href: "/swap" },
                { name: "Fee Schedule", href: "#" },
            ],
        },
        {
            title: "PROTOCOL",
            links: [
                { name: "Security Audit", href: "/security" },
                { name: "Proof of Reserves", href: "#" },
                { name: "Developer SDK", href: "/api" },
                { name: "Whitepaper", href: "#" },
            ],
        },
        {
            title: "RESOURCES",
            links: [
                { name: "Command Center", href: "/support" },
                { name: "Status Page", href: "#" },
                { name: "Compliance", href: "#" },
                { name: "Registry", href: "#" },
            ],
        },
    ];

    return (
        <footer className="bg-background border-t border-white/5 pt-24 pb-12 overflow-hidden relative">
            {/* Background Texture */}
            <div className="absolute top-0 right-0 w-1/2 h-full bg-primary/5 blur-[120px] pointer-events-none" />

            <div className="container mx-auto px-4 md:px-6 relative z-10">
                <div className="grid grid-cols-1 md:grid-cols-12 gap-16 mb-20">
                    {/* Brand Column */}
                    <div className="md:col-span-4 space-y-8">
                        <Link href="/" className="flex items-center space-x-3 group">
                            <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center">
                                <Wallet className="text-background w-6 h-6" />
                            </div>
                            <span className="text-2xl font-black italic uppercase tracking-tighter">
                                CRYP<span className="text-primary">LIO</span>
                            </span>
                        </Link>
                        <p className="text-text-dim text-sm font-bold uppercase tracking-widest leading-loose max-w-sm">
                            THE INSTITUTIONAL CLEARING LAYER FOR DECENTRALIZED TRADE. BUILT FOR SCALE, PROTECTED BY SECURE ENCLAVES.
                        </p>

                        <div className="flex items-center gap-4">
                            <div className="px-4 py-2 bg-white/5 border border-white/5 rounded-xl flex items-center gap-2">
                                <Activity className="w-3.5 h-3.5 text-primary animate-pulse" />
                                <span className="text-[10px] font-black text-white uppercase tracking-widest">NETWORK: ACTIVE</span>
                            </div>
                            <div className="px-4 py-2 bg-white/5 border border-white/5 rounded-xl flex items-center gap-2">
                                <Shield className="w-3.5 h-3.5 text-primary" />
                                <span className="text-[10px] font-black text-white uppercase tracking-widest">V2.4.0</span>
                            </div>
                        </div>
                    </div>

                    {/* Links Columns */}
                    <div className="md:col-span-8 grid grid-cols-2 md:grid-cols-3 gap-12">
                        {footerLinks.map((group) => (
                            <div key={group.title} className="space-y-8">
                                <h3 className="text-text-dim text-[10px] font-black uppercase tracking-[0.3em]">{group.title}</h3>
                                <ul className="space-y-4">
                                    {group.links.map((link) => (
                                        <li key={link.name}>
                                            <Link
                                                href={link.href}
                                                className="text-xs font-bold text-white uppercase tracking-widest hover:text-primary transition-all flex items-center gap-2 group"
                                            >
                                                <span className="w-0 h-[1px] bg-primary group-hover:w-3 transition-all" />
                                                {link.name}
                                            </Link>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="pt-12 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-6">
                    <div className="flex items-center gap-6">
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">
                            © {new Date().getFullYear()} CRYPLIO PROTOCOL. ALL NODES REGISTERED.
                        </p>
                    </div>

                    <div className="flex items-center gap-8">
                        {["Twitter", "Discord", "GitHub", "Telegram"].map(social => (
                            <Link key={social} href="#" className="text-[10px] font-black text-text-dim hover:text-white uppercase tracking-widest transition-colors">
                                {social}
                            </Link>
                        ))}
                    </div>

                    <div className="flex items-center gap-8 text-[10px] font-black text-text-dim uppercase tracking-widest">
                        <Link href="#" className="hover:text-white transition-all">Privacy</Link>
                        <Link href="#" className="hover:text-white transition-all">Terms</Link>
                        <Link href="#" className="hover:text-white transition-all">Cookies</Link>
                    </div>
                </div>
            </div>

            {/* Terminal Decoration */}
            <div className="absolute -bottom-10 right-10 flex items-center gap-2 opacity-10">
                <Terminal className="w-4 h-4" />
                <span className="text-[10px] font-mono tracking-tighter">SECURE_UPLINK_ESTABLISHED // PORT_443_SSL_V3</span>
            </div>
        </footer>
    );
};

export default Footer;
