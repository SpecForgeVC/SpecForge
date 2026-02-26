import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useRoadmapItems } from "@/hooks/use-roadmap-items";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Plus, Target, AlertCircle, Pencil, Trash2 } from "lucide-react";
import { CreateRoadmapItemModal } from "./components/CreateRoadmapItemModal";
import { EditRoadmapItemModal } from "./components/EditRoadmapItemModal";
import { useDeleteRoadmapItem } from "@/hooks/use-roadmap-items";
import type { components } from "@/api/generated/schema";

function getStatusVariant(status: string): "default" | "secondary" | "destructive" | "outline" {
    switch (status) {
        case "COMPLETED": return "default";
        case "IN_PROGRESS": return "secondary";
        case "BLOCKED": return "destructive";
        default: return "outline";
    }
}

export function RoadmapListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const navigate = useNavigate();
    const { data: project } = useProject(projectId);
    const { data: items, isLoading } = useRoadmapItems(projectId);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedItem, setSelectedItem] = useState<components["schemas"]["RoadmapItem"] | null>(null);

    const deleteMutation = useDeleteRoadmapItem(projectId!);

    const handleDelete = async (e: React.MouseEvent, id: string) => {
        e.stopPropagation();
        if (window.confirm("Are you sure you want to delete this roadmap item?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (e: React.MouseEvent, item: components["schemas"]["RoadmapItem"]) => {
        e.stopPropagation();
        setSelectedItem(item);
        setIsEditModalOpen(true);
    };

    if (isLoading) return <div className="p-8">Loading roadmap...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Roadmap Items</h1>
                    <p className="text-muted-foreground">
                        Manage Epic, Feature, and Task items for {project?.name}
                    </p>
                </div>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Item
                </Button>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {items?.map((item) => (
                    <Card
                        key={item.id}
                        className="hover:shadow-md transition-shadow cursor-pointer hover:border-primary/50"
                        onClick={() => navigate(`/roadmap/${item.id}`)}
                    >
                        <CardHeader className="pb-3">
                            <div className="flex justify-between items-start mb-2">
                                <Badge variant={getStatusVariant(item.status!)} className="capitalize">
                                    {item.status?.toLowerCase().replace("_", " ")}
                                </Badge>
                                <Badge variant="secondary" className="capitalize">
                                    {item.type?.toLowerCase()}
                                </Badge>
                                <div className="flex gap-1 ml-auto">
                                    <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground" onClick={(e) => handleEdit(e, item)}>
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                    <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive hover:bg-destructive/10" onClick={(e) => handleDelete(e, item.id!)}>
                                        <Trash2 className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                            <CardTitle className="text-xl leading-snug">{item.title}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-muted-foreground text-sm line-clamp-3 mb-4">
                                {item.description}
                            </p>
                            <div className="flex items-center gap-4 text-xs font-medium text-muted-foreground uppercase tracking-wider">
                                <span className="flex items-center gap-1.5">
                                    <Target className="h-3.5 w-3.5" />
                                    {item.priority}
                                </span>
                                {item.breaking_change && (
                                    <span className="flex items-center gap-1.5 text-destructive">
                                        <AlertCircle className="h-3.5 w-3.5" />
                                        Breaking
                                    </span>
                                )}
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            <CreateRoadmapItemModal
                projectId={projectId!}
                open={isCreateModalOpen}
                onOpenChange={setIsCreateModalOpen}
            />

            <EditRoadmapItemModal
                projectId={projectId!}
                item={selectedItem}
                open={isEditModalOpen}
                onOpenChange={setIsEditModalOpen}
            />
        </div>
    );
}
