class WebSocketService {
    private ws: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000;
    private listeners: Map<string, ((data: any) => void)[]> = new Map();
    private isConnecting = false;

    connect() {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
            return;
        }

        const userId = this.getUserId();
        if (!userId) {
            console.log('WebSocket: No user ID found, skipping connection');
            return;
        }

        this.isConnecting = true;
        const baseUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
        const wsUrl = `${baseUrl}?user_id=${userId}`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.isConnecting = false;
                this.reconnectAttempts = 0;
                this.emit('connected', { message: 'Connected to server' });
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.emit(data.type || 'message', data);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.ws.onclose = (event) => {
                console.log('WebSocket disconnected:', event.code, event.reason);
                this.isConnecting = false;
                this.emit('disconnected', { code: event.code, reason: event.reason });
                
                if (this.reconnectAttempts < this.maxReconnectAttempts) {
                    this.scheduleReconnect();
                }
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.isConnecting = false;
                this.emit('error', { error: 'WebSocket connection error' });
            };

        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            this.isConnecting = false;
        }
    }

    private getUserId(): string | null {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('user_id');
        }
        return null;
    }

    private scheduleReconnect() {
        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
        
        console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`);
        
        setTimeout(() => {
            this.connect();
        }, delay);
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    send(data: any) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket is not connected');
        }
    }

    on(event: string, callback: (data: any) => void) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event)!.push(callback);
    }

    off(event: string, callback: (data: any) => void) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            const index = callbacks.indexOf(callback);
            if (index > -1) {
                callbacks.splice(index, 1);
            }
        }
    }

    private emit(event: string, data: any) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            callbacks.forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error('Error in WebSocket callback:', error);
                }
            });
        }
    }

    isConnected() {
        return this.ws && this.ws.readyState === WebSocket.OPEN;
    }
}

export const wsService = new WebSocketService();
