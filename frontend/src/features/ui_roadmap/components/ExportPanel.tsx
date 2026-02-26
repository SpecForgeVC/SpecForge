import { useState } from "react";
import { useExportUIRoadmapItem } from "@/hooks/use-ui-roadmap";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
    Download,
    Copy,
    Check,
    FileJson,
    MessageSquare,
    Figma,
    BookOpen,
    Terminal
} from "lucide-react";

export function ExportPanel({ id }: { id: string }) {
    const exportMutation = useExportUIRoadmapItem();
    const [copiedField, setCopiedField] = useState<string | null>(null);

    const handleExport = async () => {
        await exportMutation.mutateAsync(id);
    };

    const copyToClipboard = (text: string, field: string) => {
        navigator.clipboard.writeText(text);
        setCopiedField(field);
        setTimeout(() => setCopiedField(null), 2000);
    };

    const bundle = exportMutation.data;

    return (
        <div className="space-y-6">
            {!bundle ? (
                <Card className="border-dashed flex flex-col items-center justify-center p-12 text-center bg-muted/20">
                    <Terminal className="h-10 w-10 text-muted-foreground/40 mb-4" />
                    <CardTitle className="mb-2">Artifacts Not Generated</CardTitle>
                    <p className="text-muted-foreground mb-6 max-w-xs text-sm">
                        Generate your deterministic export bundle to get LLM prompts, Figma Make instructions, and Storybook scaffolds.
                    </p>
                    <Button onClick={handleExport} disabled={exportMutation.isPending}>
                        {exportMutation.isPending ? "Generating..." : "Generate Export Bundle"}
                    </Button>
                </Card>
            ) : (
                <div className="grid gap-6 md:grid-cols-2">
                    <ExportCard
                        title="LLM Implementation Prompt"
                        description="Strict React/TS/Vite prompt for 100% accurate implementation."
                        icon={<MessageSquare className="h-5 w-5 text-blue-500" />}
                        content={bundle.llm_prompt}
                        onCopy={() => copyToClipboard(bundle.llm_prompt, 'llm')}
                        isCopied={copiedField === 'llm'}
                    />
                    <ExportCard
                        title="Figma Make Prompt"
                        description="Deterministic instructions for Figma AI to generate compliant layouts."
                        icon={<Figma className="h-5 w-5 text-purple-500" />}
                        content={bundle.figma_make}
                        onCopy={() => copyToClipboard(bundle.figma_make, 'figma')}
                        isCopied={copiedField === 'figma'}
                    />
                    <ExportCard
                        title="Storybook Scaffolding"
                        description="Component and state stories with pre-configured controls."
                        icon={<BookOpen className="h-5 w-5 text-pink-500" />}
                        content={bundle.storybook_spec}
                        onCopy={() => copyToClipboard(bundle.storybook_spec, 'storybook')}
                        isCopied={copiedField === 'storybook'}
                    />
                    <Card className="flex flex-col">
                        <CardHeader className="flex flex-row items-center gap-3 border-b py-4">
                            <FileJson className="h-5 w-5 text-emerald-500" />
                            <div>
                                <CardTitle className="text-sm">JSON Specification</CardTitle>
                                <p className="text-[10px] text-muted-foreground">The full deterministic UI contract.</p>
                            </div>
                        </CardHeader>
                        <CardContent className="flex-1 flex flex-col justify-center gap-4 pt-6 text-center">
                            <div className="p-4 bg-muted rounded-lg font-mono text-[10px] text-left max-h-32 overflow-hidden opacity-50">
                                {bundle.json_spec}
                            </div>
                            <Button variant="outline" size="sm" className="w-full" onClick={() => copyToClipboard(bundle.json_spec, 'json')}>
                                {copiedField === 'json' ? <Check className="h-4 w-4 mr-2" /> : <Copy className="h-4 w-4 mr-2" />}
                                Copy JSON Spec
                            </Button>
                            <Button size="sm" className="w-full" disabled>
                                <Download className="h-4 w-4 mr-2" /> Download Full Bundle (.zip)
                            </Button>
                        </CardContent>
                    </Card>
                </div>
            )}
        </div>
    );
}

function ExportCard({ title, description, icon, content, onCopy, isCopied }: any) {
    return (
        <Card className="flex flex-col">
            <CardHeader className="flex flex-row items-center gap-3 border-b py-4">
                <div className="p-2 bg-muted rounded-lg">
                    {icon}
                </div>
                <div>
                    <CardTitle className="text-sm">{title}</CardTitle>
                    <p className="text-[10px] text-muted-foreground">{description}</p>
                </div>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col pt-6 gap-4">
                <div className="flex-1 bg-muted/30 border rounded p-3 relative h-32 overflow-hidden">
                    <pre className="text-[10px] font-mono whitespace-pre-wrap leading-relaxed opacity-60">
                        {content}
                    </pre>
                    <div className="absolute inset-0 bg-gradient-to-t from-background/80 to-transparent pointer-events-none" />
                </div>
                <Button variant="secondary" size="sm" className="w-full font-semibold" onClick={onCopy}>
                    {isCopied ? (
                        <><Check className="h-4 w-4 mr-2" /> Copied!</>
                    ) : (
                        <><Copy className="h-4 w-4 mr-2" /> Copy to Clipboard</>
                    )}
                </Button>
            </CardContent>
        </Card>
    );
}
