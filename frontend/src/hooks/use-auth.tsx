import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import { setAccessToken, apiClient } from "@/api/client";

export type Role = "OWNER" | "ADMIN" | "REVIEWER" | "ENGINEER" | "AI_AGENT";

interface User {
    id: string;
    workspace_id: string;
    role: Role;
    name?: string; // Optional for now
}

interface AuthResponse {
    access_token: string;
    refresh_token: string;
    expires_in: number;
}

interface AuthContextType {
    user: User | null;
    login: (tokens: AuthResponse) => Promise<void>;
    logout: () => void;
    isAuthenticated: boolean;
    canPerform: (action: string) => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const fetchUserInfo = async () => {
        try {
            const response = await apiClient.get("/auth/me");
            const data = response.data;
            setUser({
                id: data.user_id,
                workspace_id: data.workspace_id,
                role: data.role as Role,
                name: "SpecForge User", // Default until name is in /auth/me
            });
            console.log("[AUTH] User info fetched successfully.");
        } catch (err) {
            console.error("[AUTH] Failed to fetch user info:", err);
            // This might happen if the token just expired before we could fetch
            // logout() is called by the interceptor if 401 occurs repeatedly
        }
    };

    useEffect(() => {
        const checkAuth = async () => {
            const refreshToken = localStorage.getItem("sf_refresh_token");
            if (refreshToken) {
                console.log("[AUTH] Refresh token found, attempting silent refresh...");
                try {
                    // We call the API directly to refresh the token on mount
                    const response = await apiClient.post("/auth/refresh", {
                        refresh_token: refreshToken,
                    });
                    const { access_token, refresh_token } = response.data;

                    setAccessToken(access_token);
                    localStorage.setItem("sf_refresh_token", refresh_token);

                    // Fetch real user info
                    await fetchUserInfo();
                } catch (err) {
                    console.error("[AUTH] Silent refresh failed:", err);
                    localStorage.removeItem("sf_refresh_token");
                    setUser(null);
                }
            } else {
                console.log("[AUTH] No refresh token found.");
            }
            setIsLoading(false);
        };
        checkAuth();
    }, []);

    const login = async (tokens: AuthResponse) => {
        setAccessToken(tokens.access_token);

        try {
            localStorage.setItem("sf_refresh_token", tokens.refresh_token);
        } catch (e: any) {
            console.warn("[AUTH] localStorage quota exceeded, attempting to clear and retry storage.", e);
            localStorage.clear();
            try {
                localStorage.setItem("sf_refresh_token", tokens.refresh_token);
            } catch (retryError) {
                console.error("[AUTH] Failed to store refresh token even after clear:", retryError);
            }
        }

        await fetchUserInfo();
    };

    const logout = () => {
        setAccessToken(null);
        localStorage.removeItem("sf_refresh_token");
        setUser(null);
    };

    const canPerform = (action: string) => {
        if (!user) return false;
        // Basic RBAC logic
        if (user.role === "OWNER" || user.role === "ADMIN") return true;
        if (user.role === "AI_AGENT" && action.startsWith("mutate")) return false;
        if (user.role === "REVIEWER" && action.includes("approve")) return true;
        return true;
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user, canPerform }}>
            {!isLoading && children}
        </AuthContext.Provider>
    );
}

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
