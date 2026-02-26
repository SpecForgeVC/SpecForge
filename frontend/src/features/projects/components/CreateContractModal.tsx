import { useState } from "react";
import { useCreateContract } from "@/hooks/use-contracts";
import { useRoadmapItems } from "@/hooks/use-roadmap-items";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import { SchemaEditor } from "@/components/ui/SchemaEditor";

interface CreateContractModalProps {
    projectId: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function CreateContractModal({ projectId, open, onOpenChange }: CreateContractModalProps) {
    const [roadmapItemId, setRoadmapItemId] = useState("");
    const [contractType, setContractType] = useState<string>("REST");
    const [version, setVersion] = useState("1.0.0");
    const [inputSchema, setInputSchema] = useState({});
    const [outputSchema, setOutputSchema] = useState({});
    const [errorSchema, setErrorSchema] = useState({});
    const [error, setError] = useState("");

    const { data: roadmapItems } = useRoadmapItems(projectId);
    const createContract = useCreateContract(projectId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!roadmapItemId) {
            setError("Roadmap Item is required.");
            return;
        }

        try {
            await createContract.mutateAsync({
                roadmap_item_id: roadmapItemId,
                contract_type: contractType as any,
                version,
                input_schema: inputSchema,
                output_schema: outputSchema,
                error_schema: errorSchema,
            });
            onOpenChange(false);
            setRoadmapItemId("");
            setContractType("REST");
            setVersion("1.0.0");
            setInputSchema({});
            setOutputSchema({});
            setErrorSchema({});
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to create contract.");
        }
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>New Contract</SheetTitle>
                    <SheetDescription>
                        Define a new interface contract for a roadmap item.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">
                    <div className="space-y-2">
                        <Label htmlFor="roadmapItem">Roadmap Item</Label>
                        <Select value={roadmapItemId} onValueChange={setRoadmapItemId}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select an item" />
                            </SelectTrigger>
                            <SelectContent>
                                {roadmapItems?.map((item) => (
                                    <SelectItem key={item.id} value={item.id!}>
                                        {item.title} ({item.type})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

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

                    <div className="space-y-2">
                        <Label htmlFor="version">Version</Label>
                        <Input
                            id="version"
                            placeholder="e.g. 1.0.0"
                            value={version}
                            onChange={(e) => setVersion(e.target.value)}
                        />
                    </div>

                    <SchemaEditor
                        label="Input Schema"
                        description="JSON schema for the expected request/input"
                        initialValue={inputSchema}
                        onChange={setInputSchema}
                    />

                    <SchemaEditor
                        label="Output Schema"
                        description="JSON schema for the successful response/output"
                        initialValue={outputSchema}
                        onChange={setOutputSchema}
                    />

                    <SchemaEditor
                        label="Error Schema"
                        description="JSON schema for error responses"
                        initialValue={errorSchema}
                        onChange={setErrorSchema}
                    />

                    {error && (
                        <p className="text-sm font-medium text-destructive">{error}</p>
                    )}

                    <SheetFooter className="pt-4">
                        <Button type="submit" disabled={createContract.isPending} className="w-full">
                            {createContract.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                "Create Contract"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
