import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScoreCard } from './ScoreCard';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Sparkles, ArrowRight, Activity, ExternalLink, AlertTriangle } from 'lucide-react';
import { intelligenceApi, type FeatureIntelligence, type AuditLog, type DriftReport } from '@/api/intelligence';
import { DriftRiskModal } from '@/features/roadmap/components/DriftRiskModal';

interface IntelligencePanelProps {
    roadmapItemId: string;
}

const formatTimeAgo = (dateStr: string): string => {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
};

const getActivityColor = (action: string): string => {
    if (action.includes('DRIFT')) return 'bg-red-500';
    if (action.includes('CREATE')) return 'bg-green-500';
    if (action.includes('UPDATE')) return 'bg-blue-500';
    if (action.includes('DELETE')) return 'bg-yellow-500';
    if (action.includes('APPROVE')) return 'bg-emerald-500';
    if (action.includes('REJECT')) return 'bg-orange-500';
    return 'bg-gray-400';
};

const formatAction = (log: AuditLog): string => {
    const entity = log.entity_type?.replace(/_/g, ' ').toLowerCase() || 'item';
    const action = log.action?.replace(/_/g, ' ').toLowerCase() || 'updated';
    return `${entity} ${action}`;
};

export const IntelligencePanel: React.FC<IntelligencePanelProps> = ({ roadmapItemId }) => {
    const [scores, setScores] = useState<FeatureIntelligence | null>(null);
    const [activity, setActivity] = useState<AuditLog[]>([]);
    const [loadingScores, setLoadingScores] = useState(true);
    const [loadingActivity, setLoadingActivity] = useState(true);
    const [driftModalOpen, setDriftModalOpen] = useState(false);
    const [driftReport] = useState<DriftReport | null>(null);

    useEffect(() => {
        const fetchScores = async () => {
            setLoadingScores(true);
            try {
                const data = await intelligenceApi.getFeatureIntelligence(roadmapItemId);
                setScores(data);
            } catch (err) {
                console.error('Failed to fetch intelligence scores:', err);
                // Fallback to default scores on error
                setScores(null);
            } finally {
                setLoadingScores(false);
            }
        };

        const fetchActivity = async () => {
            setLoadingActivity(true);
            try {
                const logs = await intelligenceApi.getRoadmapItemActivity(roadmapItemId);
                setActivity(logs);
            } catch (err) {
                console.error('Failed to fetch activity:', err);
                setActivity([]);
            } finally {
                setLoadingActivity(false);
            }
        };

        fetchScores();
        fetchActivity();
    }, [roadmapItemId]);

    const handleDriftClick = () => {
        setDriftModalOpen(true);
    };

    if (loadingScores) return <div className="p-4 text-sm text-muted-foreground">Loading Intelligence...</div>;

    const overallScore = scores?.overall_score ?? 0;
    const completenessScore = scores?.completeness_score ?? 0;
    const integrityScore = scores?.contract_integrity_score ?? 0;
    const driftScore = scores?.drift_risk_score ?? 0;
    const driftRisk = 100 - driftScore;

    return (
        <div className="space-y-6">
            <div className="space-y-4">
                <div className="flex justify-between items-center">
                    <h3 className="text-lg font-semibold">Intelligence</h3>
                    <Button variant="ghost" size="sm" className="h-8 gap-1" onClick={() => window.location.href = `/roadmap/${roadmapItemId}/intelligence`}>
                        <ExternalLink className="h-3 w-3" /> Dashboard
                    </Button>
                </div>
            </div>
            <Card className="bg-slate-50 border-slate-200 shadow-sm">
                <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-semibold flex items-center gap-2 text-slate-700">
                        <Sparkles className="h-4 w-4 text-purple-500" />
                        Feature Intelligence
                    </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex flex-col items-center">
                        <ScoreCard score={overallScore} label="Overall Health" />
                    </div>

                    <div className="space-y-3">
                        <div className="flex justify-between items-center text-sm">
                            <span className="text-slate-600">Spec Completeness</span>
                            <span className="font-medium text-slate-900">{completenessScore}%</span>
                        </div>
                        <div className="flex justify-between items-center text-sm">
                            <span className="text-slate-600">Contract Integrity</span>
                            <span className="font-medium text-slate-900">{integrityScore}%</span>
                        </div>
                        <div
                            className="flex justify-between items-center text-sm cursor-pointer hover:bg-slate-100 rounded-md px-2 py-1 -mx-2 transition-colors group"
                            onClick={handleDriftClick}
                            title="Click to view drift analysis"
                        >
                            <span className="text-slate-600 flex items-center gap-1.5">
                                Drift Risk
                                <AlertTriangle className="h-3 w-3 text-amber-500 opacity-0 group-hover:opacity-100 transition-opacity" />
                            </span>
                            <span className={`font-medium ${driftRisk <= 20 ? 'text-green-600' : driftRisk <= 50 ? 'text-yellow-600' : 'text-red-600'}`}>
                                {driftRisk}%
                            </span>
                        </div>
                    </div>

                    <Button variant="outline" className="w-full text-xs" asChild>
                        <Link to={`/roadmap/${roadmapItemId}/intelligence`}>
                            View Full Dashboard <ArrowRight className="ml-2 h-3 w-3" />
                        </Link>
                    </Button>
                </CardContent>
            </Card>

            <Card className="bg-slate-50 border-slate-200 shadow-sm">
                <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-semibold flex items-center gap-2 text-slate-700">
                        <Activity className="h-4 w-4 text-blue-500" />
                        Recent Activity
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {loadingActivity ? (
                        <p className="text-xs text-muted-foreground">Loading activity...</p>
                    ) : activity.length === 0 ? (
                        <p className="text-xs text-muted-foreground">No recent activity</p>
                    ) : (
                        <ul className="space-y-3 text-xs text-slate-600">
                            {activity.slice(0, 5).map((log) => (
                                <li key={log.id} className="flex gap-2">
                                    <span className={`w-2 h-2 rounded-full ${getActivityColor(log.action)} mt-1 shrink-0`}></span>
                                    <span className="flex-1">{formatAction(log)}</span>
                                    <span className="text-muted-foreground shrink-0">{formatTimeAgo(log.created_at)}</span>
                                </li>
                            ))}
                        </ul>
                    )}
                </CardContent>
            </Card>

            <DriftRiskModal
                open={driftModalOpen}
                onOpenChange={setDriftModalOpen}
                roadmapItemId={roadmapItemId}
                driftRiskScore={driftRisk}
                driftReport={driftReport}
            />
        </div>
    );
};
