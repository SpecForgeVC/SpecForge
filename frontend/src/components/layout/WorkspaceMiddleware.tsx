import { useEffect } from "react";
import { useParams, Outlet } from "react-router-dom";
import { useNavigation } from "@/hooks/use-navigation";

export function WorkspaceMiddleware() {
    const { workspaceId, projectId } = useParams<{ workspaceId?: string; projectId?: string }>();
    const { setActiveWorkspace, setActiveProject } = useNavigation();

    useEffect(() => {
        if (workspaceId) {
            setActiveWorkspace(workspaceId);
        }
    }, [workspaceId, setActiveWorkspace]);

    useEffect(() => {
        if (projectId) {
            setActiveProject(projectId);
        }
    }, [projectId, setActiveProject]);

    return <Outlet />;
}
