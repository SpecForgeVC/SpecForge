import { useState, useRef } from 'react';
import { refinementApi, type RefinementSession, type RefinementEvent } from '@/api/refinement';
import { getAccessToken, API_BASE_URL } from '@/api/client';

export function useRefinement() {
    const [session, setSession] = useState<RefinementSession | null>(null);
    const [events, setEvents] = useState<RefinementEvent[]>([]);
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const eventSourceRef = useRef<EventSource | null>(null);

    const startSession = async (artifactType: string, targetType: string, prompt: string, contextData: any, maxIterations: number) => {
        setError(null);
        setEvents([]);
        try {
            const newSession = await refinementApi.startSession(artifactType, targetType, prompt, contextData, maxIterations);
            if (!newSession.id) throw new Error("Session ID is missing from response");
            setSession(newSession);
            connectToSession(newSession.id);
            return newSession;
        } catch (e: any) {
            setError(e.message || "Failed to start refinement session");
            throw e;
        }
    };

    const connectToSession = (sessionId: string) => {
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
        }

        const endpoint = refinementApi.getEventsEndpoint(sessionId);
        // Using fetch-based SSE reader or standard EventSource depending on auth.
        // For now, let's use the same fetch logic as in settings to support Auth headers
        // (Assuming standard EventSource might fail if cookie not present)

        setIsConnected(true);
        const url = `${API_BASE_URL}${endpoint}`;

        fetchEvents(url);
    };

    const fetchEvents = async (url: string) => {
        try {
            const token = getAccessToken();
            const headers: HeadersInit = {};
            if (token) headers["Authorization"] = `Bearer ${token}`;

            const response = await fetch(url, { headers });
            if (!response.body) throw new Error("No response body");

            const reader = response.body.getReader();
            const decoder = new TextDecoder();

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split("\n\n");

                for (const line of lines) {
                    if (line.startsWith("data: ")) {
                        const dataStr = line.replace("data: ", "");
                        if (dataStr === "{}") continue;

                        try {
                            const event = JSON.parse(dataStr) as RefinementEvent;
                            setEvents(prev => [...prev, event]);

                            // Check for terminal states
                            if (event.type === 'SUCCESS') {
                                if (event.payload?.artifact) {
                                    setSession(prev => prev ? {
                                        ...prev,
                                        result: event.payload.artifact,
                                        status: 'VALIDATED',
                                        // @ts-ignore - evaluation added dynamically from event
                                        evaluation: event.payload.evaluation
                                    } : null);
                                }
                            }
                            if (event.type === 'ERROR' && !event.payload?.retry) { // Fatal error
                                setSession(prev => prev ? { ...prev, status: 'FAILED' } : null);
                            }

                        } catch (e) {
                            console.error("Failed to parse event data", e);
                        }
                    } else if (line.startsWith("event: done")) {
                        setIsConnected(false);
                        return;
                    }
                }
            }
        } catch (e: any) {
            console.error("SSE Error:", e);
            setError(e.message);
            setIsConnected(false);
        }
    };

    const reset = () => {
        setSession(null);
        setEvents([]);
        setError(null);
        setIsConnected(false);
        if (eventSourceRef.current) eventSourceRef.current.close();
    };

    return {
        session,
        events,
        isConnected,
        error,
        startSession,
        reset
    };
}
