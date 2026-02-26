import React from 'react';
import { useParams } from 'react-router-dom';
import {
    Loader2,
    Layers,
    Database,
    Globe,
    FileCode,
    BarChart3,
    AlertTriangle,
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useLatestBootstrapSnapshot, useBootstrapSnapshots } from '@/hooks/use-bootstrap';
import { BootstrapWizard } from './BootstrapWizard';

export const ProjectBootstrapDashboard: React.FC = () => {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: latestSnapshot, isLoading, error } = useLatestBootstrapSnapshot(projectId);
    const { data: allSnapshots } = useBootstrapSnapshots(projectId);

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-[60vh]">
                <Loader2 className="w-8 h-8 animate-spin text-primary" />
            </div>
        );
    }

    // No snapshot yet — show the wizard
    if (error || !latestSnapshot) {
        return (
            <div className="p-6 space-y-6">
                <div>
                    <h1 className="text-2xl font-bold text-foreground">Project Bootstrap Intelligence</h1>
                    <p className="text-muted-foreground mt-1">
                        Import your existing codebase to build a comprehensive project intelligence snapshot.
                    </p>
                </div>
                <BootstrapWizard projectId={projectId!} />
            </div>
        );
    }

    const snapshot = latestSnapshot;
    const snapshotData = snapshot.snapshot_json || {};
    const modules = snapshotData.modules || [];
    const apis = snapshotData.apis || [];
    const dataModels = snapshotData.data_models || [];
    const contracts = snapshotData.contracts || [];
    const risks = snapshotData.risks || [];
    const techStack = snapshotData.tech_stack || {};

    return (
        <div className="p-6 space-y-8">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-foreground">Project Intelligence</h1>
                    <p className="text-muted-foreground mt-1">
                        Snapshot v{snapshot.version} • Created {new Date(snapshot.created_at).toLocaleDateString()}
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    {allSnapshots && allSnapshots.length > 1 && (
                        <span className="text-xs text-muted-foreground">
                            {allSnapshots.length} snapshots
                        </span>
                    )}
                    <Button variant="outline" size="sm" className="text-xs">
                        Re-Bootstrap
                    </Button>
                </div>
            </div>

            {/* Score Cards */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <ScoreCard
                    label="Architecture Complexity"
                    value={snapshot.architecture_score}
                    icon={<Layers className="w-4 h-4" />}
                />
                <ScoreCard
                    label="Contract Density"
                    value={snapshot.contract_density}
                    icon={<FileCode className="w-4 h-4" />}
                />
                <ScoreCard
                    label="Risk Score"
                    value={snapshot.risk_score}
                    icon={<AlertTriangle className="w-4 h-4" />}
                    inverted
                />
                <ScoreCard
                    label="Alignment Score"
                    value={snapshot.alignment_score}
                    icon={<BarChart3 className="w-4 h-4" />}
                />
            </div>

            {/* Tech Stack */}
            {techStack && Object.keys(techStack).length > 0 && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                            <Globe className="w-4 h-4 text-blue-500" />
                            Tech Stack
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-wrap gap-2">
                            {Object.entries(techStack).flatMap(([category, items]) => {
                                if (!Array.isArray(items)) return [];
                                return (items as string[]).map((item, idx) => (
                                    <span
                                        key={`${category}-${idx}`}
                                        className="px-2 py-1 rounded-md bg-muted text-xs text-foreground border border-border"
                                    >
                                        {item}
                                    </span>
                                ));
                            })}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Modules */}
            {modules.length > 0 && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                            <Layers className="w-4 h-4 text-violet-500" />
                            Modules ({modules.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {modules.map((mod: any, idx: number) => (
                                <div key={idx} className="p-3 rounded-lg bg-muted/50 border border-border">
                                    <div className="flex items-center justify-between mb-1">
                                        <span className="text-sm font-medium text-foreground">{mod.name}</span>
                                        <RiskBadge level={mod.risk_level} />
                                    </div>
                                    {mod.description && (
                                        <p className="text-xs text-muted-foreground line-clamp-2">{mod.description}</p>
                                    )}
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* API Endpoints */}
            {apis.length > 0 && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                            <Globe className="w-4 h-4 text-emerald-500" />
                            API Endpoints ({apis.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2 max-h-64 overflow-auto">
                            {apis.map((api: any, idx: number) => (
                                <div key={idx} className="flex items-center gap-3 p-2 rounded bg-muted/50 border border-border">
                                    <MethodBadge method={api.method} />
                                    <span className="text-xs font-mono text-foreground flex-1 truncate">{api.endpoint}</span>
                                    {api.auth_type && (
                                        <span className="text-[10px] text-muted-foreground bg-muted border border-border px-1.5 py-0.5 rounded">
                                            {api.auth_type}
                                        </span>
                                    )}
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Data Models & Contracts side-by-side */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {dataModels.length > 0 && (
                    <Card className="bg-card border-border shadow-sm">
                        <CardHeader className="pb-3">
                            <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                                <Database className="w-4 h-4 text-amber-500" />
                                Data Models ({dataModels.length})
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-2 max-h-48 overflow-auto">
                                {dataModels.map((dm: any, idx: number) => (
                                    <div key={idx} className="px-3 py-2 rounded bg-muted/50 border border-border">
                                        <span className="text-xs font-medium text-foreground">{dm.name}</span>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                )}

                {contracts.length > 0 && (
                    <Card className="bg-card border-border shadow-sm">
                        <CardHeader className="pb-3">
                            <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                                <FileCode className="w-4 h-4 text-cyan-500" />
                                Contracts ({contracts.length})
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-2 max-h-48 overflow-auto">
                                {contracts.map((c: any, idx: number) => (
                                    <div key={idx} className="flex items-center justify-between px-3 py-2 rounded bg-muted/50 border border-border">
                                        <span className="text-xs font-medium text-foreground">{c.name}</span>
                                        <span className="text-[10px] text-muted-foreground">{c.contract_type}</span>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                )}
            </div>

            {/* Risks */}
            {risks.length > 0 && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-card-foreground flex items-center gap-2">
                            <AlertTriangle className="w-4 h-4 text-red-500" />
                            Identified Risks ({risks.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            {risks.map((risk: any, idx: number) => (
                                <div key={idx} className="flex items-start gap-3 p-3 rounded bg-muted/50 border border-border">
                                    <RiskBadge level={risk.severity} />
                                    <div>
                                        <span className="text-xs font-medium text-foreground">{risk.area}</span>
                                        <p className="text-xs text-muted-foreground mt-0.5">{risk.description}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
};

// --- Helper Components ---

function ScoreCard({ label, value, icon, inverted = false }: { label: string; value: number; icon: React.ReactNode; inverted?: boolean }) {
    const normalized = Math.min(100, Math.max(0, value));
    let color: string;
    if (inverted) {
        color = normalized >= 60 ? 'text-red-500' : normalized >= 30 ? 'text-amber-500' : 'text-emerald-500';
    } else {
        color = normalized >= 75 ? 'text-emerald-500' : normalized >= 40 ? 'text-amber-500' : 'text-red-500';
    }

    return (
        <Card className="bg-card border-border shadow-sm">
            <CardContent className="pt-4">
                <div className="flex items-center gap-2 text-muted-foreground mb-2">
                    {icon}
                    <span className="text-xs">{label}</span>
                </div>
                <div className={`text-3xl font-bold ${color}`}>{normalized.toFixed(0)}</div>
                <div className="w-full h-1.5 bg-secondary rounded-full mt-2">
                    <div
                        className={`h-full rounded-full transition-all ${inverted
                            ? (normalized >= 60 ? 'bg-red-500' : normalized >= 30 ? 'bg-amber-500' : 'bg-emerald-500')
                            : (normalized >= 75 ? 'bg-emerald-500' : normalized >= 40 ? 'bg-amber-500' : 'bg-red-500')
                            }`}
                        style={{ width: `${normalized}%` }}
                    />
                </div>
            </CardContent>
        </Card>
    );
}

function RiskBadge({ level }: { level: string }) {
    const colors: Record<string, string> = {
        'LOW': 'bg-emerald-100 text-emerald-700 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/30',
        'MEDIUM': 'bg-amber-100 text-amber-700 border-amber-200 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/30',
        'HIGH': 'bg-red-100 text-red-700 border-red-200 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/30',
        'CRITICAL': 'bg-red-200 text-red-800 border-red-300 dark:bg-red-600/20 dark:text-red-300 dark:border-red-500/50',
    };
    const colorClass = colors[level?.toUpperCase()] || colors['LOW'];
    return (
        <span className={`text-[10px] px-1.5 py-0.5 rounded border ${colorClass}`}>
            {level || 'UNKNOWN'}
        </span>
    );
}

function MethodBadge({ method }: { method: string }) {
    const colors: Record<string, string> = {
        'GET': 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400',
        'POST': 'bg-blue-100 text-blue-700 dark:bg-blue-500/10 dark:text-blue-400',
        'PUT': 'bg-amber-100 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400',
        'PATCH': 'bg-orange-100 text-orange-700 dark:bg-orange-500/10 dark:text-orange-400',
        'DELETE': 'bg-red-100 text-red-700 dark:bg-red-500/10 dark:text-red-400',
    };
    const colorClass = colors[method?.toUpperCase()] || 'bg-muted text-muted-foreground';
    return (
        <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded w-12 text-center ${colorClass}`}>
            {method}
        </span>
    );
}
