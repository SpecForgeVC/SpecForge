import { useCallback, useMemo, useState, useEffect } from 'react';
import ReactFlow, {
    Background,
    Controls,
    MiniMap,
    useNodesState,
    useEdgesState,
    addEdge,
    applyNodeChanges,
    applyEdgeChanges,
    Position,
    Handle,
    MarkerType,
    type Node,
    type Edge,
    type OnNodesChange,
    type OnEdgesChange
} from 'reactflow';
import 'reactflow/dist/style.css';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
    Zap,
    Play,
    RotateCcw,
    ChevronRight,
    Eye,
    Settings,
    Shield,
    Activity
} from 'lucide-react';
import { useStateMachineRunner, type StateMachineDef } from '../hooks/use-state-machine-runner';

const StateNode = ({ data }: { data: any }) => {
    const isActive = data.isActive;

    return (
        <Card className={`p-4 bg-card shadow-lg border-2 min-w-[180px] transition-all duration-300 ${isActive ? 'border-primary ring-2 ring-primary/20 scale-105 shadow-primary/10' : 'border-yellow-500/30'
            }`}>
            <Handle type="target" position={Position.Top} className="w-2 h-2" />
            <div className="flex items-center gap-2 mb-2 border-b border-yellow-500/20 pb-2">
                <Zap className={`h-4 w-4 ${isActive ? 'text-primary fill-primary animate-pulse' : 'text-yellow-500 fill-yellow-500 font-bold'}`} />
                <span className="text-sm font-bold uppercase tracking-tight">{data.label}</span>
                {isActive && <Badge className="ml-auto text-[8px] h-4 bg-primary px-1">Active</Badge>}
            </div>
            <div className="space-y-2 text-[10px]">
                <div>
                    <span className="text-muted-foreground font-semibold block uppercase tracking-tighter">Visual</span>
                    <p className="line-clamp-1 italic">{data.visual_changes || "No mutations defined"}</p>
                </div>
            </div>
            <Handle type="source" position={Position.Bottom} className="w-2 h-2" />
        </Card>
    );
};

const nodeTypes = {
    uiState: StateNode,
};

