import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Swap Crypto | Cryplio Instant Exchange",
    description: "Swap between hundreds of crypto assets instantly with zero hidden fees and institutional-grade liquidity on Cryplio.",
};

export default function SwapLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return <>{children}</>;
}
