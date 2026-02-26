import React, { useState } from 'react';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { intelligenceApi, type DriftReport, type DriftFix } from '@/api/intelligence';
import { AlertTriangle, CheckCircle2, Loader2, Sparkles, XCircle, ChevronDown, ChevronUp } from 'lucide-react';

interface DriftRiskModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    roadmapItemId: string;
    driftRiskScore: number;
    driftReport?: DriftReport | null;
}

export const DriftRiskModal: React.FC<DriftRiskModalProps> = ({
    open,
    onOpenChange,
    roadmapItemId,
    driftRiskScore,
    driftReport,
}) => {
    const [fixes, setFixes] = useState<DriftFix[]>([]);
    const [loading, setLoading] = useState(false);
    const [expandedFix, setExpandedFix] = useState<number | null>(null);
    const [appliedFixes, setAppliedFixes] = useState<Set<number>>(new Set());
    const [error, setError] = useState<string | null>(null);

    const getRiskLevel = (score: number) => {
        if (score <= 20) return { label: 'Low Risk', color: 'bg-green-500/10 text-green-600 border-green-500/20' };
        if (score <= 50) return { label: 'Medium Risk', color: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20' };
        return { label: 'High Risk', color: 'bg-red-500/10 text-red-600 border-red-500/20' };
    };

    const riskLevel = getRiskLevel(driftRiskScore);

    const handleGenerateFixes = async () => {
        if (!driftReport) return;
        setLoading(true);
        setError(null);
        try {
            const result = await intelligenceApi.generateDriftFixes(driftReport, roadmapItemId);
            setFixes(result);
        } catch (err: any) {
            setError(err.message || 'Failed to generate fixes');
        } finally {
            setLoading(false);
        }
    };

    const handleApplyFix = (index: number) => {
        setAppliedFixes(prev => new Set(prev).add(index));
    };

    const handleDismissFix = (index: number) => {
        setFixes(prev => prev.filter((_, i) => i !== index));
    };

    const handleClose = () => {
        setFixes([]);
        setAppliedFixes(new Set());
        setExpandedFix(null);
        setError(null);
        onOpenChange(false);
    };

    const breakingChanges = driftReport?.breaking_changes || [];
    const hasDrift = driftReport?.drift_detected ?? false;

    return (
        <Dialog open={open} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[600px] max-h-[85vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <AlertTriangle className="h-5 w-5 text-amber-500" />
                        Drift Risk Analysis
                    </DialogTitle>
                    <DialogDescription>
                        Detailed drift analysis for this feature with AI-suggested fixes
                    </DialogDescription>
                </DialogHeader>

                {/* Risk Score Overview */}
                <div className="rounded-lg border p-4 space-y-3">
                    <div className="flex items-center justify-between">
                        <span className="text-sm font-medium text-muted-foreground">Drift Risk Score</span>
                        <Badge variant="outline" className={riskLevel.color}>
                            {riskLevel.label}
                        </Badge>
                    </div>
                    <div className="flex items-end gap-2">
                        <span className="text-4xl font-bold tabular-nums">{driftRiskScore}</span>
                        <span className="text-sm text-muted-foreground mb-1">/ 100</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                        <div
                            className={`h-2 rounded-full transition-all duration-500 ${driftRiskScore <= 20 ? 'bg-green-500' :
                                driftRiskScore <= 50 ? 'bg-yellow-500' : 'bg-red-500'
                                }`}
                            style={{ width: `${driftRiskScore}%` }}
                        />
                    </div>
                </div>

                {/* Drift Status */}
                <div className="rounded-lg border p-4">
                    <div className="flex items-center gap-2 mb-3">
                        {hasDrift ? (
                            <XCircle className="h-4 w-4 text-red-500" />
                        ) : (
                            <CheckCircle2 className="h-4 w-4 text-green-500" />
                        )}
                        <span className="text-sm font-medium">
                            {hasDrift ? 'Drift Detected' : 'No Drift Detected'}
                        </span>
                    </div>

                    {breakingChanges.length > 0 && (
                        <div className="space-y-2">
                            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                                Breaking Changes ({breakingChanges.length})
                            </span>
                            <div className="space-y-1.5">
                                {breakingChanges.map((bc, i) => (
                                    <div key={i} className="flex items-start gap-2 text-xs p-2 rounded bg-red-500/5 border border-red-500/10">
                                        <span className="font-mono font-medium text-red-600 shrink-0">{bc.field}</span>
                                        <span className="text-muted-foreground">—</span>
                                        <span className="text-slate-700">{bc.issue}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    {!hasDrift && breakingChanges.length === 0 && (
                        <p className="text-xs text-muted-foreground">
                            All contracts are aligned with the approved specification.
                        </p>
                    )}
                </div>

                {/* Generate Fixes Button */}
                {hasDrift && fixes.length === 0 && (
                    <Button
                        onClick={handleGenerateFixes}
                        disabled={loading}
                        className="w-full gap-2"
                    >
                        {loading ? (
                            <>
                                <Loader2 className="h-4 w-4 animate-spin" />
                                Generating Fixes...
                            </>
                        ) : (
                            <>
                                <Sparkles className="h-4 w-4" />
                                Generate Suggested Fixes
                            </>
                        )}
                    </Button>
                )}

                {error && (
                    <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                        {error}
                    </div>
                )}

                {/* Fixes List */}
                {fixes.length > 0 && (
                    <div className="space-y-3">
                        <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">Suggested Fixes ({fixes.length})</span>
                            <span className="text-xs text-muted-foreground">
                                {appliedFixes.size} applied
                            </span>
                        </div>
                        {fixes.map((fix, i) => (
                            <div key={i} className="rounded-lg border p-3 space-y-2">
                                <div
                                    className="flex items-center justify-between cursor-pointer"
                                    onClick={() => setExpandedFix(expandedFix === i ? null : i)}
                                >
                                    <div className="flex items-center gap-2">
                                        {appliedFixes.has(i) ? (
                                            <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />
                                        ) : (
                                            <AlertTriangle className="h-4 w-4 text-amber-500 shrink-0" />
                                        )}
                                        <span className="text-sm font-mono font-medium">{fix.field}</span>
                                    </div>
                                    {expandedFix === i ? (
                                        <ChevronUp className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                        <ChevronDown className="h-4 w-4 text-muted-foreground" />
                                    )}
                                </div>

                                {expandedFix === i && (
                                    <div className="space-y-2 pt-2 border-t">
                                        <div>
                                            <span className="text-xs font-medium text-muted-foreground">Issue</span>
                                            <p className="text-xs text-slate-700">{fix.issue}</p>
                                        </div>
                                        <div>
                                            <span className="text-xs font-medium text-muted-foreground">Suggested Change</span>
                                            <p className="text-xs text-slate-700 font-mono bg-muted/50 rounded p-1.5">
                                                {fix.suggested_change}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="text-xs font-medium text-muted-foreground">Explanation</span>
                                            <p className="text-xs text-slate-600">{fix.explanation}</p>
                                        </div>
                                        {!appliedFixes.has(i) && (
                                            <div className="flex gap-2 pt-1">
                                                <Button size="sm" variant="default" className="h-7 text-xs" onClick={() => handleApplyFix(i)}>
                                                    Apply Fix
                                                </Button>
                                                <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => handleDismissFix(i)}>
                                                    Dismiss
                                                </Button>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                )}

                <DialogFooter>
                    <Button variant="outline" onClick={handleClose}>
                        Close
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};
