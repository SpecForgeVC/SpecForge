import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
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
import { intelligenceApi } from '../../api/intelligence';
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

const initialNodes: Node[] = [];
const initialEdges: Edge[] = [];

export const VariableLineagePage: React.FC = () => {
    const { variableId } = useParams<{ variableId: string }>();
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
    const [selectedNode, setSelectedNode] = useState<Node | null>(null);

    // Fetch lineage graph from API
    useEffect(() => {
        const fetchGraph = async () => {
            if (!variableId) return;
            try {
                const graph = await intelligenceApi.getLineageGraph(variableId);

                const dagreGraph = new dagre.graphlib.Graph();
                dagreGraph.setGraph({ rankdir: 'LR' });
                dagreGraph.setDefaultEdgeLabel(() => ({}));

                // Transform API nodes to ReactFlow nodes & add to dagre
                const flowNodes: Node[] = graph.nodes.map((n: { id: string; type: string; data: any }) => {
                    const node = {
                        id: n.id,
                        type: n.type === 'variable' ? 'default' : 'input',
                        position: { x: 0, y: 0 }, // Initial position, will be calculated by dagre
                        data: n.data,
                        width: 150, // Approximate width for layout calculation
                        height: 50  // Approximate height
                    };
                    dagreGraph.setNode(n.id, { width: 150, height: 50 });
                    return node;
                });

                // Transform API edges to ReactFlow edges & add to dagre
                const flowEdges: Edge[] = graph.edges.map((e: { id: string; source: string; target: string }) => {
                    dagreGraph.setEdge(e.source, e.target);
                    return {
                        id: e.id,
                        source: e.source,
                        target: e.target,
                        animated: true
                    };
                });

                // Calculate layout
                dagre.layout(dagreGraph);

                // Apply calculated positions
                const layoutedNodes = flowNodes.map((node) => {
                    const nodeWithPosition = dagreGraph.node(node.id);
                    return {
                        ...node,
                        position: {
                            x: nodeWithPosition.x - (node.width! / 2),
                            y: nodeWithPosition.y - (node.height! / 2),
                        },
                    };
                });

                setNodes(layoutedNodes);
                setEdges(flowEdges);
            } catch (err) {
                console.error("Failed to fetch lineage graph", err);
            }
        };
        fetchGraph();
    }, [variableId, setNodes, setEdges]);

    const onNodeClick = (_event: React.MouseEvent, node: Node) => {
        setSelectedNode(node);
    };

    return (
        <div className="h-screen flex flex-col bg-gray-900 text-white">
            <div className="p-4 border-b border-gray-800 flex justify-between items-center">
                <div>
                    <h1 className="text-xl font-bold">Variable Lineage: {variableId}</h1>
                    <p className="text-sm text-gray-400">Visualizing dependencies and impact</p>
                </div>
            </div>
            <div className="flex-1 w-full h-full relative flex">
                <div className="flex-1 h-full">
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        onNodesChange={onNodesChange}
                        onEdgesChange={onEdgesChange}
                        onNodeClick={onNodeClick}
                        fitView
                        className="bg-gray-900"
                    >
                        <Background color="#374151" gap={16} />
                        <Controls className="bg-gray-800 border-gray-700 fill-white" />
                    </ReactFlow>
                </div>

                {/* Details Side Panel */}
                {selectedNode && (
                    <div className="w-80 bg-gray-800 border-l border-gray-700 p-4 overflow-y-auto">
                        <Card className="bg-gray-900 border-gray-700 text-white">
                            <CardHeader>
                                <CardTitle className="text-lg">{selectedNode.data.label || selectedNode.id}</CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div>
                                    <Badge variant="outline" className="text-blue-400 border-blue-400">
                                        {selectedNode.type === 'input' ? 'Source' : 'Variable'}
                                    </Badge>
                                </div>
                                <div className="text-sm text-gray-300">
                                    <strong>Description:</strong>
                                    <p className="mt-1">{selectedNode.data.description || "No description available."}</p>
                                </div>
                                {selectedNode.data.type && (
                                    <div className="text-sm text-gray-300">
                                        <strong>Type:</strong> <span className="font-mono ml-2">{selectedNode.data.type}</span>
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    </div>
                )}
            </div>
        </div>
    );
};
