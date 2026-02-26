import { useState } from "react";
import { useCreateVariable } from "@/hooks/use-variables";
import { useContracts } from "@/hooks/use-contracts";
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

interface CreateVariableModalProps {
    projectId: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function CreateVariableModal({ projectId, open, onOpenChange }: CreateVariableModalProps) {
    const [contractId, setContractId] = useState("");
    const [name, setName] = useState("");
    const [type, setType] = useState("string");
    const [required, setRequired] = useState(false);
    const [defaultValue, setDefaultValue] = useState("");
    const [description, setDescription] = useState("");
    const [validationRules, setValidationRules] = useState({});
    const [error, setError] = useState("");

    const { data: contracts } = useContracts(projectId);
    const createVariable = useCreateVariable(projectId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!contractId) {
            setError("Contract is required.");
            return;
        }
        if (!name.trim()) {
            setError("Name is required.");
            return;
        }

        try {
            await createVariable.mutateAsync({
                contract_id: contractId,
                name,
                type,
                required,
                default_value: defaultValue,
                description,
                validation_rules: validationRules,
            });
            onOpenChange(false);
            setContractId("");
            setName("");
            setType("string");
            setRequired(false);
            setDefaultValue("");
            setDescription("");
            setValidationRules({});
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to create variable.");
        }
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>New Variable</SheetTitle>
                    <SheetDescription>
                        Define a new variable or constant within a contract.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">
                    <div className="space-y-2">
                        <Label htmlFor="contract">Parent Contract</Label>
                        <Select value={contractId} onValueChange={setContractId}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select a contract" />
                            </SelectTrigger>
                            <SelectContent>
                                {contracts?.map((c) => (
                                    <SelectItem key={c.id} value={c.id!}>
                                        {c.contract_type} (v{c.version})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="name">Variable Name</Label>
                        <Input
                            id="name"
                            placeholder="e.g. max_retries"
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
                            placeholder="Optional default value"
                            value={defaultValue}
                            onChange={(e) => setDefaultValue(e.target.value)}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea
                            id="description"
                            placeholder="Purpose of this variable"
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
                        <Button type="submit" disabled={createVariable.isPending} className="w-full">
                            {createVariable.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                "Create Variable"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
