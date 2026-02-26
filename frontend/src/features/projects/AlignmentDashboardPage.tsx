import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import ReactFlow, {
    type Node,
    type Edge,
    Controls,
    Background,
    useNodesState,
    useEdgesState,
} from 'reactflow';
import 'reactflow/dist/style.css';
import dagre from 'dagre';
import { alignmentApi } from '../../api/alignment';
import { roadmapApi } from '../../api/roadmap';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    AlertTriangle,
    RefreshCw,
    Layers,
    CheckCircle2,
    CircleAlert,
    Info,
} from "lucide-react";
import type { Conflict, Severity } from '../../types/alignment';

const getSeverityColor = (severity: Severity) => {
    switch (severity) {
        case 'CRITICAL': return 'bg-red-600 text-white';
        case 'ERROR': return 'bg-red-100 text-red-700 border-red-200';
        case 'WARNING': return 'bg-amber-100 text-amber-700 border-amber-200';
        case 'INFO': return 'bg-blue-100 text-blue-700 border-blue-200';
        default: return 'bg-slate-100 text-slate-700';
    }
};

const getConflictIcon = (severity: Severity) => {
    switch (severity) {
        case 'CRITICAL': return <CircleAlert className="h-5 w-5 text-red-600" />;
        case 'ERROR': return <AlertTriangle className="h-5 w-5 text-red-500" />;
        case 'WARNING': return <AlertTriangle className="h-5 w-5 text-amber-500" />;
        default: return <Info className="h-5 w-5 text-blue-500" />;
    }
};

