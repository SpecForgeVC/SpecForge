import { useState } from "react";
import { useParams } from "react-router-dom";
import { useContracts } from "@/hooks/use-contracts";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Plus, FileJson, Pencil, Trash2 } from "lucide-react";
import { CreateContractModal } from "./components/CreateContractModal";
import { EditContractModal } from "./components/EditContractModal";
import { useDeleteContract } from "@/hooks/use-contracts";
import type { components } from "@/api/generated/schema";

export function ContractListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: contracts, isLoading } = useContracts(projectId);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedContract, setSelectedContract] = useState<components["schemas"]["ContractDefinition"] | null>(null);

    const deleteMutation = useDeleteContract(projectId!);

    const handleDelete = async (id: string) => {
        if (window.confirm("Are you sure you want to delete this contract?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (contract: components["schemas"]["ContractDefinition"]) => {
        setSelectedContract(contract);
        setIsEditModalOpen(true);
    };

    if (isLoading) return <div className="p-8">Loading contracts...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Contracts</h1>
                    <p className="text-muted-foreground">
                        Interface definitions for {project?.name}
                    </p>
                </div>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Contract
                </Button>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {contracts?.map((contract) => (
                    <Card key={contract.id} className="hover:border-primary/50 transition-colors">
                        <CardHeader>
                            <div className="flex justify-between items-start">
                                <Badge variant="default">
                                    {contract.contract_type}
                                </Badge>
                                <Badge variant="outline">v{contract.version}</Badge>
                                <div className="flex gap-1 ml-auto">
                                    <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(contract)}>
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                    <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => handleDelete(contract.id!)}>
                                        <Trash2 className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                            <CardTitle className="mt-2">
                                Contract {contract.id?.slice(0, 8)}
                            </CardTitle>
                            <CardDescription>
                                Roadmap Item: {contract.roadmap_item_id?.slice(0, 8)}
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center text-sm text-muted-foreground">
                                <FileJson className="mr-2 h-4 w-4" />
                                <span>Compatible: {contract.backward_compatible ? "Yes" : "No"}</span>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            <CreateContractModal
                projectId={projectId!}
                open={isCreateModalOpen}
                onOpenChange={setIsCreateModalOpen}
            />

            <EditContractModal
                projectId={projectId!}
                contract={selectedContract}
                open={isEditModalOpen}
                onOpenChange={setIsEditModalOpen}
            />
        </div>
    );
}
