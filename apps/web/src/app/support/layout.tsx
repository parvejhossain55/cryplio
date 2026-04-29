import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Support Center | Cryplio Help & FAQ",
    description: "Get 24/7 support for your Cryplio account, P2P trading issues, and technical inquiries. Search our FAQs or open a support ticket.",
};

export default function SupportLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return <>{children}</>;
}
