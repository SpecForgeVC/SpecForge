import { useEffect, useRef } from "react";
import { type RefinementEvent } from "@/api/refinement";
import { CheckCircle2, Loader2, XCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface RefinementProgressProps {
    events: RefinementEvent[];
    status: 'IN_PROGRESS' | 'VALIDATED' | 'FAILED' | 'APPROVED' | undefined;
}

export function RefinementProgress({ events, status }: RefinementProgressProps) {
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [events]);

    return (
        <div className="space-y-4 border rounded-md p-4 bg-muted/30">
            <div className="flex items-center justify-between">
                <h4 className="text-sm font-semibold">AI Refinement Progress</h4>
                <div className="flex items-center gap-2">
                    {status === 'IN_PROGRESS' && <Loader2 className="h-4 w-4 animate-spin text-blue-500" />}
                    {status === 'VALIDATED' && <CheckCircle2 className="h-4 w-4 text-green-500" />}
                    {status === 'FAILED' && <XCircle className="h-4 w-4 text-red-500" />}
                    <span className="text-xs font-mono uppercase text-muted-foreground">{status?.replace('_', ' ')}</span>
                </div>
            </div>

            <div
                ref={scrollRef}
                className="h-48 overflow-y-auto space-y-2 font-mono text-xs bg-black text-green-400 p-3 rounded shadow-inner"
            >
                {events.length === 0 && <div className="text-gray-500 italic">Waiting for events...</div>}

                {events.map((event, i) => (
                    <div key={i} className="border-l-2 border-transparent pl-2 hover:bg-white/5 transition-colors">
                        <div className="flex gap-2">
                            <span className="text-gray-500">[{new Date(event.timestamp || Date.now()).toLocaleTimeString()}]</span>
                            <span className={cn(
                                "font-bold",
                                event.type === 'ERROR' ? "text-red-500" :
                                    event.type === 'SUCCESS' ? "text-green-500" :
                                        event.type === 'WARN' ? "text-yellow-500" :
                                            "text-blue-400"
                            )}>{event.type}</span>
                            <span>{event.message}</span>
                        </div>
                        {event.payload && (
                            <pre className="mt-1 text-gray-400 text-[10px] overflow-x-auto">
                                {JSON.stringify(event.payload, null, 2)}
                            </pre>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
}
