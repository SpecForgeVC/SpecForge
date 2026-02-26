import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { GitCommit, AlertTriangle, ArrowRight } from 'lucide-react';

interface DriftRecord {
    id: string;
    timestamp: string;
    feature: string;
    contract: string;
    driftType: 'SCHEMA_BREAKING' | 'SCHEMA_ADDITION' | 'SEMANTIC';
    severity: 'LOW' | 'MEDIUM' | 'HIGH';
    details: string;
}

import { intelligenceApi } from '@/api/intelligence';
// Removed unused AuditLog import

export const DriftHistoryPage: React.FC = () => {
    const [history, setHistory] = useState<DriftRecord[]>([]);
    // removed unused loading, error states for now as they were causing lints
    // const [loading, setLoading] = useState(true);
    // const [error, setError] = useState<string | null>(null);

    React.useEffect(() => {
        const fetchHistory = async () => {
            try {
                const logs = await intelligenceApi.getDriftHistory();
                if (!Array.isArray(logs)) {
                    console.warn('Drift history API returned non-array:', logs);
                    setHistory([]);
                    return;
                }
                const records: DriftRecord[] = logs.map(log => {
                    const report = log.new_data.drift_report;
                    // Improved typing for drift classification
                    const isBreaking = report?.breaking_changes && report.breaking_changes.length > 0;
                    const riskScore = report?.risk_score || 0;

                    const driftType = isBreaking ? 'SCHEMA_BREAKING' :
                        riskScore > 0.5 ? 'SEMANTIC' : 'SCHEMA_ADDITION';

                    const severity = riskScore > 0.7 ? 'HIGH' :
                        riskScore > 0.3 ? 'MEDIUM' : 'LOW';

                    return {
                        id: log.id,
                        timestamp: new Date(log.created_at).toLocaleString(),
                        feature: log.entity_type === 'CONTRACT' ? `Contract ${log.entity_id}` : 'Unknown Feature',
                        contract: log.entity_id,
                        driftType: driftType,
                        severity: severity,
                        details: report?.breaking_changes?.[0]?.issue || 'Drift detected'
                    };
                });
                setHistory(records);
            } catch (err) {
                console.error('Failed to fetch drift history', err);
                // setError('Failed to load drift history');
            } finally {
                // setLoading(false);
            }
        };

        fetchHistory();
    }, []);

    const getSeverityColor = (severity: string) => {
        switch (severity) {
            case 'HIGH': return 'bg-red-100 text-red-800 border-red-200';
            case 'MEDIUM': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
            case 'LOW': return 'bg-blue-100 text-blue-800 border-blue-200';
            default: return 'bg-gray-100 text-gray-800';
        }
    };

    const [selectedRecord, setSelectedRecord] = useState<DriftRecord | null>(null);

    // ... (fetch logic same as before) ...

    return (
        <div className="p-8 max-w-6xl mx-auto">
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Drift History</h1>
                    <p className="text-muted-foreground">Track divergence between Specs and Implementation</p>
                </div>
                <Button variant="outline">
                    <GitCommit className="mr-2 h-4 w-4" />
                    Force Drift Check
                </Button>
            </div>

            <div className="grid gap-4">
                {history.map((record) => (
                    <Card key={record.id} className="hover:shadow-md transition-shadow">
                        <CardHeader className="pb-2">
                            <div className="flex justify-between items-start">
                                <div className="space-y-1">
                                    <CardTitle className="text-lg font-medium flex items-center gap-2">
                                        {record.contract}
                                        <span className="text-sm font-normal text-muted-foreground">in {record.feature}</span>
                                    </CardTitle>
                                    <div className="text-sm text-muted-foreground">{record.timestamp}</div>
                                </div>
                                <Badge variant="outline" className={getSeverityColor(record.severity)}>
                                    {record.driftType.replace('_', ' ')}
                                </Badge>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <p className="text-sm">{record.details}</p>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="gap-1"
                                    onClick={() => setSelectedRecord(record)}
                                >
                                    View Details <ArrowRight className="h-4 w-4" />
                                </Button>
                            </div>
                            {record.severity === 'HIGH' && (
                                <div className="mt-3 flex items-center gap-2 text-xs text-red-600 bg-red-50 p-2 rounded">
                                    <AlertTriangle className="h-3 w-3" />
                                    <span>Immediate attention required: Breaking change detected without version increment.</span>
                                </div>
                            )}
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Diff Details Dialog */}
            {selectedRecord && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={() => setSelectedRecord(null)}>
                    <div className="bg-white p-6 rounded-lg max-w-2xl w-full m-4 shadow-xl" onClick={e => e.stopPropagation()}>
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-bold">Drift Details</h3>
                            <Button variant="ghost" size="sm" onClick={() => setSelectedRecord(null)}>Close</Button>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <h4 className="font-semibold text-sm text-gray-500">Drift Type</h4>
                                <p>{selectedRecord.driftType}</p>
                            </div>
                            <div>
                                <h4 className="font-semibold text-sm text-gray-500">Detected At</h4>
                                <p>{selectedRecord.timestamp}</p>
                            </div>
                            <div>
                                <h4 className="font-semibold text-sm text-gray-500">Issue Description</h4>
                                <div className="bg-gray-100 p-2 rounded text-sm font-mono mt-1">
                                    {selectedRecord.details}
                                </div>
                            </div>
                            <div className="text-xs text-gray-400 mt-4">
                                * Full diff requires snapshot comparison
                            </div>
                        </div>
                        <div className="mt-6 flex justify-end">
                            <Button onClick={() => setSelectedRecord(null)}>Dismiss</Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
