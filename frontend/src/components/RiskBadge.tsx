import React from 'react';

interface RiskBadgeProps {
    level: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
}

export const RiskBadge: React.FC<RiskBadgeProps> = ({ level }) => {
    const getStyles = (l: string) => {
        switch (l) {
            case 'LOW':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'MEDIUM':
                return 'bg-yellow-100 text-yellow-800 border-yellow-200';
            case 'HIGH':
                return 'bg-orange-100 text-orange-800 border-orange-200';
            case 'CRITICAL':
                return 'bg-red-100 text-red-800 border-red-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
        }
    };

    return (
        <span className={`px-2 py-1 rounded-full text-xs font-semibold border ${getStyles(level)}`}>
            {level}
        </span>
    );
};
