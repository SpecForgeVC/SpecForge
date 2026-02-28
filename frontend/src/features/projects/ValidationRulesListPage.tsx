import { useState } from "react";
import { useParams } from "react-router-dom";
import { useValidationRules, useDeleteValidationRule } from "@/hooks/use-validation-rules";
import { useProject } from "@/hooks/use-project";
import { useQuery } from "@tanstack/react-query";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Plus, ShieldCheck, Pencil, Trash2, Wand2 } from "lucide-react";
import type { components } from "@/api/generated/schema";
import { CreateValidationRuleModal } from "./components/CreateValidationRuleModal";
import { EditValidationRuleModal } from "./components/EditValidationRuleModal";
import { RecommendationModal } from "@/features/roadmap/components/RecommendationModal";
import { useCreateValidationRule } from "@/hooks/use-validation-rules";
import { contractsApi } from "@/api/contracts";
import { uiRoadmapApi } from "@/api/ui_roadmap";

export function ValidationRulesListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: rules, isLoading: isLoadingRules } = useValidationRules(projectId);

    const { data: contracts = [] } = useQuery({
        queryKey: ["contracts-project", projectId],
        queryFn: () => contractsApi.listContractsByProject(projectId!),
        enabled: !!projectId
    });

    const { data: uiRoadmapItems = [] } = useQuery({
        queryKey: ["ui-roadmap-project", projectId],
        queryFn: () => uiRoadmapApi.list(projectId!),
        enabled: !!projectId
    });

    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedRule, setSelectedRule] = useState<components["schemas"]["ValidationRule"] | null>(null);
    const [isRecommendModalOpen, setIsRecommendModalOpen] = useState(false);

    const isLoading = isLoadingRules;

    const deleteMutation = useDeleteValidationRule(projectId!);

    const handleDelete = async (id: string) => {
        if (window.confirm("Are you sure you want to delete this validation rule?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (rule: components["schemas"]["ValidationRule"]) => {
        setSelectedRule(rule);
        setIsEditModalOpen(true);
    };

    const createRuleMutation = useCreateValidationRule(projectId!);

    const handleApplyRecommendation = async (result: any) => {
        if (result.rules && Array.isArray(result.rules)) {
            for (const rule of result.rules) {
                await createRuleMutation.mutateAsync({
                    name: rule.name,
                    rule_type: rule.rule_type,
                    description: rule.description,
                    rule_config: rule.rule_config
                });
            }
        }
    };

    if (isLoading) return <div className="p-8">Loading rules...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Validation Rules</h1>
                    <p className="text-muted-foreground">
                        Custom validation logic for {project?.name}
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={() => setIsRecommendModalOpen(true)}>
                        <Wand2 className="mr-2 h-4 w-4" /> Recommend Rules
                    </Button>
                    <Button onClick={() => setIsCreateModalOpen(true)}>
                        <Plus className="mr-2 h-4 w-4" /> New Rule
                    </Button>
                </div>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Defined Rules</CardTitle>
                    <CardDescription>Rules applied to variables and contracts</CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Name</TableHead>
                                <TableHead>Type</TableHead>
                                <TableHead>Description</TableHead>
                                <TableHead className="w-[100px] text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {rules?.map((rule) => (
                                <TableRow key={rule.id}>
                                    <TableCell className="font-medium">
                                        <div className="flex items-center">
                                            <ShieldCheck className="mr-2 h-4 w-4 text-muted-foreground" />
                                            {rule.name}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant="secondary">{rule.rule_type}</Badge>
                                    </TableCell>
                                    <TableCell className="text-muted-foreground">{rule.description}</TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-1">
                                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(rule)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => handleDelete(rule.id!)}>
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {(!rules || rules.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={4} className="text-center py-8 text-muted-foreground">
                                        No validation rules found.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            <CreateValidationRuleModal
                projectId={projectId!}
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
            />
            <EditValidationRuleModal
                projectId={projectId!}
                rule={selectedRule}
                isOpen={isEditModalOpen}
                onClose={() => {
                    setSelectedRule(null);
                }}
            />
            <RecommendationModal
                isOpen={isRecommendModalOpen}
                onClose={() => setIsRecommendModalOpen(false)}
                title="Recommend Validation Rules"
                description="AI will analyze your project context to suggest relevant validation rules, security constraints, and data integrity checks."
                targetType="validation_rule"
                contracts={contracts}
                uiRoadmapItems={uiRoadmapItems}
                contextData={{
                    projectName: project?.name,
                    description: project?.description,
                    techStack: project?.tech_stack,
                    rulesCount: rules?.length || 0
                }}
                onApply={handleApplyRecommendation}
            />
        </div>
    );
}
