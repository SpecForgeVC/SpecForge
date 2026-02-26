import { useParams, useNavigate } from "react-router-dom";
import { useUIRoadmapItem } from "@/hooks/use-ui-roadmap";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
    Zap,
    ShieldCheck,
    Monitor,
    Download,
    AlertTriangle,
    ChevronLeft,
    Pencil,
    Layout,
    Layers,
    Sparkles,
    Figma
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { ComponentTreeEditor } from "./components/ComponentTreeEditor";
import { StateMachineEditor } from "./components/StateMachineEditor";
import { ExportPanel } from "@/features/ui_roadmap/components/ExportPanel";
import { FigmaPluginInstructions } from "./components/FigmaPluginInstructions";

export default function UIRoadmapItemPage() {
    const { projectId, id } = useParams<{ projectId: string; id: string }>();
    const navigate = useNavigate();
    const { data: item, isLoading } = useUIRoadmapItem(id);

    if (isLoading) return <div className="p-8 text-center italic mt-20">Loading specification details...</div>;
    if (!item) return <div className="p-8 text-center text-destructive">UI Specification not found.</div>;

    return (
        <div className="flex flex-col h-[calc(100vh-4rem)]">
            {/* Header */}
            <div className="border-b bg-muted/30 p-6 flex justify-between items-center">
                <div className="flex items-center gap-4">
                    <Button variant="ghost" size="icon" onClick={() => navigate(`/projects/${projectId}/ui-roadmap`)}>
                        <ChevronLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <div className="flex items-center gap-2">
                            <h1 className="text-2xl font-bold">{item.name}</h1>
                            <Badge variant="outline" className="capitalize">{item.screen_type}</Badge>
                            <Badge className="bg-emerald-500 hover:bg-emerald-600">v{item.version}.0</Badge>
                        </div>
                        <p className="text-sm text-muted-foreground mt-1">{item.description}</p>
                    </div>
                </div>
                <div className="flex items-center gap-3">
                    <Button variant="outline" onClick={() => navigate(`/projects/${projectId}/ui-roadmap/${item.id}/edit`)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit Spec
                    </Button>
                    <Button className="bg-primary hover:bg-primary/90">
                        <Sparkles className="mr-2 h-4 w-4" /> Run Intelligence
                    </Button>
                </div>
            </div>

            <div className="flex-1 flex overflow-hidden">
                {/* Main Content */}
                <div className="flex-1 overflow-y-auto p-8 space-y-8">
                    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
                        <ScoreCard title="Readiness Score" value={item.intelligence_score} icon={<Zap className="h-4 w-4 text-yellow-500" />} />
                        <ScoreCard title="Visual Contract" value={85} icon={<Layout className="h-4 w-4 text-blue-500" />} />
                        <ScoreCard title="A11y Compliance" value={92} icon={<ShieldCheck className="h-4 w-4 text-emerald-500" />} />
                        <ScoreCard title="Drift Risk" value={10} icon={<AlertTriangle className="h-4 w-4 text-rose-500" />} colorInverse />
                    </div>

                    <Tabs defaultValue="visual" className="w-full">
                        <TabsList className="bg-muted/50 p-1">
                            <TabsTrigger value="visual" className="gap-2"><Layers className="h-4 w-4" /> Component Tree</TabsTrigger>
                            <TabsTrigger value="flow" className="gap-2"><Zap className="h-4 w-4" /> State Machine</TabsTrigger>
                            <TabsTrigger value="compliance" className="gap-2"><ShieldCheck className="h-4 w-4" /> Compliance</TabsTrigger>
                            <TabsTrigger value="integrations" className="gap-2"><Figma className="h-4 w-4" /> Integrations</TabsTrigger>
                            <TabsTrigger value="export" className="gap-2"><Download className="h-4 w-4" /> Exports</TabsTrigger>
                        </TabsList>

                        <TabsContent value="visual" className="pt-4 h-[600px]">
                            <ComponentTreeEditor data={item.component_tree} />
                        </TabsContent>

                        <TabsContent value="flow" className="pt-4 h-[600px]">
                            <StateMachineEditor data={item.state_machine} />
                        </TabsContent>

                        <TabsContent value="compliance" className="pt-4 space-y-6">
                            <div className="grid gap-6 md:grid-cols-2">
                                <ComplianceCard
                                    title="Accessibility Spec"
                                    icon={<ShieldCheck className="h-5 w-5 text-emerald-500" />}
                                    details={[
                                        { label: "ARIA Role", value: item.accessibility_spec.role },
                                        { label: "Keyboard Nav", value: item.accessibility_spec.keyboard_nav },
                                        { label: "Focus Management", value: item.accessibility_spec.focus_management },
                                        { label: "Contrast Compliant", value: item.accessibility_spec.contrast_compliance ? "Yes" : "No" }
                                    ]}
                                />
                                <ComplianceCard
                                    title="Responsive Engine"
                                    icon={<Monitor className="h-5 w-5 text-blue-500" />}
                                    details={[
                                        { label: "Mobile Columns", value: item.responsive_spec.mobile.columns },
                                        { label: "Tablet Columns", value: item.responsive_spec.tablet.columns },
                                        { label: "Desktop Columns", value: item.responsive_spec.desktop.columns }
                                    ]}
                                />
                            </div>
                        </TabsContent>

                        <TabsContent value="integrations" className="pt-4">
                            <FigmaPluginInstructions id={item.id} />
                        </TabsContent>

                        <TabsContent value="export" className="pt-4">
                            <ExportPanel id={item.id} />
                        </TabsContent>
                    </Tabs>
                </div>

                {/* Sidebar Info */}
                <div className="w-80 border-l bg-muted/10 p-6 space-y-8 overflow-y-auto">
                    <section className="space-y-4">
                        <h3 className="text-sm font-bold uppercase tracking-widest text-muted-foreground">Context</h3>
                        <div className="space-y-3">
                            <div>
                                <span className="text-xs font-semibold text-muted-foreground block">User Persona</span>
                                <p className="text-sm font-medium">{item.user_persona}</p>
                            </div>
                            <div>
                                <span className="text-xs font-semibold text-muted-foreground block">Primary Use Case</span>
                                <p className="text-sm font-medium">{item.use_case}</p>
                            </div>
                        </div>
                    </section>

                    <section className="space-y-4">
                        <h3 className="text-sm font-bold uppercase tracking-widest text-muted-foreground flex items-center justify-between">
                            Backend Bindings
                            <Badge variant="outline" className="text-[10px]">{Array.isArray(item.backend_bindings) ? item.backend_bindings.length : 0}</Badge>
                        </h3>
                        <div className="space-y-2">
                            {Array.isArray(item.backend_bindings) && item.backend_bindings.map((b: any, i: number) => (
                                <div key={i} className="p-2 border rounded bg-card text-xs flex flex-col gap-1">
                                    <span className="font-mono text-[10px] bg-muted px-1.5 py-0.5 rounded w-fit">{b.method || "GET"}</span>
                                    <span className="font-medium truncate">{b.endpoint}</span>
                                </div>
                            ))}
                            {(!item.backend_bindings || item.backend_bindings.length === 0) && (
                                <p className="text-xs text-muted-foreground italic">No bindings defined</p>
                            )}
                        </div>
                    </section>

                    <section className="p-4 rounded-xl bg-orange-500/5 border border-orange-500/20 space-y-3">
                        <div className="flex items-center gap-2 text-orange-600">
                            <AlertTriangle className="h-4 w-4" />
                            <h3 className="text-xs font-bold uppercase tracking-wider">Drift Warning</h3>
                        </div>
                        <p className="text-[11px] text-orange-700 leading-relaxed">
                            SPEC-UI-004: State machine missing 'partial_success' state required for high-risk data mutation patterns.
                        </p>
                    </section>
                </div>
            </div>
        </div>
    );
}

function ScoreCard({ title, value, icon, colorInverse = false }: any) {
    const displayValue = colorInverse ? 100 - value : value;
    const color = getScoreColor(displayValue);

    return (
        <Card className="p-4 shadow-sm border-primary/5">
            <div className="flex justify-between items-start mb-2">
                <span className="text-xs font-bold text-muted-foreground uppercase tracking-wider">{title}</span>
                {icon}
            </div>
            <div className="space-y-2">
                <div className="text-2xl font-bold">{Math.round(value)}%</div>
                <Progress value={value} className="h-1" style={{ backgroundColor: `${color}20` }} />
            </div>
        </Card>
    );
}

function ComplianceCard({ title, icon, details }: any) {
    return (
        <Card className="shadow-sm">
            <CardHeader className="flex flex-row items-center gap-2 border-b py-4">
                {icon}
                <CardTitle className="text-sm">{title}</CardTitle>
            </CardHeader>
            <CardContent className="pt-4 grid grid-cols-2 gap-4">
                {details.map((d: any, i: number) => (
                    <div key={i}>
                        <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">{d.label}</span>
                        <p className="text-sm font-medium">{d.value}</p>
                    </div>
                ))}
            </CardContent>
        </Card>
    );
}

function getScoreColor(score: number) {
    if (score >= 90) return "#10b981";
    if (score >= 75) return "#3b82f6";
    if (score >= 50) return "#f59e0b";
    return "#ef4444";
}
