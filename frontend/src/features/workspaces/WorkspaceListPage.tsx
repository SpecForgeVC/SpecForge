import { useWorkspaces } from "@/hooks/use-workspaces";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Layout } from "lucide-react";
import { Link } from "react-router-dom";
import { useState } from "react";
import { CreateWorkspaceModal } from "./CreateWorkspaceModal";

export default function WorkspaceListPage() {
    const { data: workspaces, isLoading, error } = useWorkspaces();
    const [isModalOpen, setIsModalOpen] = useState(false);

    if (isLoading) {
        return <div className="p-8 text-center italic text-muted-foreground">Loading workspaces...</div>;
    }

    if (error) {
        return (
            <div className="p-8 text-center text-red-600">
                Error loading workspaces. Please check if the backend is running.
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Workspaces</h2>
                    <p className="text-muted-foreground">Manage your spec-first governance environments.</p>
                </div>
                <Button onClick={() => setIsModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Workspace
                </Button>
            </div>

            <CreateWorkspaceModal open={isModalOpen} onOpenChange={setIsModalOpen} />

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {workspaces?.map((workspace) => (
                    <Card key={workspace.id} className="group hover:border-primary/50 transition-colors">
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <Layout className="h-5 w-5 text-primary" />
                                <span className="text-[10px] font-mono text-muted-foreground uppercase">{workspace.id?.slice(0, 8)}</span>
                            </div>
                            <CardTitle className="mt-2 text-xl">{workspace.name}</CardTitle>
                            <CardDescription>{workspace.description || "No description provided."}</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Button asChild variant="secondary" className="w-full">
                                <Link to={`/workspaces/${workspace.id}/projects`}>Open Workspace</Link>
                            </Button>
                        </CardContent>
                    </Card>
                ))}
                {workspaces?.length === 0 && (
                    <div className="col-span-full p-12 text-center border-2 border-dashed rounded-lg">
                        <p className="text-muted-foreground">No workspaces found. Create your first one to get started.</p>
                    </div>
                )}
            </div>
        </div>
    );
}
