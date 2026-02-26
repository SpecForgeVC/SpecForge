import { useParams } from "react-router-dom";
import { useSnapshots } from "@/hooks/use-snapshots";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Camera, Clock, User } from "lucide-react";
import { format } from "date-fns";

export function SnapshotListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: snapshots, isLoading } = useSnapshots(projectId);

    if (isLoading) return <div className="p-8">Loading snapshots...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Snapshots</h1>
                    <p className="text-muted-foreground">
                        Version history and snapshots for {project?.name}
                    </p>
                </div>
            </div>

            <div className="space-y-4">
                {snapshots?.map((snap) => (
                    <Card key={snap.id}>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0">
                            <div className="space-y-1">
                                <CardTitle className="flex items-center">
                                    <Camera className="mr-2 h-4 w-4 text-primary" />
                                    Snapshot {snap.id?.slice(0, 8)}
                                </CardTitle>
                                <CardDescription className="flex items-center gap-4">
                                    <span className="flex items-center">
                                        <Clock className="mr-1 h-3 w-3" />
                                        {format(new Date(snap.created_at!), "PPP p")}
                                    </span>
                                    <span className="flex items-center">
                                        <User className="mr-1 h-3 w-3" />
                                        {snap.created_by?.slice(0, 8)}
                                    </span>
                                </CardDescription>
                            </div>
                            <Button variant="outline" size="sm">View Details</Button>
                        </CardHeader>
                        <CardContent>
                            <div className="text-sm text-muted-foreground bg-muted p-4 rounded-md font-mono overflow-hidden text-ellipsis whitespace-nowrap">
                                {JSON.stringify(snap.snapshot_data)}
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
