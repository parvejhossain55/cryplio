"use client";

import React from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

const Pagination = ({ currentPage, totalPages, onPageChange }: PaginationProps) => {
    if (totalPages <= 1) return null;

    const getPageNumbers = () => {
        const pages = [];
        const maxVisiblePages = 5;

        if (totalPages <= maxVisiblePages) {
            for (let i = 1; i <= totalPages; i++) {
                pages.push(i);
            }
        } else {
            pages.push(1);
            if (currentPage > 3) {
                pages.push("...");
            }

            const start = Math.max(2, currentPage - 1);
            const end = Math.min(totalPages - 1, currentPage + 1);

            for (let i = start; i <= end; i++) {
                pages.push(i);
            }

            if (currentPage < totalPages - 2) {
                pages.push("...");
            }
            pages.push(totalPages);
        }
        return pages;
    };

    return (
        <div className="flex items-center justify-center space-x-1 sm:space-x-2 mt-12 px-2">
            <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl sm:rounded-2xl bg-surface border border-border flex items-center justify-center text-text-dim hover:text-white hover:border-primary transition-all disabled:opacity-50 disabled:cursor-not-allowed group"
            >
                <ChevronLeft className="w-4 h-4 sm:w-5 sm:h-5 group-hover:-translate-x-1 transition-transform" />
            </button>

            <div className="flex items-center space-x-1 sm:space-x-2">
                {getPageNumbers().map((page, index) => (
                    <React.Fragment key={index}>
                        {page === "..." ? (
                            <span className="w-6 sm:w-10 text-center text-text-dim font-black text-[10px] sm:text-xs tracking-widest">...</span>
                        ) : (
                            <button
                                onClick={() => onPageChange(page as number)}
                                className={`w-10 h-10 sm:w-12 sm:h-12 rounded-xl sm:rounded-2xl border font-black text-xs sm:text-sm transition-all ${currentPage === page
                                    ? "bg-primary border-primary text-white shadow-lg shadow-primary/20 scale-105 sm:scale-110 z-10"
                                    : "bg-surface border-border text-text-dim hover:text-white hover:border-primary"
                                    }`}
                            >
                                {page}
                            </button>
                        )}
                    </React.Fragment>
                ))}
            </div>

            <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl sm:rounded-2xl bg-surface border border-border flex items-center justify-center text-text-dim hover:text-white hover:border-primary transition-all disabled:opacity-50 disabled:cursor-not-allowed group"
            >
                <ChevronRight className="w-4 h-4 sm:w-5 sm:h-5 group-hover:translate-x-1 transition-transform" />
            </button>
        </div>
    );
};

export default Pagination;
