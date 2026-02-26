import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCreateProject } from "@/hooks/use-projects";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Loader2, PlusCircle, Import } from "lucide-react";
import { NewProjectWizard } from "./NewProjectWizard";

interface CreateProjectModalProps {
    workspaceId: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function CreateProjectModal({ workspaceId, open, onOpenChange }: CreateProjectModalProps) {
    const [projectType, setProjectType] = useState<"NEW" | "EXISTING" | null>(null);
    const [name, setName] = useState("");
    const [description, setDescription] = useState("");
    const [repositoryUrl, setRepositoryUrl] = useState("");
    const [enableSelfEvaluation, setEnableSelfEvaluation] = useState(true);
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const createProject = useCreateProject(workspaceId);

    const handleSubmitExisting = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!name.trim()) {
            setError("Project name is required.");
            return;
        }

        try {
            const result = await createProject.mutateAsync({
                name,
                description,
                repository_url: repositoryUrl,
                tech_stack: {},
                settings: {
                    enable_self_evaluation: enableSelfEvaluation,
                },
                project_type: "EXISTING",
            });
            onOpenChange(false);
            navigate(`/projects/${result.data.id}/bootstrap`);
        } catch (err: any) {
            setError(err.response?.data?.error?.message || "Failed to create project.");
        }
    };

    const handleNewProjectComplete = (projectId: string) => {
        onOpenChange(false);
        navigate(`/projects/${projectId}`);
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="sm:max-w-2xl overflow-y-auto">
                <SheetHeader className="mb-6">
                    <SheetTitle>Add Project to Workspace</SheetTitle>
                    <SheetDescription>
                        Start a new governance journey or import an existing codebase.
                    </SheetDescription>
                </SheetHeader>

                {!projectType ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pt-4">
                        <button
                            onClick={() => setProjectType("NEW")}
                            className="flex flex-col items-center justify-center p-8 rounded-xl border-2 border-border bg-card hover:border-violet-500/50 hover:bg-violet-500/5 transition-all group"
                        >
                            <div className="w-12 h-12 rounded-full bg-violet-500/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                                <PlusCircle className="w-6 h-6 text-violet-400" />
                            </div>
                            <h3 className="text-lg font-bold text-foreground mb-2">New Project</h3>
                            <p className="text-sm text-muted-foreground text-center">
                                Use guided AI to define your purpose and tech stack from scratch.
                            </p>
                        </button>

                        <button
                            onClick={() => setProjectType("EXISTING")}
                            className="flex flex-col items-center justify-center p-8 rounded-xl border-2 border-border bg-card hover:border-emerald-500/50 hover:bg-emerald-500/5 transition-all group"
                        >
                            <div className="w-12 h-12 rounded-full bg-emerald-500/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                                <Import className="w-6 h-6 text-emerald-400" />
                            </div>
                            <h3 className="text-lg font-bold text-foreground mb-2">Import Existing</h3>
                            <p className="text-sm text-muted-foreground text-center">
                                Connect an existing repository or analyze your local codebase via IDE.
                            </p>
                        </button>
                    </div>
                ) : projectType === "NEW" ? (
                    <div className="pt-2">
                        <Button
                            variant="ghost"
                            size="sm"
                            className="mb-4 text-muted-foreground hover:text-foreground pl-0"
                            onClick={() => setProjectType(null)}
                        >
                            ← Back to selection
                        </Button>
                        <NewProjectWizard
                            workspaceId={workspaceId}
                            onComplete={handleNewProjectComplete}
                            onCancel={() => setProjectType(null)}
                        />
                    </div>
                ) : (
                    <form onSubmit={handleSubmitExisting} className="space-y-6 pt-2">
                        <Button
                            variant="ghost"
                            size="sm"
                            className="mb-4 text-muted-foreground hover:text-foreground pl-0"
                            onClick={() => setProjectType(null)}
                        >
                            ← Back to selection
                        </Button>
                        <div className="space-y-2">
                            <Label htmlFor="name">Project Name</Label>
                            <Input
                                id="name"
                                placeholder="e.g. Acme API"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="bg-background border-input"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="description">Description</Label>
                            <Input
                                id="description"
                                placeholder="Briefly describe the project's purpose"
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                className="bg-background border-input"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="repo">Repository URL (Optional)</Label>
                            <Input
                                id="repo"
                                placeholder="https://github.com/..."
                                value={repositoryUrl}
                                onChange={(e) => setRepositoryUrl(e.target.value)}
                                className="bg-background border-input"
                            />
                        </div>

                        <div className="flex items-center justify-between rounded-lg border border-border p-4 bg-muted/30">
                            <div className="space-y-0.5">
                                <Label className="text-base">AI Self-Evaluation</Label>
                                <p className="text-sm text-muted-foreground">
                                    Artifacts undergo AI self-critique before persistence.
                                </p>
                            </div>
                            <Switch
                                checked={enableSelfEvaluation}
                                onCheckedChange={setEnableSelfEvaluation}
                            />
                        </div>

                        {error && (
                            <p className="text-sm font-medium text-destructive">{error}</p>
                        )}

                        <Button type="submit" disabled={createProject.isPending} className="w-full bg-emerald-600 hover:bg-emerald-700">
                            {createProject.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                "Next: Bootstrap Intelligence"
                            )}
                        </Button>
                    </form>
                )}
            </SheetContent>
        </Sheet>
    );
}
