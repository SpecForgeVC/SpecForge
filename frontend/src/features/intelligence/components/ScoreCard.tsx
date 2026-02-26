import React from 'react';

interface ScoreCardProps {
    score: number;
    label: string;
    trend?: 'up' | 'down' | 'stable';
}

export const ScoreCard: React.FC<ScoreCardProps> = ({ score, label, trend }) => {
    const getColor = (s: number) => {
        if (s >= 80) return 'text-green-500';
        if (s >= 50) return 'text-yellow-500';
        return 'text-red-500';
    };

    const ringColor = (s: number) => {
        if (s >= 80) return 'stroke-green-500';
        if (s >= 50) return 'stroke-yellow-500';
        return 'stroke-red-500';
    }

    // Circular progress calculation
    const radius = 30;
    const circumference = 2 * Math.PI * radius;
    const offset = circumference - (score / 100) * circumference;

    return (
        <div className="bg-gray-800 p-4 rounded-lg shadow-md flex flex-col items-center justify-center w-32 h-40">
            <div className="relative w-20 h-20 mb-2">
                {/* Background Ring */}
                <svg className="w-full h-full transform -rotate-90">
                    <circle
                        cx="40"
                        cy="40"
                        r={radius}
                        stroke="currentColor"
                        strokeWidth="6"
                        fill="transparent"
                        className="text-gray-700"
                    />
                    {/* Progress Ring */}
                    <circle
                        cx="40"
                        cy="40"
                        r={radius}
                        stroke="currentColor"
                        strokeWidth="6"
                        fill="transparent"
                        className={`${ringColor(score)} transition-all duration-500 ease-in-out`}
                        strokeDasharray={circumference}
                        strokeDashoffset={offset}
                        strokeLinecap="round"
                    />
                </svg>
                <span className={`absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-xl font-bold ${getColor(score)}`}>
                    {score}
                </span>
            </div>
            <span className="text-xs text-gray-400 font-medium text-center">{label}</span>
            {trend && (
                <span className={`text-xs mt-1 ${trend === 'up' ? 'text-green-400' : trend === 'down' ? 'text-red-400' : 'text-gray-500'}`}>
                    {trend === 'up' ? '↑' : trend === 'down' ? '↓' : '•'}
                </span>
            )}
        </div>
    );
};