export function StateMachineEditor({ data, onChange }: { data: any, onChange?: (data: any) => void }) {
    const [isSimulating, setIsSimulating] = useState(false);

    const smDef = useMemo<StateMachineDef>(() => ({
        nodes: Object.entries(data.states || {}).map(([key, config]: [string, any]) => ({
            id: key,
            name: key,
            type: config.type || 'normal',
            visual_mutations: config.visual_changes || {},
            interaction_mutations: config.interaction_changes || {}
        })),
        transitions: (data.transitions || []).map((t: any, i: number) => ({
            id: `t-${i}`,
            source: t.from,
            target: t.to,
            event: t.trigger
        }))
    }), [data]);

    const {
        activeStateId,
        activeState,
        availableTransitions,
        trigger,
        reset
    } = useStateMachineRunner(smDef);

    const initialNodes: Node[] = useMemo(() => {
        if (data.nodes) return data.nodes;
        return Object.entries(data.states || {}).map(([key, config]: [string, any], idx) => ({
            id: key,
            type: 'uiState',
            position: { x: (idx % 3) * 250, y: Math.floor(idx / 3) * 200 },
            data: { ...config, label: key, isActive: key === activeStateId },
        }));
    }, [data, activeStateId]);

    const initialEdges: Edge[] = useMemo(() => {
        if (data.edges) return data.edges;
        return (data.transitions || []).map((t: any, i: number) => ({
            id: `e-${i}`,
            source: t.from,
            target: t.to,
            label: t.trigger,
            animated: isSimulating && t.from === activeStateId,
            style: isSimulating && t.from === activeStateId ? { stroke: 'hsl(var(--primary))', strokeWidth: 2 } : {},
            markerEnd: { type: MarkerType.ArrowClosed, color: isSimulating && t.from === activeStateId ? 'hsl(var(--primary))' : undefined }
        }));
    }, [data, activeStateId, isSimulating]);

    const [nodes, setNodes] = useNodesState(initialNodes);
    const [edges, setEdges] = useEdgesState(initialEdges);

    useEffect(() => {
        setNodes(initialNodes);
        setEdges(initialEdges);
    }, [initialNodes, initialEdges, setNodes, setEdges]);

    const syncCanonicalData = useCallback((currentNodes: Node[], currentEdges: Edge[]) => {
        const newStates = { ...data.states };

        // Update state configs from nodes (keeping metadata like label/id)
        currentNodes.forEach(node => {
            if (newStates[node.id]) {
                newStates[node.id] = {
                    ...newStates[node.id],
                    ...node.data,
                };
            }
        });

        // Convert edges back to transitions
        const newTransitions = currentEdges.map(edge => ({
            from: edge.source,
            to: edge.target,
            trigger: edge.label || ''
        }));

        return {
            ...data,
            states: newStates,
            transitions: newTransitions,
            nodes: currentNodes,
            edges: currentEdges
        };
    }, [data]);

    const onNodesChange: OnNodesChange = useCallback(
        (changes) => {
            setNodes((nds) => {
                const updatedNodes = applyNodeChanges(changes, nds);
                const updatedData = syncCanonicalData(updatedNodes, edges);
                setTimeout(() => onChange?.(updatedData), 0);
                return updatedNodes;
            });
        },
        [setNodes, edges, onChange, syncCanonicalData]
    );

    const onEdgesChange: OnEdgesChange = useCallback(
        (changes) => {
            setEdges((eds) => {
                const updatedEdges = applyEdgeChanges(changes, eds);
                const updatedData = syncCanonicalData(nodes, updatedEdges);
                setTimeout(() => onChange?.(updatedData), 0);
                return updatedEdges;
            });
        },
        [setEdges, nodes, onChange, syncCanonicalData]
    );

    const onConnect = useCallback((params: any) => {
        setEdges((eds) => {
            const updatedEdges = addEdge(params, eds);
            const updatedData = syncCanonicalData(nodes, updatedEdges);
            setTimeout(() => onChange?.(updatedData), 0);
            return updatedEdges;
        });
    }, [setEdges, nodes, onChange, syncCanonicalData]);

    return (
        <div className="w-full h-full min-h-[600px] border rounded-xl overflow-hidden bg-muted/20 flex flex-col relative">
            {/* Toolbar */}
            <div className="absolute top-4 left-4 z-10 flex gap-2">
                <Button
                    size="sm"
                    variant={isSimulating ? "default" : "outline"}
                    className="shadow-md"
                    onClick={() => setIsSimulating(!isSimulating)}
                >
                    {isSimulating ? <Settings className="mr-2 h-4 w-4" /> : <Play className="mr-2 h-4 w-4" />}
                    {isSimulating ? "Edit Mode" : "Simulate Flow"}
                </Button>
                {isSimulating && (
                    <Button size="sm" variant="outline" className="shadow-md bg-background" onClick={reset}>
                        <RotateCcw className="mr-2 h-4 w-4 text-orange-500" /> Reset
                    </Button>
                )}
            </div>

            <div className="flex-1 flex overflow-hidden">
                <div className="flex-1 relative">
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        onNodesChange={onNodesChange}
                        onEdgesChange={onEdgesChange}
                        onConnect={onConnect}
                        nodeTypes={nodeTypes}
                        fitView
                    >
                        <Background />
                        <Controls />
                        <MiniMap />
                    </ReactFlow>
                </div>

                {isSimulating && (
                    <div className="w-80 border-l bg-background p-4 space-y-4 overflow-y-auto animate-in slide-in-from-right duration-300">
                        <section className="space-y-3">
                            <h3 className="text-xs font-bold uppercase tracking-widest text-muted-foreground flex items-center gap-2">
                                <Activity className="h-3 w-3" /> Runner Active
                            </h3>
                            <Card className="bg-primary/5 border-primary/20">
                                <CardHeader className="p-3 pb-0">
                                    <div className="flex items-center justify-between">
                                        <CardTitle className="text-xs text-primary">Current State</CardTitle>
                                        <Badge variant="outline" className="text-[8px] uppercase">{activeStateId}</Badge>
                                    </div>
                                </CardHeader>
                                <CardContent className="p-3 pt-2">
                                    <div className="space-y-2">
                                        <div>
                                            <span className="text-[10px] font-bold opacity-70">Visual Mutations</span>
                                            <div className="text-[11px] font-medium leading-tight mt-1 text-primary/80 italic">
                                                {typeof activeState?.visual_mutations === 'string' ? activeState.visual_mutations : "Default state styling"}
                                            </div>
                                        </div>
                                        {activeState?.interaction_mutations && (
                                            <div>
                                                <span className="text-[10px] font-bold opacity-70">Interaction Restrictions</span>
                                                <div className="text-[11px] font-medium mt-1 p-1 bg-black/5 rounded">
                                                    {typeof activeState.interaction_mutations === 'string' ? activeState.interaction_mutations : "No restrictions"}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                </CardContent>
                            </Card>
                        </section>

                        <section className="space-y-3 pt-4 border-t">
                            <h3 className="text-xs font-bold uppercase tracking-widest text-muted-foreground flex items-center gap-2">
                                <ChevronRight className="h-3 w-3" /> Available Transitions
                            </h3>
                            <div className="space-y-2">
                                {availableTransitions.map((t) => (
                                    <Button
                                        key={t.id}
                                        variant="outline"
                                        size="sm"
                                        className="w-full justify-between h-9 text-xs group hover:border-primary hover:bg-primary/5"
                                        onClick={() => trigger(t.event)}
                                    >
                                        <span className="font-mono">{t.event}</span>
                                        <ChevronRight className="h-3 w-3 opacity-30 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                                    </Button>
                                ))}
                                {availableTransitions.length === 0 && (
                                    <div className="p-4 rounded-lg bg-muted/40 border-dashed border-2 text-center">
                                        <Shield className="h-6 w-6 mx-auto mb-2 text-muted-foreground opacity-20" />
                                        <p className="text-[10px] text-muted-foreground italic">No exit transitions available from this state.</p>
                                    </div>
                                )}
                            </div>
                        </section>

                        <section className="pt-6">
                            <Alert variant="secondary" className="bg-muted/30 border-none shadow-none">
                                <Eye className="h-3 w-3" />
                                <AlertTitle className="text-[10px] font-bold">Observer Mode</AlertTitle>
                                <AlertDescription className="text-[9px]">
                                    Simulating the deterministic state machine allows you to audit UI logic before generation.
                                </AlertDescription>
                            </Alert>
                        </section>
                    </div>
                )}
            </div>
        </div>
    );
}

function Alert({ children, className }: any) {
    return <div className={`p-3 rounded-lg flex gap-3 ${className}`}>{children}</div>;
}

function AlertTitle({ children, className }: any) {
    return <div className={className}>{children}</div>;
}

function AlertDescription({ children, className }: any) {
    return <div className={className}>{children}</div>;
}
