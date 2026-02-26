import { useState, useEffect } from "react";
import { useUpdateRoadmapItem } from "@/hooks/use-roadmap-items";
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
import { Loader2 } from "lucide-react";
import type { components } from "@/api/generated/schema";

interface EditRoadmapItemModalProps {
    projectId: string;
    item: components["schemas"]["RoadmapItem"] | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function EditRoadmapItemModal({ projectId, item, open, onOpenChange }: EditRoadmapItemModalProps) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [status, setStatus] = useState("DRAFT");
    const [error, setError] = useState("");

    const updateItem = useUpdateRoadmapItem(projectId);

    useEffect(() => {
        if (item) {
            setTitle(item.title || "");
            setDescription(item.description || "");
            setStatus(item.status || "DRAFT");
        }
    }, [item]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!item?.id) return;

        try {
            await updateItem.mutateAsync({
                id: item.id,
                updates: {
                    title,
                    description,
                    status: status as any,
                },
            });
            onOpenChange(false);
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to update item.");
        }
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-md">
                <SheetHeader>
                    <SheetTitle>Edit Roadmap Item</SheetTitle>
                    <SheetDescription>
                        Update the title, description, or status of this item.
                    </SheetDescription>
                </SheetHeader>
                <form onSubmit={handleSubmit} className="space-y-6 pt-6">
                    <div className="space-y-2">
                        <Label htmlFor="title">Title</Label>
                        <Input
                            id="title"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea
                            id="description"
                            className="min-h-[100px]"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="status">Status</Label>
                        <Select value={status} onValueChange={setStatus}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="DRAFT">Draft</SelectItem>
                                <SelectItem value="IN_PROGRESS">In Progress</SelectItem>
                                <SelectItem value="COMPLETED">Completed</SelectItem>
                                <SelectItem value="BLOCKED">Blocked</SelectItem>
                                <SelectItem value="DEPRECATED">Deprecated</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {error && (
                        <p className="text-sm font-medium text-destructive">{error}</p>
                    )}

                    <SheetFooter className="pt-4">
                        <Button type="submit" disabled={updateItem.isPending} className="w-full">
                            {updateItem.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Updating...
                                </>
                            ) : (
                                "Update Item"
                            )}
                        </Button>
                    </SheetFooter>
                </form>
            </SheetContent>
        </Sheet>
    );
}
