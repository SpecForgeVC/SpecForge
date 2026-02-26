import { useState } from "react";
import { useCreateWebhook } from "@/hooks/use-webhooks";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle, Webhook } from "lucide-react";

interface CreateWebhookModalProps {
    projectId: string;
    isOpen: boolean;
    onClose: () => void;
}

export function CreateWebhookModal({ projectId, isOpen, onClose }: CreateWebhookModalProps) {
    const [name, setName] = useState("");
    const [url, setUrl] = useState("");
    const [events, setEvents] = useState("");
    const [secret, setSecret] = useState("");
    const [error, setError] = useState<string | null>(null);

    const createMutation = useCreateWebhook(projectId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        try {
            await createMutation.mutateAsync({
                name,
                url,
                events: events.split(",").map(e => e.trim()).filter(e => e !== ""),
                secret,
            });
            onClose();
            setName("");
            setUrl("");
            setEvents("");
            setSecret("");
        } catch (err: any) {
            setError(err.response?.data?.error || err.message || "Failed to create webhook");
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Webhook className="h-5 w-5" />
                        Create New Webhook
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
                            placeholder="For HMAC signature verification"
                            type="password"
                        />
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={onClose}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={createMutation.isPending}>
                            {createMutation.isPending ? "Creating..." : "Create Webhook"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
