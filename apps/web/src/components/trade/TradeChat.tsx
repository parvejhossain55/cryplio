"use client";

import React, { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Send,
    Paperclip,
    Download,
    Image,
    FileText,
    X,
    Check,
    CheckCheck,
    Loader2
} from "lucide-react";
import { wsService } from "@/services/websocketService";
import { toast } from "sonner";

interface ChatMessage {
    id: string;
    trade_id: string;
    sender_id: string;
    content?: string;
    file_url?: string;
    mime_type?: string;
    file_size?: number;
    created_at: string;
    sender_username?: string;
}

interface TradeChatProps {
    tradeId: string;
    currentUserId: string;
    counterpartUsername?: string;
}

const TradeChat: React.FC<TradeChatProps> = ({
    tradeId,
    currentUserId,
    counterpartUsername
}) => {
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [newMessage, setNewMessage] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [isSending, setIsSending] = useState(false);
    const [isUploading, setIsUploading] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        fetchChatHistory();
        subscribeToTradeChat();
        
        return () => {
            unsubscribeFromTradeChat();
        };
    }, [tradeId]);

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const fetchChatHistory = async () => {
        try {
            const response = await fetch(`/api/trades/${tradeId}/messages`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch chat history");
            }
            
            setMessages(data.messages || []);
        } catch (error) {
            console.error("Error fetching chat history:", error);
            toast.error("Failed to load chat history");
        } finally {
            setIsLoading(false);
        }
    };

    const subscribeToTradeChat = () => {
        // Subscribe to trade chat messages
        wsService.send({
            type: "subscribe_trade",
            trade_id: tradeId
        });

        // Listen for chat messages
        wsService.on("chat_message", (data: ChatMessage) => {
            if (data.trade_id === tradeId) {
                setMessages(prev => [...prev, data]);
            }
        });
    };

    const unsubscribeFromTradeChat = () => {
        wsService.send({
            type: "unsubscribe_trade",
            trade_id: tradeId
        });
    };

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    const sendMessage = async () => {
        if (!newMessage.trim()) return;

        setIsSending(true);
        try {
            const response = await fetch(`/api/trades/${tradeId}/messages`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ content: newMessage.trim() }),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to send message");
            }

            setNewMessage("");
        } catch (error) {
            console.error("Error sending message:", error);
            toast.error("Failed to send message");
        } finally {
            setIsSending(false);
        }
    };

    const sendFileMessage = async (file: File) => {
        setIsUploading(true);
        try {
            const formData = new FormData();
            formData.append("file", file);

            const response = await fetch(`/api/trades/${tradeId}/messages/file`, {
                method: "POST",
                body: formData,
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to send file");
            }

            toast.success("File sent successfully");
        } catch (error) {
            console.error("Error sending file:", error);
            toast.error("Failed to send file");
        } finally {
            setIsUploading(false);
        }
    };

    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            // Check file size (10MB limit)
            if (file.size > 10 * 1024 * 1024) {
                toast.error("File size must be less than 10MB");
                return;
            }
            
            sendFileMessage(file);
        }
    };

    const formatTime = (timestamp: string) => {
        return new Date(timestamp).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit"
        });
    };

    const isOwnMessage = (message: ChatMessage) => {
        return message.sender_id === currentUserId;
    };

    const getFileIcon = (mimeType?: string) => {
        if (!mimeType) return FileText;
        
        if (mimeType.startsWith("image/")) {
            return Image;
        }
        
        return FileText;
    };

    const formatFileSize = (bytes?: number) => {
        if (!bytes) return "";
        
        const sizes = ["B", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + " " + sizes[i];
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64">
                <Loader2 className="w-8 h-8 animate-spin text-primary" />
            </div>
        );
    }

    return (
        <div className="flex flex-col h-full bg-surface border border-white/10 rounded-2xl">
            {/* Chat Header */}
            <div className="flex items-center justify-between p-4 border-b border-white/10">
                <div>
                    <h3 className="text-sm font-bold text-white">Trade Chat</h3>
                    <p className="text-xs text-text-dim">
                        {counterpartUsername ? `Chat with ${counterpartUsername}` : "Chat with trade partner"}
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                    <span className="text-xs text-text-dim">Connected</span>
                </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                <AnimatePresence>
                    {messages.map((message) => {
                        const FileIcon = getFileIcon(message.mime_type);
                        const isOwn = isOwnMessage(message);
                        
                        return (
                            <motion.div
                                key={message.id}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                className={`flex ${isOwn ? "justify-end" : "justify-start"}`}
                            >
                                <div className={`max-w-xs lg:max-w-md ${
                                    isOwn ? "order-2" : "order-1"
                                }`}>
                                    <div className={`px-4 py-2 rounded-2xl ${
                                        isOwn 
                                            ? "bg-primary text-white" 
                                            : "bg-white/10 text-white"
                                    }`}>
                                        {message.content && (
                                            <p className="text-sm whitespace-pre-wrap break-words">
                                                {message.content}
                                            </p>
                                        )}
                                        
                                        {message.file_url && (
                                            <div className="flex items-center space-x-2 mt-2">
                                                <FileIcon className="w-4 h-4" />
                                                <div className="flex-1 min-w-0">
                                                    <p className="text-xs truncate">
                                                        {message.file_url.split("/").pop()}
                                                    </p>
                                                    {message.file_size && (
                                                        <p className="text-xs opacity-75">
                                                            {formatFileSize(message.file_size)}
                                                        </p>
                                                    )}
                                                </div>
                                                <a
                                                    href={message.file_url}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="p-1 hover:bg-white/20 rounded transition-colors"
                                                >
                                                    <Download className="w-3 h-3" />
                                                </a>
                                            </div>
                                        )}
                                        
                                        <div className={`flex items-center justify-end space-x-1 mt-1 ${
                                            isOwn ? "text-primary-100" : "text-text-dim"
                                        }`}>
                                            <span className="text-xs">
                                                {formatTime(message.created_at)}
                                            </span>
                                            {isOwn && (
                                                <CheckCheck className="w-3 h-3" />
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </motion.div>
                        );
                    })}
                </AnimatePresence>
                <div ref={messagesEndRef} />
            </div>

            {/* Message Input */}
            <div className="p-4 border-t border-white/10">
                <div className="flex items-center space-x-2">
                    <input
                        type="file"
                        ref={fileInputRef}
                        onChange={handleFileSelect}
                        className="hidden"
                        accept="image/*,.pdf,.doc,.docx,.txt"
                    />
                    
                    <button
                        onClick={() => fileInputRef.current?.click()}
                        disabled={isUploading}
                        className="p-2 text-text-dim hover:text-white transition-colors disabled:opacity-50"
                        title="Attach file"
                    >
                        {isUploading ? (
                            <Loader2 className="w-5 h-5 animate-spin" />
                        ) : (
                            <Paperclip className="w-5 h-5" />
                        )}
                    </button>
                    
                    <input
                        type="text"
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        onKeyPress={(e) => e.key === "Enter" && sendMessage()}
                        placeholder="Type a message..."
                        className="flex-1 px-4 py-2 bg-white/10 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                        disabled={isSending}
                    />
                    
                    <button
                        onClick={sendMessage}
                        disabled={!newMessage.trim() || isSending}
                        className="p-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isSending ? (
                            <Loader2 className="w-5 h-5 animate-spin" />
                        ) : (
                            <Send className="w-5 h-5" />
                        )}
                    </button>
                </div>
                
                <div className="mt-2 text-xs text-text-dim">
                    Files: Images, PDF, DOC, TXT (Max 10MB)
                </div>
            </div>
        </div>
    );
};

export default TradeChat;
