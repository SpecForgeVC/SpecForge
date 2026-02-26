import { useParams, Link } from "react-router-dom";
import { useProject } from "@/hooks/use-project";
import { useRoadmapItems } from "@/hooks/use-roadmap-items";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { useUpdateProject } from "@/hooks/use-projects";
import {
    AlertTriangle,
    CheckCircle2,
    Clock,
    ArrowRight,
    ShieldCheck,
    Zap
} from "lucide-react";
import { alignmentApi } from "@/api/alignment";
import { useQuery } from "@tanstack/react-query";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MCPSettingsCard } from "./components/MCPSettingsCard";
import { MCPConnectionGuide } from "./components/MCPConnectionGuide";
import { ArchitectureRiskHeatmap } from "./components/ArchitectureRiskHeatmap";
import { ContractMutationAlerts } from "./components/ContractMutationAlerts";
import { DriftTimeline } from "./components/DriftTimeline";
import { useBootstrapSnapshots, useLatestBootstrapSnapshot } from "@/hooks/use-bootstrap";
import { formatDistanceToNow } from "date-fns";

export default function DashboardPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project, isLoading: projectLoading } = useProject(projectId);
    const { data: roadmapItems, isLoading: roadmapLoading } = useRoadmapItems(projectId);
    const updateProject = useUpdateProject(projectId || "");

    const { data: alignmentReport } = useQuery({
        queryKey: ["alignment", projectId],
        queryFn: () => alignmentApi.getAlignmentReport(projectId!),
        enabled: !!projectId
    });

    const { data: snapshots } = useBootstrapSnapshots(projectId);
    const { data: latestSnapshot } = useLatestBootstrapSnapshot(projectId);

    if (projectLoading || roadmapLoading) {
        return <div className="p-8 text-center italic text-muted-foreground">Loading dashboard...</div>;
    }

    const highRiskItems = roadmapItems?.filter(item => item.risk_level === "HIGH") || [];
    const activeItems = roadmapItems?.filter(item => ["IN_PROGRESS", "IN_REVIEW"].includes(item.status || "")) || [];
    const completeItems = roadmapItems?.filter(item => item.status === "COMPLETE") || [];

    // Map snapshots to timeline events
    const timelineEvents: any[] = snapshots?.map(snap => ({
        id: snap.id,
        version: `v${snap.version}`,
        type: snap.version === 1 ? 'initial' : 'drift',
        description: snap.version === 1 ? 'Initial codebase ingestion completed.' : 'Drift analysis snapshot.',
        timestamp: formatDistanceToNow(new Date(snap.created_at), { addSuffix: true }),
        score: snap.alignment_score
    })) || [];

    // Add alignment report as the latest event if available
    if (alignmentReport) {
        timelineEvents.unshift({
            id: alignmentReport.id,
            version: 'Current',
            type: 'alignment',
            description: 'Post-import alignment verified.',
            timestamp: formatDistanceToNow(new Date(alignmentReport.created_at), { addSuffix: true }),
            score: alignmentReport.alignment_score
        });
    }

    // Map modules to risk data
    const snapshotData = latestSnapshot?.snapshot_json || {};
    const modules = snapshotData.modules || [];
    const riskData = modules.map((m: any) => ({
        module: m.name,
        risk: m.risk_level === 'HIGH' ? 85 : m.risk_level === 'MEDIUM' ? 45 : 15,
        trend: 'stable' as const
    }));

    const alerts = alignmentReport?.conflicts.map(c => ({
        id: c.id,
        type: c.type,
        description: c.description,
        severity: c.severity,
        timestamp: formatDistanceToNow(new Date(c.created_at), { addSuffix: true })
    })) || [];

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">{project?.name || "Project Dashboard"}</h2>
                    <p className="text-muted-foreground">{project?.description || "Project overview and status."}</p>
                </div>
                <Button asChild className="gap-2 bg-indigo-600 hover:bg-indigo-700 text-white">
                    <Link to={`/projects/${projectId}/alignment`}>
                        <Zap className="h-4 w-4" />
                        Alignment Engine
                    </Link>
                </Button>
            </div>

            <Tabs defaultValue="overview" className="space-y-6">
                <TabsList className="bg-slate-100/50 p-1 border">
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="roadmap">Roadmap</TabsTrigger>
                    <TabsTrigger value="reality-anchor" className="flex items-center gap-1.5">
                        <ShieldCheck className="h-3.5 w-3.5" />
                        Reality Anchor
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-6 pt-2">
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Active Items</CardTitle>
                                <Clock className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{activeItems.length}</div>
                                <p className="text-xs text-muted-foreground">Currently in progress or review</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">High Risk</CardTitle>
                                <AlertTriangle className="h-4 w-4 text-red-500" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold text-red-600">{highRiskItems.length}</div>
                                <p className="text-xs text-muted-foreground">Items requiring attention</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Completed</CardTitle>
                                <CheckCircle2 className="h-4 w-4 text-green-500" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold text-green-600">{completeItems.length}</div>
                                <p className="text-xs text-muted-foreground">Successfully closed items</p>
                            </CardContent>
                        </Card>
                        <Card className="cursor-pointer hover:bg-slate-50 transition-colors" onClick={() => window.location.href = `/projects/${projectId}/alignment`}>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Alignment Score</CardTitle>
                                <Zap className="h-4 w-4 text-indigo-500" />
                            </CardHeader>
                            <CardContent>
                                <div className={`text-2xl font-bold ${(alignmentReport?.alignment_score || 0) > 80 ? 'text-green-600' :
                                    (alignmentReport?.alignment_score || 0) > 50 ? 'text-amber-600' : 'text-red-600'
                                    }`}>
                                    {alignmentReport?.alignment_score ?? '--'}%
                                </div>
                                <p className="text-xs text-muted-foreground">Cross-artifact consistency</p>
                            </CardContent>
                        </Card>
                    </div>

                    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
                        <div className="md:col-span-4">
                            <ArchitectureRiskHeatmap data={riskData} />
                        </div>
                        <div className="md:col-span-3">
                            <ContractMutationAlerts alerts={alerts} />
                        </div>
                    </div>

                    <Card>
                        <CardHeader>
                            <CardTitle>Recent Roadmap Items</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {roadmapItems?.slice(0, 5).map((item) => (
                                    <div key={item.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-slate-50 transition-colors">
                                        <div className="flex flex-col">
                                            <span className="font-semibold">{item.title}</span>
                                            <span className="text-xs text-muted-foreground uppercase">{item.type} • {item.priority} PRIORITY</span>
                                        </div>
                                        <div className="flex items-center gap-4">
                                            <div className={`px-2 py-1 rounded text-[10px] font-bold ${item.status === 'COMPLETE' ? 'bg-green-100 text-green-700' :
                                                item.status === 'IN_PROGRESS' ? 'bg-blue-100 text-blue-700' :
                                                    'bg-slate-100 text-slate-700'
                                                }`}>
                                                {item.status}
                                            </div>
                                            <Button asChild variant="ghost" size="icon">
                                                <Link to={`/roadmap/${item.id}`}>
                                                    <ArrowRight className="h-4 w-4" />
                                                </Link>
                                            </Button>
                                        </div>
                                    </div>
                                ))}
                                {roadmapItems?.length === 0 && (
                                    <div className="text-center py-8 text-muted-foreground italic">
                                        No roadmap items found.
                                    </div>
                                )}
                            </div>
                            {roadmapItems && roadmapItems.length > 5 && (
                                <Button variant="link" className="mt-4 p-0">View all items</Button>
                            )}
                        </CardContent>
                    </Card>

                    <Card className="border-indigo-100 bg-indigo-50/10">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <ShieldCheck className="h-5 w-5 text-indigo-600" />
                                Project Settings
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between rounded-lg border bg-background p-4">
                                <div className="space-y-0.5">
                                    <Label className="text-base">AI Self-Evaluation</Label>
                                    <p className="text-sm text-muted-foreground">
                                        Critically review AI-generated artifacts before they are saved.
                                    </p>
                                </div>
                                <Switch
                                    checked={!!project?.settings?.enable_self_evaluation}
                                    onCheckedChange={(checked) => {
                                        updateProject.mutate({
                                            settings: {
                                                ...project?.settings,
                                                enable_self_evaluation: checked
                                            }
                                        });
                                    }}
                                    disabled={updateProject.isPending}
                                />
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="roadmap" className="pt-2">
                    <Card>
                        <CardHeader>
                            <CardTitle>Full Roadmap</CardTitle>
                            <CardDescription>All items planned for this project.</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {roadmapItems?.map((item) => (
                                    <div key={item.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-slate-50 transition-colors">
                                        <div className="flex flex-col">
                                            <span className="font-semibold">{item.title}</span>
                                            <span className="text-xs text-muted-foreground uppercase">{item.type} • {item.priority} PRIORITY</span>
                                        </div>
                                        <div className="flex items-center gap-4">
                                            <div className="px-2 py-1 rounded text-[10px] font-bold bg-slate-100 text-slate-700">
                                                {item.status}
                                            </div>
                                            <Button asChild variant="ghost" size="icon">
                                                <Link to={`/roadmap/${item.id}`}>
                                                    <ArrowRight className="h-4 w-4" />
                                                </Link>
                                            </Button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="reality-anchor" className="space-y-6 pt-2">
                    <div className="grid gap-6 lg:grid-cols-3">
                        <div className="lg:col-span-2">
                            <DriftTimeline events={timelineEvents} projectId={projectId} />
                        </div>
                        <div className="space-y-6">
                            {project && <MCPSettingsCard project={project} />}
                            {project && <MCPConnectionGuide project={project} />}
                        </div>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}
