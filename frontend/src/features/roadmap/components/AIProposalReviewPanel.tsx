import ReactDiffViewer from "react-diff-viewer-continued";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetFooter
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Check, X, Sparkles, Brain } from "lucide-react";
import type { components } from "@/api/generated/schema";

type AIProposal = components["schemas"]["AIProposal"];

interface AIProposalReviewPanelProps {
    proposal: AIProposal;
    isOpen: boolean;
    onClose: () => void;
    onApprove: (id: string) => void;
    onReject: (id: string) => void;
    currentValue?: string;
}

export function AIProposalReviewPanel({
    proposal,
    isOpen,
    onClose,
    onApprove,
    onReject,
    currentValue = "{}"
}: AIProposalReviewPanelProps) {
    // Mocking the proposed value from the diff object for visualization
    // In a real app, the diff might be a JSON patch or a full new object
    const newValue = JSON.stringify(proposal.diff, null, 2);

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent className="sm:max-w-xl md:max-w-3xl overflow-y-auto">
                <SheetHeader>
                    <div className="flex items-center gap-2 text-primary">
                        <Sparkles className="h-5 w-5" />
                        <SheetTitle>Review AI Proposal</SheetTitle>
                    </div>
                    <SheetDescription>
                        AI suggested a {proposal.proposal_type?.toLowerCase().replace("_", " ")} change.
                    </SheetDescription>
                </SheetHeader>

                <div className="py-6 space-y-6">
                    <div className="p-4 bg-slate-50 border rounded-lg space-y-2">
                        <div className="flex items-center gap-2 font-semibold text-sm">
                            <Brain className="h-4 w-4 text-primary" /> Reasoning
                        </div>
                        <p className="text-sm text-muted-foreground leading-relaxed">
                            {proposal.reasoning || "No detailed reasoning provided."}
                        </p>
                        <div className="flex items-center gap-2 mt-4">
                            <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">Confidence:</span>
                            <div className="flex-1 h-1.5 bg-slate-200 rounded-full overflow-hidden">
                                <div
                                    className="h-full bg-primary"
                                    style={{ width: `${(proposal.confidence_score || 0) * 100}%` }}
                                />
                            </div>
                            <span className="text-xs font-mono">{(proposal.confidence_score || 0) * 100}%</span>
                        </div>
                    </div>

                    <div className="space-y-3">
                        <h3 className="font-semibold text-sm px-1">Visual Change Preview</h3>
                        <div className="border rounded-md overflow-hidden text-xs">
                            <ReactDiffViewer
                                oldValue={currentValue}
                                newValue={newValue}
                                splitView={true}
                                useDarkTheme={false}
                                leftTitle="Current Spec"
                                rightTitle="AI Proposal"
                            />
                        </div>
                    </div>
                </div>

                <SheetFooter className="flex-row gap-2 pt-6 border-t mt-auto">
                    <Button variant="outline" className="flex-1" onClick={() => onReject(proposal.id!)}>
                        <X className="mr-2 h-4 w-4" /> Reject
                    </Button>
                    <Button className="flex-1" onClick={() => onApprove(proposal.id!)}>
                        <Check className="mr-2 h-4 w-4" /> Approve & Apply
                    </Button>
                </SheetFooter>
            </SheetContent>
        </Sheet>
    );
}
