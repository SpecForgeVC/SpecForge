import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { ShieldAlert, ArrowRight, Zap } from "lucide-react";
import { Button } from "@/components/ui/button";

interface MutationAlert {
    id: string;
    type: string;
    description: string;
    severity: 'INFO' | 'WARNING' | 'ERROR' | 'CRITICAL' | 'LOW' | 'MEDIUM' | 'HIGH';
    timestamp: string;
}

interface ContractMutationAlertsProps {
    alerts?: MutationAlert[];
}

export function ContractMutationAlerts({ alerts }: ContractMutationAlertsProps) {
    if (!alerts || alerts.length === 0) {
        return (
            <Card className="h-full">
                <CardHeader className="pb-3">
                    <CardTitle className="text-lg font-bold flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            <ShieldAlert className="h-5 w-5 text-rose-500" />
                            Intelligence Alerts
                        </div>
                    </CardTitle>
                    <CardDescription>
                        Real-time alignment conflicts detected via MCP reality snapshots.
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center justify-center py-10 text-center">
                    <div className="h-10 w-10 rounded-full bg-emerald-50 flex items-center justify-center mb-3">
                        <Zap className="h-5 w-5 text-emerald-500" />
                    </div>
                    <p className="text-sm font-medium text-slate-900">System aligned</p>
                    <p className="text-xs text-slate-500 mt-1 max-w-[200px]">
                        No active conflicts detected. Your contracts match the codebase reality.
                    </p>
                </CardContent>
            </Card>
        );
    }

    const displayAlerts = alerts;

    const getSeverityColor = (severity: string) => {
        switch (severity.toUpperCase()) {
            case 'CRITICAL': return "bg-red-600";
            case 'HIGH': return "bg-rose-500";
            case 'MEDIUM': return "bg-amber-500";
            case 'ERROR': return "bg-red-500";
            default: return "bg-slate-500";
        }
    };

    return (
        <Card className="h-full">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg font-bold flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <ShieldAlert className="h-5 w-5 text-rose-500" />
                        Intelligence Alerts
                    </div>
                    <span className="text-[10px] bg-rose-100 text-rose-600 px-2 py-0.5 rounded-full font-bold">
                        {displayAlerts.length} ACTIVE
                    </span>
                </CardTitle>
                <CardDescription>
                    Real-time alignment conflicts detected via MCP reality snapshots.
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-3">
                    {displayAlerts.map((alert) => (
                        <div key={alert.id || Math.random()} className="relative pl-4 border-l-2 py-1">
                            <div className={`absolute left-[-2.5px] top-2 h-1.5 w-1.5 rounded-full ${getSeverityColor(alert.severity)}`} />
                            <div className="flex items-start justify-between gap-4">
                                <div className="flex flex-col">
                                    <div className="flex items-center gap-2">
                                        <span className="text-xs font-black uppercase tracking-tighter opacity-70 italic">{alert.type}</span>
                                        {alert.timestamp && (
                                            <span className="text-[10px] text-muted-foreground">• {alert.timestamp}</span>
                                        )}
                                    </div>
                                    <p className="text-sm mt-0.5 leading-snug">{alert.description}</p>
                                </div>
                                <Button variant="ghost" size="icon" className="h-8 w-8 shrink-0 hover:bg-slate-100 text-indigo-600">
                                    <ArrowRight className="h-4 w-4" />
                                </Button>
                            </div>
                        </div>
                    ))}
                </div>
                <Button variant="link" className="w-full mt-4 text-xs font-bold uppercase tracking-widest text-indigo-600 hover:text-indigo-700 h-auto p-0 flex items-center gap-1.5">
                    View Alignment Report
                    <Zap className="h-3 w-3" />
                </Button>
            </CardContent>
        </Card>
    );
}
