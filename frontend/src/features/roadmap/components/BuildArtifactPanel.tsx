import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Download, FileJson, FileText, FolderArchive, Loader2, ShieldCheck, Zap } from "lucide-react";
import { roadmapApi as api } from '@/api/roadmap';
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";

interface BuildArtifactPanelProps {
    roadmapItemId: string;
    completenessScore?: number;
}

export const BuildArtifactPanel: React.FC<BuildArtifactPanelProps> = ({
    roadmapItemId,
    completenessScore = 0
}) => {
    const [format, setFormat] = useState<'json' | 'markdown' | 'zip'>('zip');
    const [includeDeps, setIncludeDeps] = useState(true);
    const [includeGov, setIncludeGov] = useState(true);
    const [isExporting, setIsExporting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleExport = async () => {
        setIsExporting(true);
        setError(null);
        try {
            await api.exportBuildArtifact(roadmapItemId, {
                format,
                include_dependencies: includeDeps,
                include_governance: includeGov
            });
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to export artifacts");
        } finally {
            setIsExporting(false);
        }
    };

    return (
        <Card className="border-primary/20 bg-primary/5">
            <CardHeader>
                <div className="flex items-center justify-between">
                    <div>
                        <CardTitle className="flex items-center gap-2">
                            <Zap className="h-5 w-5 text-primary" />
                            Build Artifact Export Engine
                        </CardTitle>
                        <CardDescription>
                            Generate an agent-ready implementation package for this roadmap item.
                        </CardDescription>
                    </div>
                    <Badge variant={completenessScore > 80 ? "default" : "outline"} className="h-fit">
                        Ready Score: {completenessScore}%
                    </Badge>
                </div>
            </CardHeader>
            <CardContent className="space-y-6">
                {completenessScore < 60 && (
                    <Alert variant="destructive" className="bg-destructive/10 border-destructive/20 text-destructive-foreground">
                        <ShieldCheck className="h-4 w-4" />
                        <AlertTitle>Low Readiness Score</AlertTitle>
                        <AlertDescription>
                            This item has low spec completeness. Autonomous implementation might require more manual guidance.
                        </AlertDescription>
                    </Alert>
                )}

                <div className="grid gap-6 md:grid-cols-2">
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label>Export Format</Label>
                            <Select value={format} onValueChange={(v: any) => setFormat(v)}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Select format" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="zip">
                                        <div className="flex items-center gap-2">
                                            <FolderArchive className="h-4 w-4" />
                                            <span>Zipped Package (.specforge-build)</span>
                                        </div>
                                    </SelectItem>
                                    <SelectItem value="json">
                                        <div className="flex items-center gap-2">
                                            <FileJson className="h-4 w-4" />
                                            <span>Structured JSON (Agent Native)</span>
                                        </div>
                                    </SelectItem>
                                    <SelectItem value="markdown">
                                        <div className="flex items-center gap-2">
                                            <FileText className="h-4 w-4" />
                                            <span>Human Readable Markdown</span>
                                        </div>
                                    </SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-3 pt-2">
                            <div className="flex items-center space-x-2">
                                <Checkbox
                                    id="deps"
                                    checked={includeDeps}
                                    onCheckedChange={(v) => setIncludeDeps(!!v)}
                                />
                                <label htmlFor="deps" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                    Include Roadmap Dependency Graph
                                </label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <Checkbox
                                    id="gov"
                                    checked={includeGov}
                                    onCheckedChange={(v) => setIncludeGov(!!v)}
                                />
                                <label htmlFor="gov" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                    Include Governance & Security Constraints
                                </label>
                            </div>
                        </div>
                    </div>

                    <div className="bg-background/50 rounded-lg p-4 border space-y-3">
                        <h4 className="text-sm font-medium">Included in this export:</h4>
                        <ul className="text-xs space-y-2 text-muted-foreground">
                            <li className="flex items-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-primary" />
                                All linked contracts (REST/Schema)
                            </li>
                            <li className="flex items-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-primary" />
                                All validation rules and variables
                            </li>
                            <li className="flex items-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-primary" />
                                Primary implementation AI prompts
                            </li>
                            <li className="flex items-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-primary" />
                                Verification & Refinement instructions
                            </li>
                            <li className="flex items-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-primary" />
                                Deterministic integrity hash (tamper-proof)
                            </li>
                        </ul>
                    </div>
                </div>

                {error && (
                    <p className="text-sm text-destructive">{error}</p>
                )}

                <Button
                    className="w-full h-12 text-lg font-semibold gap-2 shadow-lg shadow-primary/20"
                    onClick={handleExport}
                    disabled={isExporting}
                >
                    {isExporting ? <Loader2 className="h-5 w-5 animate-spin" /> : <Download className="h-5 w-5" />}
                    {isExporting ? "Generating Package..." : "Export Build Artifacts"}
                </Button>
            </CardContent>
        </Card>
    );
};
