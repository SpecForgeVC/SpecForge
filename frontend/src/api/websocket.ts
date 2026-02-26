import { getAccessToken } from './client';

export type WebSocketCallback = (data: any) => void;

class WebSocketService {
    private ws: WebSocket | null = null;
    private listeners: Map<string, Set<WebSocketCallback>> = new Map();
    private reconnectInterval: number = 5000;

    connect() {
        if (this.ws) return;

        const token = getAccessToken();
        if (!token) {
            console.warn('[WS] No access token available, delaying connection...');
            setTimeout(() => this.connect(), 2000);
            return;
        }

        const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/v1/ws';

        // Pass token in query param for handshake authentication
        this.ws = new WebSocket(`${wsUrl}?token=${token}`);

        this.ws.onopen = () => {
            console.log('[WS] Connected successfully');
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.dispatch(message.type, message.payload);
            } catch (e) {
                console.error('[WS] Failed to parse message', e);
            }
        };

        this.ws.onclose = (event) => {
            console.log(`[WS] Disconnected (code: ${event.code}), reconnecting in ${this.reconnectInterval}ms...`);
            this.ws = null;
            // If we got a 4001 (custom code maybe?) or just a normal close, we retry
            setTimeout(() => this.connect(), this.reconnectInterval);
        };

        this.ws.onerror = (err) => {
            console.error('[WS] Error occurred', err);
            if (this.ws) {
                this.ws.close();
            }
        };
    }

    subscribe(eventType: string, callback: WebSocketCallback) {
        if (!this.listeners.has(eventType)) {
            this.listeners.set(eventType, new Set());
        }
        this.listeners.get(eventType)?.add(callback);
    }

    unsubscribe(eventType: string, callback: WebSocketCallback) {
        this.listeners.get(eventType)?.delete(callback);
    }

    private dispatch(eventType: string, payload: any) {
        const callbacks = this.listeners.get(eventType);
        if (callbacks) {
            callbacks.forEach(cb => cb(payload));
        }
    }
}

export const webSocketService = new WebSocketService();
