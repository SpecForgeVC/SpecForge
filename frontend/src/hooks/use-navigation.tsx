import { createContext, useContext, useState, type ReactNode } from "react";

interface NavigationContextType {
    activeWorkspaceId: string | null;
    activeProjectId: string | null;
    setActiveWorkspace: (id: string | null) => void;
    setActiveProject: (id: string | null) => void;
    getProjectLink: (path: string) => string;
}

const NavigationContext = createContext<NavigationContextType | undefined>(undefined);

export function NavigationProvider({ children }: { children: ReactNode }) {
    const [activeWorkspaceId, setActiveWorkspaceId] = useState<string | null>(() =>
        localStorage.getItem("sf_active_workspace")
    );
    const [activeProjectId, setActiveProjectId] = useState<string | null>(() =>
        localStorage.getItem("sf_active_project")
    );

    const setActiveWorkspace = (id: string | null) => {
        setActiveWorkspaceId(id);
        if (id) {
            localStorage.setItem("sf_active_workspace", id);
        } else {
            localStorage.removeItem("sf_active_workspace");
        }
    };

    const setActiveProject = (id: string | null) => {
        setActiveProjectId(id);
        if (id) {
            localStorage.setItem("sf_active_project", id);
        } else {
            localStorage.removeItem("sf_active_project");
        }
    };

    const getProjectLink = (path: string) => {
        if (!activeProjectId) return "/workspaces";
        return path.replace(":id", activeProjectId).replace(":projectId", activeProjectId);
    };

    return (
        <NavigationContext.Provider
            value={{
                activeWorkspaceId,
                activeProjectId,
                setActiveWorkspace,
                setActiveProject,
                getProjectLink
            }}
        >
            {children}
        </NavigationContext.Provider>
    );
}

export const useNavigation = () => {
    const context = useContext(NavigationContext);
    if (context === undefined) {
        throw new Error("useNavigation must be used within a NavigationProvider");
    }
    return context;
};
