import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { AlertCircle, TrendingUp, TrendingDown, Minus, BarChart3 } from "lucide-react";

interface RiskData {
    module: string;
    risk: number; // 0-100
    trend: 'up' | 'down' | 'stable';
}

interface ArchitectureRiskHeatmapProps {
    data?: RiskData[];
}

export function ArchitectureRiskHeatmap({ data }: ArchitectureRiskHeatmapProps) {
    if (!data || data.length === 0) {
        return (
            <Card className="h-full">
                <CardHeader>
                    <CardTitle className="text-lg font-bold flex items-center gap-2">
                        <AlertCircle className="h-5 w-5 text-indigo-500" />
                        Architecture Risk Heatmap
                    </CardTitle>
                    <CardDescription>
                        Intelligence-derived risk levels across architectural boundaries.
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center justify-center py-12 text-center">
                    <div className="h-12 w-12 rounded-full bg-slate-100 flex items-center justify-center mb-4">
                        <BarChart3 className="h-6 w-6 text-slate-400" />
                    </div>
                    <p className="text-sm font-medium text-slate-900">Baseline analysis in progress</p>
                    <p className="text-xs text-slate-500 mt-1 max-w-[200px]">
                        Submit your first project snapshot to see architectural risk analysis.
                    </p>
                </CardContent>
            </Card>
        );
    }

    const displayData = data;

    const getRiskColor = (risk: number) => {
        if (risk === 0) return "bg-slate-500/10 text-slate-700 border-slate-200";
        if (risk < 30) return "bg-emerald-500/10 text-emerald-700 border-emerald-200";
        if (risk < 60) return "bg-amber-500/10 text-amber-700 border-amber-200";
        return "bg-rose-500/10 text-rose-700 border-rose-200";
    };

    const getTrendIcon = (trend: string) => {
        switch (trend) {
            case 'up': return <TrendingUp className="h-3 w-3 text-rose-500" />;
            case 'down': return <TrendingDown className="h-3 w-3 text-emerald-500" />;
            default: return <Minus className="h-3 w-3 text-slate-400" />;
        }
    };

    return (
        <Card className="h-full">
            <CardHeader>
                <CardTitle className="text-lg font-bold flex items-center gap-2">
                    <AlertCircle className="h-5 w-5 text-indigo-500" />
                    Architecture Risk Heatmap
                </CardTitle>
                <CardDescription>
                    Intelligence-derived risk levels across architectural boundaries.
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {displayData.sort((a, b) => b.risk - a.risk).map((item) => (
                        <div key={item.module} className={`p-3 rounded-lg border flex items-center justify-between ${getRiskColor(item.risk)}`}>
                            <div className="flex flex-col">
                                <span className="text-sm font-semibold">{item.module}</span>
                                <div className="flex items-center gap-1.5 mt-0.5">
                                    {getTrendIcon(item.trend)}
                                    <span className="text-[10px] uppercase font-bold tracking-wider opacity-70">
                                        {item.trend === 'up' ? 'Worsening' : item.trend === 'down' ? 'Improving' : 'Stable'}
                                    </span>
                                </div>
                            </div>
                            <div className="flex flex-col items-end">
                                <span className="text-lg font-bold leading-none">{item.risk}%</span>
                                <span className="text-[10px] uppercase font-bold tracking-tight opacity-60">Risk Score</span>
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}
