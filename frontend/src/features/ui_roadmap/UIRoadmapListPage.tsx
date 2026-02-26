import { useParams, useNavigate } from "react-router-dom";
import { useUIRoadmapItems } from "@/hooks/use-ui-roadmap";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Plus, Layout, Zap, BarChart3, Layers } from "lucide-react";
import { Progress } from "@/components/ui/progress";

export function UIRoadmapListPage() {
    const { projectId } = useParams<{ projectId: string }>();
    const navigate = useNavigate();
    const { data: project } = useProject(projectId);
    const { data: items, isLoading } = useUIRoadmapItems(projectId);

    if (isLoading) return <div className="p-8">Loading UI Roadmap...</div>;

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center bg-card p-6 rounded-lg border border-border shadow-sm">
                <div className="flex items-center gap-4">
                    <div className="p-3 bg-primary/10 rounded-full">
                        <Layers className="h-6 w-6 text-primary" />
                    </div>
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">UI Feature Roadmap</h1>
                        <p className="text-muted-foreground flex items-center gap-1.5 mt-1">
                            Deterministic UI governance and drift detection for <span className="font-semibold text-foreground">{project?.name}</span>
                        </p>
                    </div>
                </div>
                <Button onClick={() => navigate(`/projects/${projectId}/ui-roadmap/new`)} className="shadow-sm">
                    <Plus className="mr-2 h-4 w-4" /> New UI Item
                </Button>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {items?.map((item) => (
                    <Card
                        key={item.id}
                        className="group hover:shadow-lg transition-all cursor-pointer border-t-0 border-r-0 border-b-0 border-l-4"
                        style={{ borderLeftColor: getScoreColor(item.intelligence_score) }}
                        onClick={() => navigate(`/projects/${projectId}/ui-roadmap/${item.id}`)}
                    >
                        <CardHeader className="pb-3">
                            <div className="flex justify-between items-start mb-1">
                                <Badge variant="secondary" className="capitalize px-2 py-0 text-xs font-semibold">
                                    {item.screen_type}
                                </Badge>
                                <div className="flex items-center gap-1.5 px-2 py-1 bg-muted rounded-md text-xs font-bold transition-colors group-hover:bg-primary/10">
                                    <Zap className="h-3 w-3 fill-yellow-500 text-yellow-500" />
                                    {Math.round(item.intelligence_score)}%
                                </div>
                            </div>
                            <CardTitle className="text-xl group-hover:text-primary transition-colors">{item.name}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-muted-foreground text-sm line-clamp-2 mb-6 h-10">
                                {item.description || "No description provided."}
                            </p>
                            <div className="space-y-2.5">
                                <div className="flex justify-between items-center text-xs font-medium text-muted-foreground">
                                    <span className="flex items-center gap-1">
                                        <BarChart3 className="h-3 w-3" /> Readiness Score
                                    </span>
                                    <span className="text-foreground">{Math.round(item.intelligence_score)}%</span>
                                </div>
                                <Progress value={item.intelligence_score} className="h-2" />
                            </div>
                            <div className="mt-6 pt-4 border-t border-border flex justify-between items-center">
                                <div className="flex gap-1.5">
                                    <div className="h-2 w-2 rounded-full bg-emerald-500" title="Accessibility Validated" />
                                    <div className="h-2 w-2 rounded-full bg-blue-500" title="State Machine Complete" />
                                    <div className="h-2 w-2 rounded-full bg-purple-500" title="Backend Linked" />
                                </div>
                                <span className="text-[10px] uppercase font-bold text-muted-foreground/60 tracking-wider">
                                    v{item.version}.0
                                </span>
                            </div>
                        </CardContent>
                    </Card>
                ))}
                {items?.length === 0 && (
                    <Card className="col-span-full border-dashed border-2 flex flex-col items-center justify-center p-16 text-center bg-muted/30">
                        <div className="p-4 bg-muted rounded-full mb-6">
                            <Layout className="h-10 w-10 text-muted-foreground/40" />
                        </div>
                        <CardTitle className="text-2xl mb-2">No UI Roadmap Items</CardTitle>
                        <p className="text-muted-foreground mb-8 max-w-md">
                            Start by defining your first deterministic UI feature spec. We'll enforce state machine completeness, design token compliance, and backend contract alignment.
                        </p>
                        <Button variant="default" size="lg" onClick={() => navigate(`/projects/${projectId}/ui-roadmap/new`)}>
                            Create First UI Item
                        </Button>
                    </Card>
                )}
            </div>
        </div>
    );
}

function getScoreColor(score: number) {
    if (score >= 90) return "#10b981"; // emerald-500
    if (score >= 75) return "#3b82f6"; // blue-500
    if (score >= 50) return "#f59e0b"; // amber-500
    return "#ef4444"; // red-500
}
