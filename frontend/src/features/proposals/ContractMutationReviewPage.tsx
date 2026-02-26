import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { diffLines } from 'diff';
import { Check, X, Info as InfoIcon } from 'lucide-react';
import type { AiProposal } from '@/api/proposals';
import { proposalsApi } from '@/api/proposals';

export const ContractMutationReviewPage: React.FC = () => {
    const { proposalId } = useParams<{ proposalId: string }>();
    const navigate = useNavigate();
    const [proposal, setProposal] = useState<AiProposal | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchProposal = async () => {
            if (!proposalId) return;
            try {
                const data = await proposalsApi.getProposal(proposalId);
                setProposal(data);
            } catch (err) {
                console.error("Failed to fetch proposal", err);
                setError("Failed to load proposal");
            } finally {
                setLoading(false);
            }
        };
        fetchProposal();
    }, [proposalId]);

    const handleApprove = async () => {
        if (!proposal) return;
        try {
            await proposalsApi.approve(proposal.id);
            navigate(-1); // Go back
        } catch (err) {
            console.error("Failed to approve", err);
        }
    };

    const handleReject = async () => {
        if (!proposal) return;
        try {
            await proposalsApi.reject(proposal.id);
            navigate(-1); // Go back
        } catch (err) {
            console.error("Failed to reject", err);
        }
    };

    if (loading) return <div className="p-8 text-center text-muted-foreground">Loading proposal...</div>;
    if (error || !proposal) return <div className="p-8 text-center text-red-600">{error || "Proposal not found"}</div>;

    // Adapt Diff: Assuming diff.original and diff.modified exist, or explicitly passed
    // If not present, we default to empty strings to avoid crash
    const oldSchema = (proposal.diff?.original as string) || '';
    const newSchema = (proposal.diff?.modified as string) || '';

    const diff = diffLines(oldSchema, newSchema);

    return (
        <div className="p-8 max-w-5xl mx-auto space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold">Contract Mutation Proposal</h1>
                    <p className="text-muted-foreground">Review changes proposed by AI Agent</p>
                </div>
                <Badge variant="outline" className="text-lg px-4 py-1">{proposal.status}</Badge>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="flex justify-between items-center">
                        <span>{proposal.proposal_type}</span>
                        <span className="text-sm font-normal text-muted-foreground">Proposal ID: {proposal.id}</span>
                    </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="bg-slate-50 p-4 rounded-md border text-sm">
                        <strong>Reason:</strong> {proposal.reasoning}
                        <br />
                        <strong>Confidence:</strong> {(proposal.confidence_score * 100).toFixed(1)}%
                    </div>

                    <div className="border rounded-md overflow-hidden font-mono text-sm bg-slate-900 text-slate-300">
                        {diff.map((part, index) => {
                            const color = part.added ? 'bg-green-900/30 text-green-400' :
                                part.removed ? 'bg-red-900/30 text-red-400' :
                                    'bg-transparent';
                            const prefix = part.added ? '+' : part.removed ? '-' : ' ';
                            return (
                                <div key={index} className={`${color} px-4 py-1 whitespace-pre-wrap flex`}>
                                    <span className="w-4 select-none opacity-50">{prefix}</span>
                                    <span>{part.value}</span>
                                </div>
                            );
                        })}
                    </div>
                </CardContent>
                {proposal.status === 'PENDING' && (
                    <CardFooter className="flex justify-end gap-3 border-t pt-6">
                        <Button variant="destructive" className="flex gap-2" onClick={handleReject}>
                            <X className="h-4 w-4" /> Reject Proposal
                        </Button>
                        <Button className="flex gap-2 bg-green-600 hover:bg-green-700 text-white" onClick={handleApprove}>
                            <Check className="h-4 w-4" /> Approve Changes
                        </Button>
                    </CardFooter>
                )}
            </Card>

            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 flex items-start gap-3">
                <InfoIcon className="h-5 w-5 text-blue-600 mt-0.5" />
                <div>
                    <h4 className="font-semibold text-blue-900">Impact Analysis</h4>
                    <p className="text-sm text-blue-800 mt-1">
                        Review the diff carefully. Changes here will trigger a new version snapshot upon approval.
                    </p>
                </div>
            </div>
        </div>
    );
};

