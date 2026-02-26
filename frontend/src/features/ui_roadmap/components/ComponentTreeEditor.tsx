import { useCallback, useMemo } from 'react';
import ReactFlow, {
    Background,
    Controls,
    MiniMap,
    useNodesState,
    useEdgesState,
    addEdge,
    Handle,
    Position,
    type Node,
    type Edge
} from 'reactflow';
import 'reactflow/dist/style.css';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Box, Code } from 'lucide-react';

const ComponentNode = ({ data }: { data: any }) => {
    return (
        <Card className="p-3 bg-card shadow-lg border-2 min-w-[150px] border-primary/20">
            <Handle type="target" position={Position.Top} className="w-2 h-2" />
            <div className="flex items-center gap-2 mb-2 border-b pb-2">
                <Box className="h-4 w-4 text-primary" />
                <span className="text-xs font-bold uppercase tracking-wider">{data.type}</span>
            </div>
            <div className="space-y-1">
                {data.binding && (
                    <div className="text-[10px] flex items-center gap-1 text-muted-foreground">
                        <Code className="h-3 w-3" /> {data.binding}
                    </div>
                )}
                <div className="flex flex-wrap gap-1">
                    {data.validation?.map((v: string) => (
                        <Badge key={v} variant="outline" className="text-[8px] py-0 h-4">
                            {v}
                        </Badge>
                    ))}
                </div>
            </div>
            <Handle type="source" position={Position.Bottom} className="w-2 h-2" />
        </Card>
    );
};

const nodeTypes = {
    component: ComponentNode,
};

export function ComponentTreeEditor({ data }: { data: any }) {
    // Convert our tree structure to React Flow nodes/edges
    const initialNodes: Node[] = [];
    const initialEdges: Edge[] = [];

    const flattenTree = (node: any, x = 0, y = 0, parentId?: string) => {
        const id = node.id || Math.random().toString(36).substr(2, 9);
        initialNodes.push({
            id,
            type: 'component',
            position: { x, y },
            data: { ...node },
        });

        if (parentId) {
            initialEdges.push({
                id: `e-${parentId}-${id}`,
                source: parentId,
                target: id,
                animated: true,
            });
        }

        if (node.children) {
            node.children.forEach((child: any, idx: number) => {
                flattenTree(child, x + (idx - (node.children.length - 1) / 2) * 200, y + 150, id);
            });
        }
    };

    useMemo(() => {
        flattenTree(data);
    }, [data]);

    const [nodes, , onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

    const onConnect = useCallback((params: any) => setEdges((eds) => addEdge(params, eds)), [setEdges]);

    return (
        <div className="w-full h-full min-h-[500px] border rounded-xl overflow-hidden bg-muted/20 relative">
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
            <div className="absolute top-4 right-4 z-10">
                <Badge className="bg-primary hover:bg-primary shadow-sm">
                    Visual Tree Logic Active
                </Badge>
            </div>
        </div>
    );
}
