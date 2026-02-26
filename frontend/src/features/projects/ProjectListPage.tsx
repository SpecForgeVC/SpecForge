import { useParams, Link } from "react-router-dom";
import { useProjects } from "@/hooks/use-projects";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Plus, FolderGit2, ArrowLeft, Trash2 } from "lucide-react";
import { useState } from "react";
import { CreateProjectModal } from "./CreateProjectModal";
import { useDeleteProject } from "@/hooks/use-projects";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export default function ProjectListPage() {
    const { workspaceId } = useParams<{ workspaceId: string }>();
    const { data: projects, isLoading, error } = useProjects(workspaceId);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [projectToDelete, setProjectToDelete] = useState<string | null>(null);
    const deleteProject = useDeleteProject(workspaceId || "");

    const handleDelete = async () => {
        if (projectToDelete) {
            try {
                await deleteProject.mutateAsync(projectToDelete);
                setProjectToDelete(null);
            } catch (err) {
                console.error("Failed to delete project:", err);
            }
        }
    };

    if (isLoading) {
        return <div className="p-8 text-center italic text-muted-foreground">Loading projects...</div>;
    }

    if (error) {
        return (
            <div className="p-8 text-center text-red-600">
                Error loading projects.
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <Button asChild variant="ghost" size="icon">
                        <Link to="/workspaces">
                            <ArrowLeft className="h-4 w-4" />
                        </Link>
                    </Button>
                    <div>
                        <h2 className="text-3xl font-bold tracking-tight">Projects</h2>
                        <p className="text-muted-foreground">Projects in this workspace.</p>
                    </div>
                </div>
                <Button onClick={() => setIsModalOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" /> New Project
                </Button>
            </div>

            {workspaceId && (
                <CreateProjectModal
                    workspaceId={workspaceId}
                    open={isModalOpen}
                    onOpenChange={setIsModalOpen}
                />
            )}

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {projects?.map((project) => (
                    <Card key={project.id} className="group hover:border-primary/50 transition-colors">
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <FolderGit2 className="h-5 w-5 text-primary" />
                                <div className="flex items-center gap-2">
                                    <span className="text-[10px] font-mono text-muted-foreground uppercase">{project.id?.slice(0, 8)}</span>
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        className="h-8 w-8 text-muted-foreground hover:text-destructive transition-colors"
                                        onClick={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            setProjectToDelete(project.id!);
                                        }}
                                    >
                                        <Trash2 className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                            <CardTitle className="mt-2 text-xl">{project.name}</CardTitle>
                            <CardDescription>{project.description || "No description provided."}</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="text-xs text-muted-foreground mb-4">
                                Tech Stack: {project.tech_stack ? Object.keys(project.tech_stack).join(", ") : "Not specified"}
                            </div>
                            <Button asChild variant="secondary" className="w-full">
                                <Link to={`/projects/${project.id}`}>View Dashboard</Link>
                            </Button>
                        </CardContent>
                    </Card>
                ))}
                {projects?.length === 0 && (
                    <div className="col-span-full p-12 text-center border-2 border-dashed rounded-lg">
                        <p className="text-muted-foreground">No projects found in this workspace. Create one to start governing.</p>
                    </div>
                )}
            </div>

            <AlertDialog open={!!projectToDelete} onOpenChange={(open: boolean) => !open && setProjectToDelete(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                        <AlertDialogDescription>
                            This action cannot be undone. This will permanently delete the project
                            and all associated roadmap items, contracts, and snapshots.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={deleteProject.status === 'pending'}>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDelete}
                            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                            disabled={deleteProject.status === 'pending'}
                        >
                            {deleteProject.status === 'pending' ? "Deleting..." : "Delete Project"}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
}
