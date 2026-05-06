"use client";

import React from "react";
import { AuthProvider } from "@/context/AuthContext";
import { Toaster } from "sonner";

export default function RootProvider({ children }: { children: React.ReactNode }) {
    return (
        <AuthProvider>
            <Toaster
                position="top-right"
                richColors
                closeButton
                theme="dark"
                toastOptions={{
                    style: {
                        background: 'rgba(15, 15, 20, 0.9)',
                        backdropFilter: 'blur(10px)',
                        border: '1px solid rgba(255, 255, 255, 0.1)',
                        color: '#fff',
                        borderRadius: '1.5rem',
                        fontSize: '12px',
                        fontWeight: '700',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em'
                    },
                }}
            />
            {children}
        </AuthProvider>
    );
}
