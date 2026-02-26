import { useState } from "react";
import { useParams } from "react-router-dom";
import { useVariables } from "@/hooks/use-variables";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Plus, Variable, Pencil, Trash2, Network } from "lucide-react";
import { CreateVariableModal } from "./components/CreateVariableModal";
import { EditVariableModal } from "./components/EditVariableModal";
import { useDeleteVariable } from "@/hooks/use-variables";
import type { components } from "@/api/generated/schema";

export function VariableListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: variables, isLoading } = useVariables(projectId);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedVariable, setSelectedVariable] = useState<components["schemas"]["VariableDefinition"] | null>(null);

    const deleteMutation = useDeleteVariable(projectId!);

    const handleDelete = async (id: string) => {
        if (window.confirm("Are you sure you want to delete this variable?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (v: components["schemas"]["VariableDefinition"]) => {
        setSelectedVariable(v);
        setIsEditModalOpen(true);
    };

    if (isLoading) return <div className="p-8">Loading variables...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Variables</h1>
                    <p className="text-muted-foreground">
                        Global variables and constants for {project?.name}
                    </p>
                </div>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Variable
                </Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Defined Variables</CardTitle>
                    <CardDescription>All variables across project contracts</CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Name</TableHead>
                                <TableHead>Type</TableHead>
                                <TableHead>Required</TableHead>
                                <TableHead>Default</TableHead>
                                <TableHead>Description</TableHead>
                                <TableHead className="w-[100px] text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {variables?.map((v) => (
                                <TableRow key={v.id}>
                                    <TableCell className="font-medium">
                                        <div className="flex items-center">
                                            <Variable className="mr-2 h-4 w-4 text-muted-foreground" />
                                            {v.name}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant="secondary">{v.type}</Badge>
                                    </TableCell>
                                    <TableCell>{v.required ? "Yes" : "No"}</TableCell>
                                    <TableCell className="font-mono text-xs">{v.default_value || "-"}</TableCell>
                                    <TableCell className="text-muted-foreground">{v.description}</TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-1">
                                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => window.location.href = `/variables/${v.id}/lineage`}>
                                                <Network className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(v)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => handleDelete(v.id!)}>
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {(!variables || variables.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                                        No variables found for this project.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            <CreateVariableModal
                projectId={projectId!}
                open={isCreateModalOpen}
                onOpenChange={setIsCreateModalOpen}
            />

            <EditVariableModal
                projectId={projectId!}
                variable={selectedVariable}
                open={isEditModalOpen}
                onOpenChange={setIsEditModalOpen}
            />
        </div>
    );
}
