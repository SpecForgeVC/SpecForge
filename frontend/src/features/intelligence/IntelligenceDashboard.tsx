import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ScoreCard } from './components/ScoreCard';
import { intelligenceApi } from '../../api/intelligence';
import type { FeatureIntelligence } from '../../api/intelligence';
import { webSocketService } from '../../api/websocket';

export const IntelligenceDashboard: React.FC = () => {
    const { roadmapItemId } = useParams<{ roadmapItemId: string }>();
    const [intelligence, setIntelligence] = useState<FeatureIntelligence | null>(null);

    useEffect(() => {
        const fetchIntelligence = async () => {
            if (!roadmapItemId) return;
            try {
                const data = await intelligenceApi.getFeatureIntelligence(roadmapItemId);
                setIntelligence(data);
            } catch (err) {
                console.error("Failed to fetch intelligence", err);
            } finally {
            }
        };

        fetchIntelligence();

        // Subscribe to real-time updates
        const handleUpdate = (payload: FeatureIntelligence) => {
            if (payload.feature_id === roadmapItemId) {
                setIntelligence(payload);
            }
        };

        webSocketService.subscribe('FEATURE_SCORE_UPDATED', handleUpdate);
        webSocketService.connect(); // Ensure connected

        return () => {
            webSocketService.unsubscribe('FEATURE_SCORE_UPDATED', handleUpdate);
        };
    }, [roadmapItemId]);

    if (!intelligence) return <div className="text-white">Loading Intelligence...</div>;

    return (
        <div className="p-6 bg-gray-900 min-h-screen text-white">
            <h1 className="text-3xl font-bold mb-6">Feature Intelligence Dashboard</h1>

            <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                <div className="col-span-1 md:col-span-4 bg-gray-800 p-6 rounded-lg shadow-lg flex items-center justify-between">
                    <div>
                        <h2 className="text-2xl font-semibold">Overall Healthy Score</h2>
                        <p className="text-gray-400">Aggregated from all metrics</p>
                    </div>
                    <div className="w-32 h-32">
                        <ScoreCard score={intelligence.overall_score} label="Overall" />
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <ScoreCard score={intelligence.completeness_score} label="Spec Completeness" />
                <ScoreCard score={intelligence.contract_integrity_score} label="Contract Integrity" />
                <ScoreCard score={intelligence.variable_coverage_score} label="Variable Coverage" />
                <ScoreCard score={intelligence.test_coverage_score} label="Test Coverage" />
                <ScoreCard score={intelligence.dependency_stability_score} label="Dependency Stability" />
                <ScoreCard score={intelligence.drift_risk_score} label="Drift Safety" />
                <ScoreCard score={intelligence.llm_confidence_score} label="AI Confidence" />
            </div>
        </div>
    );
};
