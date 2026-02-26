import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Clock, History, Layout, CheckCircle2, ChevronRight } from "lucide-react";
import { Link } from "react-router-dom";

interface DriftEvent {
    id: string;
    version: string;
    type: 'initial' | 'drift' | 'alignment';
    description: string;
    timestamp: string;
    score: number;
}

interface DriftTimelineProps {
    events?: DriftEvent[];
}

export function DriftTimeline({ events, projectId }: DriftTimelineProps & { projectId?: string }) {
    if (!events || events.length === 0) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="text-xl font-bold flex items-center gap-2">
                        <History className="h-6 w-6 text-indigo-500" />
                        Drift & Alignment Timeline
                    </CardTitle>
                    <CardDescription>
                        Chronological history of project snapshots and intelligence analysis loops.
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center justify-center py-16 text-center border-2 border-dashed rounded-xl mx-6 mb-6">
                    <div className="h-12 w-12 rounded-full bg-indigo-50 flex items-center justify-center mb-4">
                        <Clock className="h-6 w-6 text-indigo-400" />
                    </div>
                    <p className="text-lg font-bold text-slate-900">Timeline begins here</p>
                    <p className="text-sm text-slate-500 mt-1 max-w-[280px]">
                        Start importing your project to see real-time drift analysis and alignment history.
                    </p>
                </CardContent>
            </Card>
        );
    }

    const displayEvents = events;

    const getIcon = (type: string) => {
        switch (type) {
            case 'initial': return <Layout className="h-4 w-4 text-blue-500" />;
            case 'drift': return <History className="h-4 w-4 text-amber-500" />;
            case 'alignment': return <CheckCircle2 className="h-4 w-4 text-emerald-500" />;
            default: return <Clock className="h-4 w-4 text-slate-400" />;
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-xl font-bold flex items-center gap-2">
                    <History className="h-6 w-6 text-indigo-500" />
                    Drift & Alignment Timeline
                </CardTitle>
                <CardDescription>
                    Chronological history of project snapshots and intelligence analysis loops.
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="relative space-y-8 before:absolute before:left-[11px] before:top-2 before:h-[calc(100%-16px)] before:w-[2px] before:bg-slate-100">
                    {displayEvents.map((event) => (
                        <div key={event.id} className="relative pl-8">
                            <div className="absolute left-0 top-1 z-10 flex h-6 w-6 items-center justify-center rounded-full border bg-white shadow-sm">
                                {getIcon(event.type)}
                            </div>
                            <Link
                                to={event.type === 'alignment'
                                    ? `/projects/${projectId}/alignment`
                                    : `/projects/${projectId}/bootstrap`
                                }
                                className="flex flex-col gap-1 p-4 rounded-xl border bg-slate-50/50 hover:bg-slate-50 transition-colors cursor-pointer group block"
                            >
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2">
                                        <span className="text-sm font-black text-indigo-600 tracking-tight">{event.version}</span>
                                        <span className="text-[10px] uppercase font-bold text-muted-foreground px-1.5 py-0.5 border rounded bg-white">
                                            {event.type}
                                        </span>
                                    </div>
                                    <span className="text-[10px] text-muted-foreground font-medium">{event.timestamp}</span>
                                </div>
                                <p className="text-sm text-slate-700 font-medium py-1">{event.description}</p>
                                <div className="flex items-center justify-between mt-2 pt-2 border-t border-slate-100">
                                    <div className="flex items-center gap-3">
                                        <div className="flex flex-col">
                                            <span className="text-[10px] uppercase text-muted-foreground font-bold leading-none">Alignment Score</span>
                                            <span className={`text-sm font-bold mt-1 ${event.score > 85 ? 'text-emerald-600' : event.score > 70 ? 'text-amber-600' : 'text-rose-600'}`}>
                                                {event.score}%
                                            </span>
                                        </div>
                                    </div>
                                    <ChevronRight className="h-4 w-4 text-slate-400 group-hover:translate-x-1 transition-transform" />
                                </div>
                            </Link>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}
