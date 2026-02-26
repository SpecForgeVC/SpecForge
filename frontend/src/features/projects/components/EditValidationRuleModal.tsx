import { useState, useEffect } from "react";
import { useUpdateValidationRule } from "@/hooks/use-validation-rules";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { components } from "@/api/generated/schema";

interface EditValidationRuleModalProps {
    projectId: string;
    rule: components["schemas"]["ValidationRule"] | null;
    isOpen: boolean;
    onClose: () => void;
}

export function EditValidationRuleModal({ projectId, rule, isOpen, onClose }: EditValidationRuleModalProps) {
    const [name, setName] = useState("");
    const [description, setDescription] = useState("");
    const [ruleType, setRuleType] = useState("SCHEMA");
    const [ruleDefinition, setRuleDefinition] = useState("{}");
    const [error, setError] = useState<string | null>(null);

    const updateMutation = useUpdateValidationRule(projectId);

    useEffect(() => {
        if (rule) {
            setName(rule.name || "");
            setDescription(rule.description || "");
            setRuleType(rule.rule_type || "SCHEMA");
            setRuleDefinition(JSON.stringify(rule.rule_config, null, 2) || "{}");
        }
    }, [rule]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!rule?.id) return;

        setError(null);
        try {
            await updateMutation.mutateAsync({
                id: rule.id,
                updates: {
                    name,
                    description,
                    rule_type: ruleType,
                    rule_config: JSON.parse(ruleDefinition),
                },
            });
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.error || err.message || "Failed to update rule");
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Edit Validation Rule</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit} className="space-y-4 py-4">
                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}
                    <div className="space-y-2">
                        <Label htmlFor="name">Rule Name</Label>
                        <Input
                            id="name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="e.g. Email Format"
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="description">Description (Optional)</Label>
                        <Textarea
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="What does this rule validate?"
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="type">Rule Type</Label>
                        <Select value={ruleType} onValueChange={setRuleType}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select type" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="SCHEMA">Schema Validation</SelectItem>
                                <SelectItem value="REGEX">Regex Match</SelectItem>
                                <SelectItem value="CUSTOM">Custom Logic</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="definition">Rule Definition (JSON)</Label>
                        <Textarea
                            id="definition"
                            value={ruleDefinition}
                            onChange={(e) => setRuleDefinition(e.target.value)}
                            className="font-mono text-xs h-[150px]"
                            placeholder='{"type": "string", "format": "email"}'
                            required
                        />
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
