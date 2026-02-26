import type { components } from "@/api/generated/schema";
import { AlertCircle, CheckCircle2, AlertTriangle, ShieldAlert } from "lucide-react";
import { Progress } from "@/components/ui/progress";

type DriftReport = components["schemas"]["DriftReport"];

interface DriftReportViewProps {
    report: DriftReport;
}

export function DriftReportView({ report }: DriftReportViewProps) {
    const riskColor = report.risk_score && report.risk_score > 70 ? "text-red-600" :
        report.risk_score && report.risk_score > 30 ? "text-amber-600" : "text-green-600";

    return (
        <div className="space-y-6">
            <div className="grid gap-4 md:grid-cols-2">
                <div className={`p-6 border rounded-lg flex flex-col items-center justify-center gap-4 ${report.drift_detected ? 'bg-amber-50 border-amber-200' : 'bg-green-50 border-green-200'}`}>
                    {report.drift_detected ? (
                        <ShieldAlert className="h-12 w-12 text-amber-600" />
                    ) : (
                        <CheckCircle2 className="h-12 w-12 text-green-600" />
                    )}
                    <div className="text-center">
                        <h3 className="text-xl font-bold">{report.drift_detected ? "Drift Detected" : "In Sync"}</h3>
                        <p className="text-sm text-muted-foreground">
                            {report.drift_detected ? "The implementation has deviated from the contract." : "Implementation matches the approved spec."}
                        </p>
                    </div>
                </div>
                <div className="p-6 border rounded-lg space-y-4">
                    <h3 className="font-semibold flex items-center gap-2">
                        <AlertCircle className="h-4 w-4" /> Risk Score
                    </h3>
                    <div className="space-y-2">
                        <div className="flex items-center justify-between text-sm">
                            <span>Probability of Failure</span>
                            <span className={`font-bold ${riskColor}`}>{report.risk_score}%</span>
                        </div>
                        <Progress value={report.risk_score} className="h-2" />
                    </div>
                    <p className="text-xs text-muted-foreground">
                        Based on field-level drift, breaking changes, and dependency depth.
                    </p>
                </div>
            </div>

            {report.breaking_changes && report.breaking_changes.length > 0 && (
                <div className="space-y-3">
                    <h3 className="font-semibold text-red-700 flex items-center gap-2">
                        <AlertTriangle className="h-4 w-4" /> Breaking Changes
                    </h3>
                    <div className="border rounded-md divide-y overflow-hidden">
                        {report.breaking_changes.map((change, idx) => (
                            <div key={idx} className="p-4 bg-white flex items-start gap-4">
                                <div className="font-mono text-xs bg-slate-100 px-1.5 py-0.5 rounded border">
                                    {change.field}
                                </div>
                                <div className="text-sm">
                                    {change.issue}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
