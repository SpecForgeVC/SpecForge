import { useState } from "react";
import { useCreateRequirement } from "@/hooks/use-requirements";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

interface CreateRequirementModalProps {
    roadmapItemId: string;
    isOpen: boolean;
    onClose: () => void;
}

export function CreateRequirementModal({ roadmapItemId, isOpen, onClose }: CreateRequirementModalProps) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [acceptanceCriteria, setAcceptanceCriteria] = useState("");
    const [testable, setTestable] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const createMutation = useCreateRequirement(roadmapItemId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        try {
            await createMutation.mutateAsync({
                title,
                description,
                acceptance_criteria: acceptanceCriteria,
                testable,
            });
            onClose();
            setTitle("");
            setDescription("");
            setAcceptanceCriteria("");
            setTestable(true);
        } catch (err: any) {
            setError(err.response?.data?.error || "Failed to create requirement");
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Create New Requirement</DialogTitle>
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
                        <Button type="submit" disabled={createMutation.isPending}>
                            {createMutation.isPending ? "Creating..." : "Create Requirement"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
