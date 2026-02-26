import { useState, useCallback, useMemo } from 'react';

export interface StateMachineNode {
    id: string;
    name: string;
    type: 'initial' | 'normal' | 'final' | 'error';
    visual_mutations: Record<string, any>;
    interaction_mutations: Record<string, any>;
}

export interface StateMachineTransition {
    id: string;
    source: string;
    target: string;
    event: string;
    condition?: string;
    action?: string;
}

export interface StateMachineDef {
    nodes: StateMachineNode[];
    transitions: StateMachineTransition[];
}

export const useStateMachineRunner = (definition: StateMachineDef) => {
    const [activeStateId, setActiveStateId] = useState<string | null>(() => {
        const initialNode = definition.nodes.find(n => n.type === 'initial');
        return initialNode ? initialNode.id : (definition.nodes[0]?.id || null);
    });

    const [history, setHistory] = useState<string[]>([]);
    const [lastEvent, setLastEvent] = useState<string | null>(null);

    const activeState = useMemo(() =>
        definition.nodes.find(n => n.id === activeStateId),
        [activeStateId, definition.nodes]);

    const availableTransitions = useMemo(() =>
        definition.transitions.filter(t => t.source === activeStateId),
        [activeStateId, definition.transitions]);

    const trigger = useCallback((event: string) => {
        const transition = definition.transitions.find(
            t => t.source === activeStateId && t.event === event
        );

        if (transition) {
            if (activeStateId) {
                setHistory(prev => [...prev, activeStateId]);
            }
            setActiveStateId(transition.target);
            setLastEvent(event);
            return true;
        }
        return false;
    }, [activeStateId, definition.transitions]);

    const reset = useCallback(() => {
        const initialNode = definition.nodes.find(n => n.type === 'initial');
        setActiveStateId(initialNode ? initialNode.id : (definition.nodes[0]?.id || null));
        setHistory([]);
        setLastEvent(null);
    }, [definition.nodes]);

    return {
        activeStateId,
        activeState,
        availableTransitions,
        history,
        lastEvent,
        trigger,
        reset
    };
};
