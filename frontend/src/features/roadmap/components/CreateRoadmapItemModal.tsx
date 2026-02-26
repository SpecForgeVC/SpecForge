import { useState, useEffect } from "react";
import { useCreateRoadmapItem } from "@/hooks/use-roadmap-items";
import { useRefinement } from "@/hooks/use-refinement";
import { RefinementProgress } from "./RefinementProgress";
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
import { Loader2, Sparkles } from "lucide-react";

interface CreateRoadmapItemModalProps {
    projectId: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function CreateRoadmapItemModal({ projectId, open, onOpenChange }: CreateRoadmapItemModalProps) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [type, setType] = useState<string>("FEATURE");
    const [priority, setPriority] = useState<string>("MEDIUM");
    const [error, setError] = useState("");

    // AI Refinement State
    const [isAIEnabled, setIsAIEnabled] = useState(false);
    const [maxIterations, setMaxIterations] = useState(3);
    const { startSession, session, events, isConnected, reset: resetRefinement } = useRefinement();

    // Auto-fill form when refinement succeeds
    useEffect(() => {
        if (session?.status === 'VALIDATED' && session.result) {
            const result = session.result as Record<string, any>;
            if (result.title) setTitle(String(result.title));
            if (result.description) setDescription(String(result.description));
            if (result.type) setType(String(result.type));
            if (result.priority) setPriority(String(result.priority));
        }
    }, [session?.status, session?.result]);
    const createItem = useCreateRoadmapItem(projectId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (isAIEnabled && !session?.result) {
            // Start Refinement
            if (!description.trim()) {
                setError("Please provide a description/prompt for the AI.");
                return;
            }
            try {
                await startSession("roadmap_item", "roadmap_item", description, { projectId }, maxIterations);
            } catch (err: any) {
                setError(err.message);
            }
            return;
        }

        if (!title.trim()) {
            setError("Title is required.");
            return;
        }

        try {
            const result = session?.result as Record<string, any> || {};
            await createItem.mutateAsync({
                title,
                description,
                type: type as any,
                priority: priority as any,
                status: "DRAFT",
                risk_level: "LOW",
                breaking_change: false,
                regression_sensitive: false,
                business_context: String(result.business_context || ""),
                technical_context: String(result.technical_context || ""),
            });
            onOpenChange(false);
            resetForm();
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to create roadmap item.");
        }
    };

    const resetForm = () => {
        setTitle("");
        setDescription("");
        setType("FEATURE");
        setPriority("MEDIUM");
        setIsAIEnabled(false);
        resetRefinement();
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>New Roadmap Item</SheetTitle>
                    <SheetDescription>
                        Define a new Epic, Feature, or Task for this project.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">

                    <div className="flex items-center justify-between bg-muted/50 p-3 rounded-lg border">
                        <div className="space-y-0.5">
                            <Label htmlFor="ai-mode" className="text-base flex items-center gap-2">
                                <Sparkles className="h-4 w-4 text-purple-500" />
                                AI Assistant
                            </Label>
                            <p className="text-xs text-muted-foreground">
                                Use Multi-Pass Refinement to generate a high-quality spec directly from your description.
                            </p>
                        </div>
                        <Switch
                            id="ai-mode"
                            checked={isAIEnabled}
                            onCheckedChange={(checked) => {
                                setIsAIEnabled(checked);
                                if (!checked) resetRefinement();
                            }}
                        />
                    </div>

                    {isAIEnabled && (
                        <div className="space-y-4 animate-in fade-in zoom-in-95 duration-200">
                            <div className="space-y-2">
                                <Label htmlFor="iterations">Max Iterations</Label>
                                <Select value={String(maxIterations)} onValueChange={(v) => setMaxIterations(Number(v))}>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="1">1 (Fast)</SelectItem>
                                        <SelectItem value="3">3 (Balanced)</SelectItem>
                                        <SelectItem value="5">5 (Thorough)</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            {/* Prompt/Description input is needed first if AI is enabled */}
                        </div>
                    )}

                    {/* Progress View */}
                    {(isConnected || session) && isAIEnabled && (
                        <RefinementProgress events={events} status={session?.status} />
                    )}

                    <div className="space-y-2">
                        <Label htmlFor="title">Title {isAIEnabled && !title && <span className="text-xs text-muted-foreground">(Generated automatically)</span>}</Label>
                        <Input
                            id="title"
                            placeholder="e.g. Implement User Authentication"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            disabled={isAIEnabled && !session?.result && !title} // Only disable if waiting for AI and no title exists
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="type">Type</Label>
                        <Select value={type} onValueChange={setType}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="EPIC">Epic</SelectItem>
                                <SelectItem value="FEATURE">Feature</SelectItem>
                                <SelectItem value="TASK">Task</SelectItem>
                                <SelectItem value="BUGFIX">Bugfix</SelectItem>
                                <SelectItem value="REFACTOR">Refactor</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="priority">Priority</Label>
                        <Select value={priority} onValueChange={setPriority}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="LOW">Low</SelectItem>
                                <SelectItem value="MEDIUM">Medium</SelectItem>
                                <SelectItem value="HIGH">High</SelectItem>
                                <SelectItem value="CRITICAL">Critical</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label htmlFor="description">Description</Label>
                            <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                disabled={!title.trim() || (isConnected && !session?.result)}
                                onClick={async () => {
                                    if (!title.trim()) return;
                                    setIsAIEnabled(true);
                                    try {
                                        await startSession("roadmap_item", "roadmap_item", title, { projectId }, maxIterations);
                                    } catch (err: any) {
                                        setError(err.message);
                                    }
                                }}
                            >
                                <Sparkles className="h-3.5 w-3.5 mr-2 text-purple-500" />
                                Generate Details
                            </Button>
                        </div>
                        <Textarea
                            id="description"
                            placeholder="Detailed description of the item"
                            value={description}
                            onChange={(e: any) => setDescription(e.target.value)}
                            rows={4}
                        />
                    </div>

                    {error && (
                        <p className="text-sm font-medium text-destructive">{error}</p>
                    )}

                    <SheetFooter className="pt-4">
                        <Button type="submit" disabled={createItem.isPending || (isAIEnabled && isConnected && !session?.result)} className="w-full">
                            {createItem.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : isAIEnabled && !session?.result ? (
                                <>
                                    <Sparkles className="mr-2 h-4 w-4" />
                                    start Refinement
                                </>
                            ) : (
                                "Create Item"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
