import { useCallback, useMemo, useEffect } from 'react';
import ReactFlow, {
    Background,
    Controls,
    MiniMap,
    useNodesState,
    useEdgesState,
    addEdge,
    applyNodeChanges,
    applyEdgeChanges,
    Handle,
    Position,
    type Node,
    type Edge,
    type OnNodesChange,
    type OnEdgesChange
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

const NODE_TYPES = {
    component: ComponentNode,
};

export function ComponentTreeEditor({ data, onChange }: { data: any, onChange?: (data: any) => void }) {
    const nodeTypes = useMemo(() => NODE_TYPES, []);

    // Helper to flatten our tree structure to React Flow nodes/edges
    const generateLayout = useCallback((node: any) => {
        const initialNodes: Node[] = [];
        const initialEdges: Edge[] = [];

        if (node?.nodes && node?.edges) {
            return { nodes: node.nodes, edges: node.edges };
        }

        const flattenTree = (currentNode: any, x = 0, y = 0, parentId?: string) => {
            const id = currentNode.id || Math.random().toString(36).substr(2, 9);
            initialNodes.push({
                id,
                type: 'component',
                position: { x, y },
                data: { ...currentNode },
            });

            if (parentId) {
                initialEdges.push({
                    id: `e-${parentId}-${id}`,
                    source: parentId,
                    target: id,
                    animated: true,
                });
            }

            if (currentNode.children) {
                currentNode.children.forEach((child: any, idx: number) => {
                    const offsetX = (idx - (currentNode.children.length - 1) / 2) * 200;
                    flattenTree(child, x + offsetX, y + 150, id);
                });
            }
        };

        if (node) {
            flattenTree(node);
        }
        return { nodes: initialNodes, edges: initialEdges };
    }, []);

    const layout = useMemo(() => generateLayout(data), [data, generateLayout]);

    const [nodes, setNodes] = useNodesState(layout.nodes);
    const [edges, setEdges] = useEdgesState(layout.edges);

    // Update nodes and edges when the layout changes (e.g. from AI)
    useEffect(() => {
        setNodes(layout.nodes);
        setEdges(layout.edges);
    }, [layout, setNodes, setEdges]);

    const reconstructTree = useCallback((currentNodes: Node[], currentEdges: Edge[]) => {
        // Find root (node with no incoming edges)
        const rootNode = currentNodes.find(n => !currentEdges.some(e => e.target === n.id));
        if (!rootNode) return null;

        const buildNode = (node: Node): any => {
            const children = currentEdges
                .filter(e => e.source === node.id)
                .map(e => currentNodes.find(n => n.id === e.target))
                .filter(Boolean)
                .map(n => buildNode(n!));

            return {
                ...node.data,
                children: children.length > 0 ? children : undefined
            };
        };

        return buildNode(rootNode);
    }, []);

    const onNodesChange: OnNodesChange = useCallback(
        (changes) => {
            setNodes((nds) => {
                const updatedNodes = applyNodeChanges(changes, nds);
                // Execute onChange in a microtask or after render to avoid the "update while rendering" warning
                const fullTree = reconstructTree(updatedNodes, edges);
                setTimeout(() => onChange?.({ ...fullTree, nodes: updatedNodes, edges }), 0);
                return updatedNodes;
            });
        },
        [setNodes, edges, onChange, reconstructTree]
    );

    const onEdgesChange: OnEdgesChange = useCallback(
        (changes) => {
            setEdges((eds) => {
                const updatedEdges = applyEdgeChanges(changes, eds);
                const fullTree = reconstructTree(nodes, updatedEdges);
                setTimeout(() => onChange?.({ ...fullTree, nodes, edges: updatedEdges }), 0);
                return updatedEdges;
            });
        },
        [setEdges, nodes, onChange, reconstructTree]
    );

    const onConnect = useCallback((params: any) => {
        setEdges((eds) => {
            const updatedEdges = addEdge(params, eds);
            const fullTree = reconstructTree(nodes, updatedEdges);
            setTimeout(() => onChange?.({ ...fullTree, nodes, edges: updatedEdges }), 0);
            return updatedEdges;
        });
    }, [setEdges, nodes, onChange, reconstructTree]);

    return (
        <div className="w-full h-[500px] border rounded-xl overflow-hidden bg-muted/20 relative">
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
