import { useParams } from "react-router-dom";
import { useProjectProposals } from "@/hooks/use-ai-proposals";
import { useApproveProposal, useRejectProposal } from "@/hooks/use-proposal-actions";
import { useProject } from "@/hooks/use-project";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Check, X, Sparkles, Clock, Eye } from "lucide-react";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";

export function ProposalQueuePage() {
    const { projectId } = useParams<{ projectId: string }>();
    const { data: project } = useProject(projectId);
    const { data: proposals, isLoading } = useProjectProposals(projectId);

    if (isLoading) return <div className="p-8">Loading proposals...</div>;

    const pendingProposals = proposals?.filter(p => p.status === "PENDING") || [];
    const processedProposals = proposals?.filter(p => p.status !== "PENDING") || [];

    return (
        <div className="p-8 space-y-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">AI Proposal Queue</h1>
                    <p className="text-muted-foreground">
                        Review and approve AI-generated improvements for {project?.name}
                    </p>
                </div>
                <Badge variant="secondary" className="px-3 py-1">
                    <Sparkles className="mr-2 h-4 w-4 text-primary" />
                    {pendingProposals.length} Pending
                </Badge>
            </div>

            {pendingProposals.length === 0 && (
                <Alert>
                    <Clock className="h-4 w-4" />
                    <AlertTitle>All caught up!</AlertTitle>
                    <AlertDescription>
                        There are no pending AI proposals to review at this time.
                    </AlertDescription>
                </Alert>
            )}

            <div className="grid gap-6">
                {pendingProposals.map((proposal) => (
                    <ProposalCard key={proposal.id} proposal={proposal} />
                ))}
            </div>

            {processedProposals.length > 0 && (
                <div className="pt-8">
                    <h2 className="text-xl font-semibold mb-4 opacity-50">Recent History</h2>
                    <div className="grid gap-4 opacity-75">
                        {processedProposals.map((proposal) => (
                            <ProposalCard key={proposal.id} proposal={proposal} compact />
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}

function ProposalCard({ proposal, compact = false }: { proposal: any, compact?: boolean }) {
    const approveMutation = useApproveProposal(proposal.id!, proposal.roadmap_item_id!);
    const rejectMutation = useRejectProposal(proposal.id!, proposal.roadmap_item_id!);

    const handleApprove = async () => {
        await approveMutation.mutateAsync();
    };

    const handleReject = async () => {
        await rejectMutation.mutateAsync();
    };

    return (
        <Card className={compact ? "bg-slate-50/50" : ""}>
            <CardHeader className={compact ? "py-3" : ""}>
                <div className="flex justify-between items-start">
                    <div className="flex items-center gap-2">
                        <Badge variant="outline" className="capitalize">
                            {proposal.proposal_type?.replace("_", " ")}
                        </Badge>
                        <Badge variant={proposal.status === "PENDING" ? "secondary" : proposal.status === "APPROVED" ? "default" : "destructive"}>
                            {proposal.status}
                        </Badge>
                    </div>
                    {!compact && proposal.confidence_score && (
                        <div className="flex items-center gap-1 text-sm font-medium text-primary">
                            <Sparkles className="h-4 w-4" />
                            {Math.round(proposal.confidence_score * 100)}% Confidence
                        </div>
                    )}
                </div>
                <CardTitle className={compact ? "text-base mt-2" : "mt-2"}>
                    {proposal.reasoning}
                </CardTitle>
                {!compact && (
                    <CardDescription className="mt-2">
                        Item ID: {proposal.roadmap_item_id}
                    </CardDescription>
                )}
            </CardHeader>
            {!compact && (
                <CardContent>
                    <div className="bg-slate-950 rounded-lg p-4 font-mono text-xs text-slate-50 overflow-auto max-h-[300px]">
                        <pre>{JSON.stringify(proposal.diff, null, 2)}</pre>
                    </div>
                </CardContent>
            )}
            {proposal.status === "PENDING" && !compact && (
                <CardFooter className="flex justify-end gap-3 pt-4 border-t">
                    <Button variant="outline" onClick={handleReject} disabled={rejectMutation.isPending}>
                        <X className="mr-2 h-4 w-4" /> Reject
                    </Button>
                    <Button variant="secondary" onClick={() => window.location.href = `/proposals/${proposal.id}/review`}>
                        <Eye className="mr-2 h-4 w-4" /> Review
                    </Button>
                    <Button onClick={handleApprove} disabled={approveMutation.isPending}>
                        <Check className="mr-2 h-4 w-4" /> Approve & Apply
                    </Button>
                </CardFooter>
            )}
        </Card>
    );
}
