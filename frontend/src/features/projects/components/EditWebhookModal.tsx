import { useState, useEffect } from "react";
import { useUpdateWebhook } from "@/hooks/use-webhooks";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle, Webhook } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import type { components } from "@/api/generated/schema";

interface EditWebhookModalProps {
    projectId: string;
    webhook: components["schemas"]["Webhook"] | null;
    isOpen: boolean;
    onClose: () => void;
}

export function EditWebhookModal({ projectId, webhook, isOpen, onClose }: EditWebhookModalProps) {
    const [name, setName] = useState("");
    const [url, setUrl] = useState("");
    const [events, setEvents] = useState("");
    const [secret, setSecret] = useState("");
    const [isActive, setIsActive] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const updateMutation = useUpdateWebhook(projectId);

    useEffect(() => {
        if (webhook) {
            setName(webhook.name || "");
            setUrl(webhook.url || "");
            setEvents(webhook.events?.join(", ") || "");
            setIsActive(webhook.active ?? true);
        }
    }, [webhook]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!webhook?.id) return;

        setError(null);
        try {
            await updateMutation.mutateAsync({
                id: webhook.id,
                updates: {
                    name,
                    url,
                    events: events.split(",").map(e => e.trim()).filter(e => e !== ""),
                    secret: secret || undefined,
                    active: isActive,
                },
            });
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.error || err.message || "Failed to update webhook");
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Webhook className="h-5 w-5" />
                        Edit Webhook
                    </DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit} className="space-y-4 py-4">
                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}
                    <div className="space-y-2">
                        <Label htmlFor="name">Name</Label>
                        <Input
                            id="name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="e.g. Production Alerts"
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="url">Payload URL</Label>
                        <Input
                            id="url"
                            value={url}
                            onChange={(e) => setUrl(e.target.value)}
                            placeholder="https://example.com/webhook"
                            required
                            type="url"
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="events">Events (Comma separated)</Label>
                        <Input
                            id="events"
                            value={events}
                            onChange={(e) => setEvents(e.target.value)}
                            placeholder="proposal.created, snapshot.completed"
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="secret">Secret Key (Optional)</Label>
                        <Input
                            id="secret"
                            value={secret}
                            onChange={(e) => setSecret(e.target.value)}
                            placeholder="Leave empty to keep existing"
                            type="password"
                        />
                    </div>
                    <div className="flex items-center space-x-2">
                        <Checkbox
                            id="active"
                            checked={isActive}
                            onCheckedChange={(checked: boolean) => setIsActive(!!checked)}
                        />
                        <Label htmlFor="active" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                            Active
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
