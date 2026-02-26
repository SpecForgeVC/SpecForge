import { useState } from "react";
import { useFigmaPluginAssets } from "@/hooks/use-ui-roadmap";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
    Download,
    ExternalLink,
    Info,
    CheckCircle2,
    Figma,
    FileCode,
    Terminal,
    ChevronRight,
    Monitor
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export function FigmaPluginInstructions({ id }: { id: string }) {
    const { data: assets, isLoading } = useFigmaPluginAssets(id);
    const [activeStep, setActiveStep] = useState(1);

    const downloadFile = (filename: string, content: string) => {
        const blob = new Blob([content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    const steps = [
        {
            id: 1,
            title: "Download Plugin Assets",
            description: "Get the manifest and logic files generated for this roadmap item.",
            icon: <Download className="h-4 w-4" />
        },
        {
            id: 2,
            title: "Open Figma Desktop",
            description: "Navigate to Plugins > Development > Import plugin from manifest...",
            icon: <Figma className="h-4 w-4" />
        },
        {
            id: 3,
            title: "Select Manifest",
            description: "Pick the 'manifest.json' file you just downloaded.",
            icon: <FileCode className="h-4 w-4" />
        },
        {
            id: 4,
            title: "Sync Layers",
            description: "Run the plugin and select Figma frames to sync to SpecForge.",
            icon: <Monitor className="h-4 w-4" />
        }
    ];

    if (isLoading) return <div className="p-8 text-center animate-pulse italic">Preparing plugin assets...</div>;

    return (
        <div className="space-y-6">
            <Alert className="bg-blue-500/5 border-blue-500/20">
                <Info className="h-4 w-4 text-blue-500" />
                <AlertTitle className="text-blue-700">Live Sync Required</AlertTitle>
                <AlertDescription className="text-blue-600/80 text-xs">
                    To synchronize your Figma designs directly into the Component Tree, you must install the SpecForge developer plugin.
                </AlertDescription>
            </Alert>

            <div className="grid gap-6 md:grid-cols-3">
                {/* Left: Steps */}
                <div className="md:col-span-2 space-y-4">
                    {steps.map((step) => (
                        <div
                            key={step.id}
                            className={`flex gap-4 p-4 rounded-xl border transition-all ${activeStep === step.id ? 'bg-card border-primary ring-1 ring-primary/20 shadow-sm' : 'bg-muted/30 border-transparent opacity-60'
                                }`}
                            onClick={() => setActiveStep(step.id)}
                        >
                            <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm ${activeStep === step.id ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'
                                }`}>
                                {step.id}
                            </div>
                            <div className="space-y-1">
                                <h4 className="font-semibold text-sm flex items-center gap-2">
                                    {step.title}
                                    {activeStep > step.id && <CheckCircle2 className="h-3 w-3 text-emerald-500" />}
                                </h4>
                                <p className="text-xs text-muted-foreground">{step.description}</p>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Right: Actions */}
                <Card className="shadow-sm border-primary/5 h-fit sticky top-4">
                    <CardHeader>
                        <CardTitle className="text-sm">Plugin Package</CardTitle>
                        <CardDescription className="text-[10px]">Version: Alpha 0.1</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-3">
                        <Button
                            variant="outline"
                            size="sm"
                            className="w-full justify-start text-xs"
                            onClick={() => assets && downloadFile('manifest.json', assets['manifest.json'])}
                        >
                            <FileCode className="mr-2 h-4 w-4 text-purple-500" /> manifest.json
                            <Download className="ml-auto h-3 w-3 opacity-50" />
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            className="w-full justify-start text-xs"
                            onClick={() => assets && downloadFile('code.js', assets['code.js'])}
                        >
                            <Terminal className="mr-2 h-4 w-4 text-blue-500" /> code.js
                            <Download className="ml-auto h-3 w-3 opacity-50" />
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            className="w-full justify-start text-xs"
                            onClick={() => assets && downloadFile('ui.html', assets['ui.html'])}
                        >
                            <Monitor className="mr-2 h-4 w-4 text-emerald-500" /> ui.html
                            <Download className="ml-auto h-3 w-3 opacity-50" />
                        </Button>
                        <div className="pt-4 border-t mt-4 space-y-2">
                            <Button className="w-full text-xs font-bold" onClick={() => setActiveStep(2)}>
                                Integrate with Figma <ChevronRight className="ml-2 h-4 w-4" />
                            </Button>
                            <a href="https://www.figma.com/plugin-docs/" target="_blank" rel="noopener noreferrer" className="text-[10px] text-muted-foreground flex items-center justify-center gap-1 hover:text-primary transition-colors">
                                Figma Plugin Docs <ExternalLink className="h-2 w-2" />
                            </a>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
