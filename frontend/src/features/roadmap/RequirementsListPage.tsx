import { useState } from "react";
import { useParams } from "react-router-dom";
import { useRoadmapItems } from "@/hooks/use-roadmap-items";
import { useRequirements, useDeleteRequirement, useCreateRequirement } from "@/hooks/use-requirements";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { FileText, Pencil, Trash2, ChevronRight, ChevronDown, Plus, Sparkles } from "lucide-react";
import type { components } from "@/api/generated/schema";
import { CreateRequirementModal } from "./components/CreateRequirementModal";
import { EditRequirementModal } from "./components/EditRequirementModal";
import { RecommendationModal } from "./components/RecommendationModal";
import { variablesApi } from "@/api/variables";
import { contractsApi } from "@/api/contracts";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export function RequirementsListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: roadmapItems, isLoading: loadingRoadmap } = useRoadmapItems(projectId);

    if (loadingRoadmap) return <div className="p-8">Loading requirements...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Technical Requirements</h1>
                    <p className="text-muted-foreground">
                        Specific implementation requirements for {project?.name}
                    </p>
                </div>
            </div>

            <div className="space-y-6">
                {roadmapItems?.map((item) => (
                    <RoadmapItemRequirements key={item.id} item={item} />
                ))}
                {(!roadmapItems || roadmapItems.length === 0) && (
                    <Card className="bg-slate-50/50 border-dashed">
                        <CardContent className="flex flex-col items-center justify-center py-12 text-muted-foreground">
                            <FileText className="h-12 w-12 mb-4 opacity-20" />
                            <p>No roadmap items found. Create one to add requirements.</p>
                        </CardContent>
                    </Card>
                )}
            </div>
        </div>
    );
}

function RoadmapItemRequirements({ item }: { item: components["schemas"]["RoadmapItem"] }) {
    const [isExpanded, setIsExpanded] = useState(true);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedReq, setSelectedReq] = useState<components["schemas"]["Requirement"] | null>(null);

    const { data: requirements } = useRequirements(item.id!);
    const createRequirement = useCreateRequirement(item.id!);
    const deleteMutation = useDeleteRequirement(item.id!);
    const queryClient = useQueryClient();

    const [recModal, setRecModal] = useState<{
        isOpen: boolean;
        title: string;
        description: string;
        targetType: string;
    }>({
        isOpen: false,
        title: "",
        description: "",
        targetType: "requirement"
    });

    const { data: contracts = [] } = useQuery({
        queryKey: ["contracts", item.id],
        queryFn: () => contractsApi.listContracts(item.id!),
        enabled: !!item.id && recModal.isOpen
    });

    const handleDelete = async (id: string) => {
        if (window.confirm("Are you sure you want to delete this requirement?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (req: components["schemas"]["Requirement"]) => {
        setSelectedReq(req);
        setIsEditModalOpen(true);
    };

    return (
        <Card className="overflow-hidden">
            <CardHeader className="bg-slate-50/50 cursor-pointer py-4" onClick={() => setIsExpanded(!isExpanded)}>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                        <CardTitle className="text-lg">{item.title}</CardTitle>
                        <Badge variant="outline">{item.type}</Badge>
                    </div>
                    <div className="flex gap-2">
                        <Button
                            size="sm"
                            variant="secondary"
                            onClick={(e) => {
                                e.stopPropagation();
                                setIsCreateModalOpen(true);
                            }}
                        >
                            <Plus className="h-4 w-4 mr-2" />
                            Add Requirement
                        </Button>
                        <Button
                            size="sm"
                            variant="outline"
                            onClick={(e) => {
                                e.stopPropagation();
                                setRecModal({
                                    isOpen: true,
                                    title: "Recommend Requirements",
                                    description: "Generate technical requirements based on the roadmap item details.",
                                    targetType: "requirement"
                                });
                            }}
                        >
                            <Sparkles className="h-4 w-4 mr-2" />
                            Recommend
                        </Button>
                    </div>
                </div>
            </CardHeader>
            {isExpanded && (
                <CardContent className="p-0">
                    <Table>
                        <TableHeader>
                            <TableRow className="hover:bg-transparent">
                                <TableHead className="pl-6">Title</TableHead>
                                <TableHead>Testable</TableHead>
                                <TableHead>Criteria</TableHead>
                                <TableHead className="w-[100px] text-right pr-6">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {requirements?.map((req) => (
                                <TableRow key={req.id}>
                                    <TableCell className="font-medium pl-6">{req.title}</TableCell>
                                    <TableCell>
                                        <Badge variant={req.testable ? "default" : "outline"}>
                                            {req.testable ? "Testable" : "Descriptive"}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-muted-foreground max-w-md truncate">
                                        {req.acceptance_criteria}
                                    </TableCell>
                                    <TableCell className="text-right pr-6">
                                        <div className="flex justify-end gap-1">
                                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(req)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => handleDelete(req.id!)}>
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {(!requirements || requirements.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={4} className="text-center py-6 text-muted-foreground italic">
                                        No requirements defined for this item.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            )}

            <CreateRequirementModal
                roadmapItemId={item.id!}
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
            />
            <EditRequirementModal
                roadmapItemId={item.id!}
                requirement={selectedReq}
                isOpen={isEditModalOpen}
                onClose={() => {
                    setIsEditModalOpen(false);
                    setSelectedReq(null);
                }}
            />
            <RecommendationModal
                isOpen={recModal.isOpen}
                onClose={() => setRecModal({ ...recModal, isOpen: false })}
                title={recModal.title}
                description={recModal.description}
                targetType={recModal.targetType}
                contextData={item}
                contracts={contracts}
                onApply={async (result) => {
                    try {
                        let reqCount = 0;
                        let varCount = 0;

                        // Handle Requirements
                        if (result.requirements && Array.isArray(result.requirements)) {
                            for (const req of result.requirements) {
                                await createRequirement.mutateAsync({
                                    title: req.title,
                                    acceptance_criteria: req.acceptance_criteria + (req.priority ? `\n\n[Priority: ${req.priority}]` : ""),
                                    testable: req.testable,
                                    order_index: 0
                                });
                                reqCount++;
                            }
                        }

                        // Handle Variables
                        if (result.variables && Array.isArray(result.variables) && item.project_id) {
                            if (contracts.length === 0) {
                                console.warn("Skipping variable creation: No contracts found.");
                                alert(`Created ${reqCount} requirements. Variables were skipped because no API Contracts exist for this item.`);
                            } else {
                                const targetContract = contracts[0];
                                for (const v of result.variables) {
                                    try {
                                        await variablesApi.createVariable(item.project_id, {
                                            ...v,
                                            project_id: item.project_id,
                                            contract_id: targetContract.id
                                        });
                                        varCount++;
                                    } catch (e) {
                                        console.error("Failed to create variable", v.name, e);
                                    }
                                }
                                await queryClient.invalidateQueries({ queryKey: ["variables", item.project_id] });
                            }
                        }

                        setRecModal({ ...recModal, isOpen: false });
                        if (varCount > 0) {
                            alert(`Successfully created ${reqCount} requirements and ${varCount} variables (attached to ${contracts[0]?.contract_type}).`);
                        } else if (contracts.length > 0 || !result.variables || result.variables.length === 0) {
                            if (reqCount > 0) alert(`Successfully created ${reqCount} requirements.`);
                        }
                    } catch (err) {
                        console.error("Failed to apply recommendations:", err);
                        alert("Failed to apply changes. See console.");
                    }
                }}
            />
        </Card>
    );
}
