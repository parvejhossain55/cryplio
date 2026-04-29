import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import RootProvider from "@/components/providers/RootProvider";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export const metadata: Metadata = {
  title: "Cryplio",
  description: "P2P Crypto Exchange",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.variable} antialiased`}>
        <RootProvider>
          {children}
        </RootProvider>
      </body>
    </html>
  );
}
