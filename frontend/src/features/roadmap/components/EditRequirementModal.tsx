import { useState, useEffect } from "react";
import { useUpdateRequirement } from "@/hooks/use-requirements";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import type { components } from "@/api/generated/schema";

interface EditRequirementModalProps {
    roadmapItemId: string;
    requirement: components["schemas"]["Requirement"] | null;
    isOpen: boolean;
    onClose: () => void;
}

export function EditRequirementModal({ roadmapItemId, requirement, isOpen, onClose }: EditRequirementModalProps) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [acceptanceCriteria, setAcceptanceCriteria] = useState("");
    const [testable, setTestable] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const updateMutation = useUpdateRequirement(roadmapItemId);

    useEffect(() => {
        if (requirement) {
            setTitle(requirement.title || "");
            setDescription(requirement.description || "");
            setAcceptanceCriteria(requirement.acceptance_criteria || "");
            setTestable(requirement.testable ?? true);
        }
    }, [requirement]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!requirement?.id) return;

        setError(null);
        try {
            await updateMutation.mutateAsync({
                id: requirement.id,
                updates: {
                    title,
                    description,
                    acceptance_criteria: acceptanceCriteria,
                    testable,
                },
            });
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.error || "Failed to update requirement");
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Edit Requirement</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit} className="space-y-4 py-4">
                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}
                    <div className="space-y-2">
                        <Label htmlFor="title">Title</Label>
                        <Input
                            id="title"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            placeholder="e.g. Implement OIDC login"
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="description">Description (Optional)</Label>
                        <Textarea
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Detailed explanation of the requirement..."
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="criteria">Acceptance Criteria</Label>
                        <Textarea
                            id="criteria"
                            value={acceptanceCriteria}
                            onChange={(e) => setAcceptanceCriteria(e.target.value)}
                            placeholder="How do we know this is done?"
                            required
                        />
                    </div>
                    <div className="flex items-center space-x-2">
                        <Checkbox
                            id="testable"
                            checked={testable}
                            onCheckedChange={(checked: boolean) => setTestable(!!checked)}
                        />
                        <Label htmlFor="testable" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                            Is this requirement testable?
                        </Label>
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={onClose}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={updateMutation.isPending}>
                            {updateMutation.isPending ? "Saving..." : "Save Changes"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