export const AlignmentDashboardPage: React.FC = () => {
    const { projectId } = useParams<{ projectId: string }>();
    const queryClient = useQueryClient();
    const [nodes, setNodes, onNodesChange] = useNodesState([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);
    const [selectedConflict, setSelectedConflict] = useState<Conflict | null>(null);

    const { data: report, isLoading: reportLoading } = useQuery({
        queryKey: ['alignment', projectId],
        queryFn: () => alignmentApi.getAlignmentReport(projectId!),
        enabled: !!projectId
    });

    const { data: roadmapItems } = useQuery({
        queryKey: ['roadmap', projectId],
        queryFn: () => roadmapApi.listRoadmapItems(projectId!),
        enabled: !!projectId
    });

    const { data: dependencies } = useQuery({
        queryKey: ['roadmap-dependencies', projectId],
        queryFn: () => alignmentApi.listRoadmapDependencies(projectId!),
        enabled: !!projectId
    });

    const triggerCheck = useMutation({
        mutationFn: () => alignmentApi.triggerAlignmentCheck(projectId!),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['alignment', projectId] });
        }
    });

    // Build Graph
    useEffect(() => {
        if (!roadmapItems || !dependencies) return;

        const dagreGraph = new dagre.graphlib.Graph();
        dagreGraph.setGraph({ rankdir: 'LR' });
        dagreGraph.setDefaultEdgeLabel(() => ({}));

        const flowNodes: Node[] = roadmapItems.map((item: any) => {
            const node = {
                id: item.id,
                data: { label: item.title, type: item.type },
                position: { x: 0, y: 0 },
                width: 180,
                height: 60,
                style: {
                    background: item.type === 'FEATURE' ? '#f0f9ff' : '#f8fafc',
                    border: '1px solid #cbd5e1',
                    borderRadius: '8px',
                    padding: '10px',
                    fontSize: '12px',
                    fontWeight: '600'
                }
            };
            dagreGraph.setNode(item.id, { width: 180, height: 60 });
            return node;
        });

        const flowEdges: Edge[] = dependencies.data.map((dep) => {
            dagreGraph.setEdge(dep.source_id, dep.target_id);
            return {
                id: dep.id,
                source: dep.source_id,
                target: dep.target_id,
                animated: true,
                label: dep.dependency_type,
                style: { stroke: '#94a3b8' }
            };
        });

        // Add conflict-based virtual edges if they exist
        report?.conflicts.filter(c => c.type === 'DEPENDENCY_LOOP').forEach((_) => {
            // We could highlight existing edges or add red ones
        });

        dagre.layout(dagreGraph);

        const layoutedNodes = flowNodes.map((node) => {
            const nodeWithPosition = dagreGraph.node(node.id);
            return {
                ...node,
                position: {
                    x: nodeWithPosition.x - 90,
                    y: nodeWithPosition.y - 30,
                },
            };
        });

        setNodes(layoutedNodes);
        setEdges(flowEdges);
    }, [roadmapItems, dependencies, report, setNodes, setEdges]);

    if (reportLoading) return <div className="p-8 text-center italic text-muted-foreground">Analyzing project alignment...</div>;

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Intelligence Alignment</h2>
                    <p className="text-muted-foreground">Cross-artifact consistency and dependency integrity.</p>
                </div>
                <div className="flex items-center gap-2">
                    <div className="flex flex-col items-end mr-4">
                        <span className="text-sm font-medium text-muted-foreground">Alignment Score</span>
                        <div className={`text-2xl font-bold ${(report?.alignment_score || 0) > 80 ? 'text-green-600' :
                            (report?.alignment_score || 0) > 50 ? 'text-amber-600' : 'text-red-600'
                            }`}>
                            {report?.alignment_score ?? '--'}%
                        </div>
                    </div>
                    <Button
                        onClick={() => triggerCheck.mutate()}
                        disabled={triggerCheck.isPending}
                        variant="outline"
                        className="gap-2"
                    >
                        <RefreshCw className={`h-4 w-4 ${triggerCheck.isPending ? 'animate-spin' : ''}`} />
                        Run Analysis
                    </Button>
                </div>
            </div>

            <Tabs defaultValue="overview" className="w-full">
                <TabsList className="grid w-full grid-cols-3 lg:w-[400px]">
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="graph">Dependency Graph</TabsTrigger>
                    <TabsTrigger value="resolutions">Resolutions</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="mt-4 space-y-4">
                    <div className="grid gap-4 md:grid-cols-2">
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <AlertTriangle className="h-5 w-5 text-amber-500" />
                                    Detected Conflicts
                                </CardTitle>
                                <CardDescription>
                                    Found {report?.conflicts.length || 0} potential issues.
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {report?.conflicts.map((conflict, idx) => (
                                        <div
                                            key={idx}
                                            className="flex flex-col border rounded-lg p-3 cursor-pointer hover:bg-slate-50 transition-colors"
                                            onClick={() => setSelectedConflict(conflict)}
                                        >
                                            <div className="flex items-center justify-between mb-2">
                                                <Badge className={getSeverityColor(conflict.severity)}>
                                                    {conflict.severity}
                                                </Badge>
                                                <span className="text-[10px] uppercase font-bold text-slate-500">{conflict.type}</span>
                                            </div>
                                            <p className="text-sm font-medium">{conflict.description}</p>
                                        </div>
                                    ))}
                                    {(!report?.conflicts || report.conflicts.length === 0) && (
                                        <div className="text-center py-8 text-muted-foreground italic">
                                            No conflicts detected. Great job!
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Layers className="h-5 w-5 text-blue-500" />
                                    Schema Overlaps
                                </CardTitle>
                                <CardDescription>
                                    Shared data models and reused properties.
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {report?.overlaps.map((overlap, idx) => (
                                        <div key={idx} className="border rounded-lg p-3 bg-slate-50">
                                            <h4 className="text-sm font-semibold mb-1">{overlap.type}</h4>
                                            <div className="flex flex-wrap gap-1 mt-2">
                                                {overlap.shared_fields.map((f, i) => (
                                                    <Badge key={i} variant="secondary" className="text-[10px]">{f}</Badge>
                                                ))}
                                            </div>
                                            <p className="text-xs text-muted-foreground mt-2">{overlap.description}</p>
                                        </div>
                                    ))}
                                    {(!report?.overlaps || report.overlaps.length === 0) && (
                                        <div className="text-center py-8 text-muted-foreground italic text-xs">
                                            No significant overlaps found.
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="graph" className="mt-4 border rounded-xl overflow-hidden bg-slate-50 relative h-[600px]">
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        onNodesChange={onNodesChange}
                        onEdgesChange={onEdgesChange}
                        fitView
                    >
                        <Background />
                        <Controls />
                    </ReactFlow>
                    <div className="absolute top-4 right-4 z-10 space-y-2">
                        <Card className="w-48 bg-white/80 backdrop-blur">
                            <CardContent className="p-3 text-[10px] space-y-1">
                                <div className="flex items-center gap-2">
                                    <div className="w-3 h-3 bg-[#f0f9ff] border border-slate-300 rounded" />
                                    <span>Feature</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <div className="w-3 h-3 bg-[#f8fafc] border border-slate-300 rounded" />
                                    <span>Sub-item</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <div className="w-6 h-[1px] bg-[#94a3b8]" />
                                    <span>Dependency</span>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="resolutions" className="mt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>AI Resolution Suggestions</CardTitle>
                            <CardDescription>
                                Automated remediation plans for detected misalignments.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="grid gap-6 md:grid-cols-2">
                                <div className="space-y-4">
                                    {report?.conflicts.map((conflict, idx) => (
                                        <div
                                            key={idx}
                                            className={`p-4 border rounded-xl transition-all ${selectedConflict === conflict ? 'ring-2 ring-indigo-500 bg-indigo-50/20' : 'hover:border-slate-300'}`}
                                            onClick={() => setSelectedConflict(conflict)}
                                        >
                                            <div className="flex items-start gap-3">
                                                {getConflictIcon(conflict.severity)}
                                                <div>
                                                    <h4 className="text-sm font-bold">{conflict.type}</h4>
                                                    <p className="text-xs mt-1 text-slate-600">{conflict.description}</p>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>

                                <div className="bg-slate-50 rounded-xl p-6 border border-dashed flex flex-col">
                                    {selectedConflict ? (
                                        <>
                                            <div className="flex items-center justify-between mb-4">
                                                <Badge className={getSeverityColor(selectedConflict.severity)}>{selectedConflict.severity}</Badge>
                                                <Badge variant="outline">{selectedConflict.type}</Badge>
                                            </div>
                                            <h3 className="text-lg font-bold mb-2">Remediation Plan</h3>
                                            <p className="text-sm text-slate-700 leading-relaxed">
                                                {selectedConflict.remediation}
                                            </p>
                                            <div className="mt-auto pt-6 flex gap-2">
                                                <Button className="flex-1 gap-2">
                                                    Apply Resolution
                                                    <CheckCircle2 className="h-4 w-4" />
                                                </Button>
                                                <Button variant="outline">Ignore</Button>
                                            </div>
                                        </>
                                    ) : (
                                        <div className="m-auto text-center space-y-2">
                                            <Info className="h-10 w-10 text-slate-300 mx-auto" />
                                            <p className="text-sm text-slate-400">Select a conflict to see remediation steps.</p>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
};
