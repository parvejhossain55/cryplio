"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import WalletCard from "@/components/dashboard/user/WalletCard";
import RecentActivity from "@/components/dashboard/user/RecentActivity";
import CryplioCard from "@/components/dashboard/user/CryplioCard";
import QuickActions from "@/components/dashboard/user/QuickActions";
import SecurityStatus from "@/components/dashboard/user/SecurityStatus";

const UserDashboard = () => {
    return (
        <DashboardLayout title="Overview" role="user">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column: Stats & Wallet */}
                <div className="lg:col-span-2 space-y-8">
                    <WalletCard />
                    <RecentActivity />
                </div>

                {/* Right Column: Cards & Offers */}
                <div className="space-y-8">
                    <CryplioCard />
                    <QuickActions />
                    <SecurityStatus />
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserDashboard;
