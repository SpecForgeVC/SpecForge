import { useState, useEffect } from "react";
import { useUpdateVariable } from "@/hooks/use-variables";
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
import { Textarea } from "@/components/ui/textarea";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Loader2 } from "lucide-react";
import { SchemaEditor } from "@/components/ui/SchemaEditor";
import type { components } from "@/api/generated/schema";

interface EditVariableModalProps {
    projectId: string;
    variable: components["schemas"]["VariableDefinition"] | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function EditVariableModal({ projectId, variable, open, onOpenChange }: EditVariableModalProps) {
    const [name, setName] = useState("");
    const [type, setType] = useState("string");
    const [required, setRequired] = useState(false);
    const [defaultValue, setDefaultValue] = useState("");
    const [description, setDescription] = useState("");
    const [validationRules, setValidationRules] = useState({});
    const [error, setError] = useState("");

    const updateVariable = useUpdateVariable(projectId);

    useEffect(() => {
        if (variable) {
            setName(variable.name || "");
            setType(variable.type || "string");
            setRequired(variable.required || false);
            setDefaultValue(variable.default_value || "");
            setDescription(variable.description || "");
            setValidationRules(variable.validation_rules || {});
        }
    }, [variable]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!variable?.id) return;

        try {
            await updateVariable.mutateAsync({
                id: variable.id,
                updates: {
                    name,
                    type,
                    required,
                    default_value: defaultValue,
                    description,
                    validation_rules: validationRules,
                },
            });
            onOpenChange(false);
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to update variable.");
        }
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>Edit Variable</SheetTitle>
                    <SheetDescription>
                        Update the variable definition.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">
                    <div className="space-y-2">
                        <Label htmlFor="name">Variable Name</Label>
                        <Input
                            id="name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="type">Type</Label>
                            <Select value={type} onValueChange={setType}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="string">String</SelectItem>
                                    <SelectItem value="number">Number</SelectItem>
                                    <SelectItem value="boolean">Boolean</SelectItem>
                                    <SelectItem value="object">Object</SelectItem>
                                    <SelectItem value="array">Array</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="flex flex-col space-y-2 justify-end pb-2">
                            <div className="flex items-center space-x-2">
                                <Switch
                                    id="required"
                                    checked={required}
                                    onCheckedChange={setRequired}
                                />
                                <Label htmlFor="required">Required</Label>
                            </div>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="defaultValue">Default Value</Label>
                        <Input
                            id="defaultValue"
                            value={defaultValue}
                            onChange={(e) => setDefaultValue(e.target.value)}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea
                            id="description"
                            value={description}
                            onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setDescription(e.target.value)}
                        />
                    </div>

                    <SchemaEditor
                        label="Validation Rules"
                        description="JSON configuration for additional validation"
                        initialValue={validationRules}
                        onChange={setValidationRules}
                    />

                    {error && (
                        <p className="text-sm font-medium text-destructive">{error}</p>
                    )}

                    <SheetFooter className="pt-4">
                        <Button type="submit" disabled={updateVariable.isPending} className="w-full">
                            {updateVariable.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Updating...
                                </>
                            ) : (
                                "Update Variable"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
