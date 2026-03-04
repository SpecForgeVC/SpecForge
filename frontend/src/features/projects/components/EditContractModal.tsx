import { useState, useEffect } from "react";
import { useUpdateContract } from "@/hooks/use-contracts";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Loader2, Sparkles } from "lucide-react";
import { SchemaEditor } from "@/components/ui/SchemaEditor";
import type { components } from "@/api/generated/schema";

import { refinementApi } from "@/api/refinement";
import { useRoadmapItem } from "@/hooks/use-roadmap-items";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { useToast } from "@/hooks/use-toast";

interface EditContractModalProps {
    projectId: string;
    contract: components["schemas"]["ContractDefinition"] | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function EditContractModal({ projectId, contract, open, onOpenChange }: EditContractModalProps) {
    const { toast } = useToast();
    const [contractType, setContractType] = useState<string>("REST");
    const [version, setVersion] = useState("");
    const [inputSchema, setInputSchema] = useState({});
    const [outputSchema, setOutputSchema] = useState({});
    const [errorSchema, setErrorSchema] = useState({});
    const [error, setError] = useState("");
    const [generatingField, setGeneratingField] = useState<string | null>(null);
    const [versionBump, setVersionBump] = useState<"patch" | "minor" | "major">("patch");

    // Fetch roadmap item for context
    const { data: roadmapItem } = useRoadmapItem(contract?.roadmap_item_id);

    const updateContract = useUpdateContract(projectId);

    useEffect(() => {
        if (contract) {
            setContractType(contract.contract_type || "REST");
            setVersion(contract.version || "");
            setInputSchema(contract.input_schema || {});
            setOutputSchema(contract.output_schema || {});
            setErrorSchema(contract.error_schema || {});
            setVersionBump("patch"); // Reset bump choice
        }
    }, [contract]);

    const bumpVersion = (current: string, type: "patch" | "minor" | "major"): string => {
        const parts = current.split('.').map(Number);
        if (parts.length !== 3 || parts.some(isNaN)) {
            // Fallback for non-semver
            return current + ".1";
        }
        let [major, minor, patch] = parts;
        if (type === "major") {
            major++;
            minor = 0;
            patch = 0;
        } else if (type === "minor") {
            minor++;
            patch = 0;
        } else {
            patch++;
        }
        return `${major}.${minor}.${patch}`;
    };

    const handleGenerateSchema = async (field: "input_schema" | "output_schema" | "error_schema") => {
        if (!roadmapItem || !contract) return;
        setGeneratingField(field);
        try {
            // Start refinement session
            const prompt = `Generate a JSON Schema for the ${field.replace('_', ' ')}.`;
            const context = {
                title: roadmapItem.title,
                description: roadmapItem.description,
                contract_type: contractType,
                target_field: field,
                existing_input: field !== "input_schema" ? inputSchema : undefined,
                existing_output: field !== "output_schema" ? outputSchema : undefined,
                existing_error: field !== "error_schema" ? errorSchema : undefined,
            };

            const session = await refinementApi.startSession("schema", "schema_suggestion", prompt, context, 5);

            // Use SSE stream for refinement events
            const eventSource = refinementApi.getEvents(session.id!);

            eventSource.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);

                    if (data.type === "SUCCESS" && data.payload?.artifact) {
                        const artifact = data.payload.artifact;
                        // Depending on the field, update the corresponding state
                        if (field === "input_schema" && artifact.input_schema) setInputSchema(artifact.input_schema);
                        if (field === "output_schema" && artifact.output_schema) setOutputSchema(artifact.output_schema);
                        if (field === "error_schema" && artifact.error_schema) setErrorSchema(artifact.error_schema);

                        // Legacy/Alternate support for direct .schema field
                        if (artifact.schema) {
                            if (field === "input_schema") setInputSchema(artifact.schema);
                            if (field === "output_schema") setOutputSchema(artifact.schema);
                            if (field === "error_schema") setErrorSchema(artifact.schema);
                        }

                        toast({
                            title: "Schema Generated",
                            description: `The ${field.replace('_', ' ')} was successfully generated by AI.`,
                            variant: "success",
                        });

                        setGeneratingField(null);
                        eventSource.close();
                    }

                    if (data.type === "ERROR" || data.type === "FAILED") {
                        toast({
                            title: "Generation Failed",
                            description: data.message || "An error occurred during generation.",
                            variant: "destructive",
                        });
                        console.error("Generation Failed:", data.message || "An error occurred during generation.");
                        setGeneratingField(null);
                        eventSource.close();
                    }
                } catch (err) {
                    console.error("Failed to parse refinement event:", err);
                }
            };

