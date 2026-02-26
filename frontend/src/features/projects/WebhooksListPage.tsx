import { useState } from "react";
import { useParams } from "react-router-dom";
import { useWebhooks, useDeleteWebhook } from "@/hooks/use-webhooks";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Plus, Webhook as WebhookIcon, Pencil, Trash2 } from "lucide-react";
import type { components } from "@/api/generated/schema";
import { CreateWebhookModal } from "./components/CreateWebhookModal";
import { EditWebhookModal } from "./components/EditWebhookModal";

export function WebhooksListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: webhooks, isLoading } = useWebhooks(projectId);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [selectedWebhook, setSelectedWebhook] = useState<components["schemas"]["Webhook"] | null>(null);

    const deleteMutation = useDeleteWebhook(projectId!);

    const handleDelete = async (id: string) => {
        if (window.confirm("Are you sure you want to delete this webhook?")) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleEdit = (webhook: components["schemas"]["Webhook"]) => {
        setSelectedWebhook(webhook);
        setIsEditModalOpen(true);
    };

    if (isLoading) return <div className="p-8">Loading webhooks...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Webhooks</h1>
                    <p className="text-muted-foreground">
                        External notifications for {project?.name}
                    </p>
                </div>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Webhook
                </Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Configured Webhooks</CardTitle>
                    <CardDescription>HTTP endpoints for system events</CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Name</TableHead>
                                <TableHead>URL</TableHead>
                                <TableHead>Status</TableHead>
                                <TableHead>Events</TableHead>
                                <TableHead className="w-[100px] text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {webhooks?.map((webhook) => (
                                <TableRow key={webhook.id}>
                                    <TableCell className="font-medium">
                                        <div className="flex items-center">
                                            <WebhookIcon className="mr-2 h-4 w-4 text-muted-foreground" />
                                            {webhook.name || "Unnamed Webhook"}
                                        </div>
                                    </TableCell>
                                    <TableCell className="max-w-[200px] truncate underline decoration-dotted">
                                        {webhook.url}
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant={webhook.active ? "default" : "outline"}>
                                            {webhook.active ? "Active" : "Disabled"}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex flex-wrap gap-1">
                                            {webhook.events?.map((event) => (
                                                <Badge key={event} variant="secondary" className="text-[10px]">
                                                    {event}
                                                </Badge>
                                            ))}
                                        </div>
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <div className="flex justify-end gap-1">
                                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(webhook)}>
                                                <Pencil className="h-4 w-4" />
                                            </Button>
                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => handleDelete(webhook.id!)}>
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {(!webhooks || webhooks.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                                        No webhooks configured.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            <CreateWebhookModal
                projectId={projectId!}
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
            />
            <EditWebhookModal
                projectId={projectId!}
                webhook={selectedWebhook}
                isOpen={isEditModalOpen}
                onClose={() => {
                    setIsEditModalOpen(false);
                    setSelectedWebhook(null);
                }}
            />
        </div>
    );
}