            eventSource.onerror = (err) => {
                console.error("SSE Connection Error:", err);
                setGeneratingField(null);
                eventSource.close();
            };
        } catch (e) {
            console.error("Failed to start refinement session:", e);
            setGeneratingField(null);
            toast({
                title: "Error",
                description: "Failed to start AI generation session.",
                variant: "destructive",
            });
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!contract?.id) return;

        try {
            await updateContract.mutateAsync({
                id: contract.id,
                updates: {
                    contract_type: contractType as any,
                    version: bumpVersion(version, versionBump), // Auto-bump
                    input_schema: inputSchema,
                    output_schema: outputSchema,
                    error_schema: errorSchema,
                },
            });
            toast({
                title: "Contract Updated",
                description: "Your changes have been saved successfully.",
                variant: "success",
            });
            onOpenChange(false);
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to update contract.");
            toast({
                title: "Update Failed",
                description: err.response?.data?.error?.message || "An error occurred while saving.",
                variant: "destructive",
            });
        }
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>Edit Contract</SheetTitle>
                    <SheetDescription>
                        Update the interface contract definition.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">
                    <div className="space-y-2">
                        <Label htmlFor="type">Contract Type</Label>
                        <Select value={contractType} onValueChange={setContractType}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="REST">REST API</SelectItem>
                                <SelectItem value="GRAPHQL">GraphQL</SelectItem>
                                <SelectItem value="CLI">CLI Command</SelectItem>
                                <SelectItem value="INTERNAL_FUNCTION">Internal Function</SelectItem>
                                <SelectItem value="EVENT">Event/Message</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-4 border p-4 rounded-md bg-muted/20">
                        <Label>Version Management</Label>
                        <div className="flex items-center gap-4">
                            <div className="text-sm font-mono bg-background px-2 py-1 rounded border">
                                Current: {contract?.version}
                            </div>
                            <div className="text-sm text-muted-foreground">→</div>
                            <div className="text-sm font-mono font-bold text-primary bg-background px-2 py-1 rounded border">
                                New: {bumpVersion(version, versionBump)}
                            </div>
                        </div>
                        <RadioGroup value={versionBump} onValueChange={(v: "patch" | "minor" | "major") => setVersionBump(v)} className="flex gap-4">
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="patch" id="patch" />
                                <Label htmlFor="patch">Patch (Bugfix)</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="minor" id="minor" />
                                <Label htmlFor="minor">Minor (Feature)</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="major" id="major" />
                                <Label htmlFor="major">Major (Breaking)</Label>
                            </div>
                        </RadioGroup>
                    </div>

                    <div className="space-y-2">
                        <div className="flex justify-between items-center">
                            <Label>Input Schema</Label>
                            {Object.keys(inputSchema).length === 0 && (
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-6 text-xs"
                                    onClick={(e) => { e.preventDefault(); handleGenerateSchema("input_schema"); }}
                                    disabled={!!generatingField}
                                >
                                    {generatingField === "input_schema" ? <Loader2 className="h-3 w-3 animate-spin" /> : <Sparkles className="h-3 w-3 mr-1" />}
                                    Auto-generate
                                </Button>
                            )}
                        </div>
                        <SchemaEditor
                            label=""
                            description="JSON schema for the expected request/input"
                            initialValue={inputSchema}
                            onChange={setInputSchema}
                        />
                    </div>

                    <div className="space-y-2">
                        <div className="flex justify-between items-center">
                            <Label>Output Schema</Label>
                            {Object.keys(outputSchema).length === 0 && (
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-6 text-xs"
                                    onClick={(e) => { e.preventDefault(); handleGenerateSchema("output_schema"); }}
                                    disabled={!!generatingField}
                                >
                                    {generatingField === "output_schema" ? <Loader2 className="h-3 w-3 animate-spin" /> : <Sparkles className="h-3 w-3 mr-1" />}
                                    Auto-generate
                                </Button>
                            )}
                        </div>
                        <SchemaEditor
                            label=""
                            description="JSON schema for the successful response/output"
                            initialValue={outputSchema}
                            onChange={setOutputSchema}
                        />
                    </div>

                    <div className="space-y-2">
                        <div className="flex justify-between items-center">
                            <Label>Error Schema</Label>
                            {Object.keys(errorSchema).length === 0 && (
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-6 text-xs"
                                    onClick={(e) => { e.preventDefault(); handleGenerateSchema("error_schema"); }}
                                    disabled={!!generatingField}
                                >
                                    {generatingField === "error_schema" ? <Loader2 className="h-3 w-3 animate-spin" /> : <Sparkles className="h-3 w-3 mr-1" />}
                                    Auto-generate
                                </Button>
                            )}
                        </div>
                        <SchemaEditor
                            label=""
                            description="JSON schema for error responses"
                            initialValue={errorSchema}
                            onChange={setErrorSchema}
                        />
                    </div>

                    {error && (
                        <p className="text-sm font-medium text-destructive">{error}</p>
                    )}

                    <SheetFooter className="pt-4">
                        <Button type="submit" disabled={updateContract.isPending} className="w-full">
                            {updateContract.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Updating...
                                </>
                            ) : (
                                "Update Contract"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
