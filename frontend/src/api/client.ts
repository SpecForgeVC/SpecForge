import axios, { AxiosError } from "axios";
import type { InternalAxiosRequestConfig } from "axios";
import type { paths } from "./generated/schema";

const protocol = window.location.protocol;
const hostname = window.location.hostname;
// Use VITE_API_BASE_URL if set, otherwise fallback to the current hostname at port 8080
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || `${protocol}//${hostname}:8080/api/v1`;

// In-memory storage for access token
let accessToken: string | null = null;

export const setAccessToken = (token: string | null) => {
    accessToken = token;
};

export interface ApiResponse<T> {
    success: boolean;
    data: T;
    meta?: any;
    error?: {
        code: string;
        message: string;
        details?: string;
    };
}

export const getAccessToken = () => accessToken;

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        "Content-Type": "application/json",
    },
    withCredentials: true, // Required if using cookies for refresh token
});

// Request Interceptor
apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    if (accessToken) {
        console.log(`[DEBUG] apiClient: Injecting Auth header for ${config.url}`);
        config.headers.Authorization = `Bearer ${accessToken}`;
    } else {
        console.log(`[DEBUG] apiClient: No accessToken for ${config.url}`);
    }
    return config;
});

// Response Interceptor for Token Refresh
apiClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

        if (error.response?.status === 401 && !originalRequest._retry) {
            console.log(`[DEBUG] apiClient: Got 401 for ${originalRequest.url}, attempting refresh...`);
            originalRequest._retry = true;

            try {
                const refreshToken = localStorage.getItem("sf_refresh_token");
                if (!refreshToken) {
                    console.error("[DEBUG] apiClient: No refresh token found in storage.");
                    throw new Error("No refresh token available");
                }

                // Call refresh endpoint
                // Note: We use the base axios instance here to avoid the interceptors
                console.log("[DEBUG] apiClient: Calling /auth/refresh...");
                const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
                    refresh_token: refreshToken,
                });

                const { access_token, refresh_token } = response.data;
                console.log("[DEBUG] apiClient: Refresh successful.");

                // Update tokens
                setAccessToken(access_token);
                localStorage.setItem("sf_refresh_token", refresh_token);

                // Update header and retry original request
                if (originalRequest.headers) {
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                }

                console.log(`[DEBUG] apiClient: Retrying original request ${originalRequest.url}`);
                return apiClient(originalRequest);
            } catch (refreshError) {
                console.error("[DEBUG] apiClient: Refresh failed:", refreshError);
                // If refresh fails, clear tokens and redirect to login
                setAccessToken(null);
                localStorage.removeItem("sf_refresh_token");
                return Promise.reject(refreshError);
            }
        }

        return Promise.reject(error);
    }
);

// Type-safe helper for API calls
export type ApiPaths = keyof paths;
export type ApiMethod<P extends ApiPaths> = keyof paths[P] & ("get" | "post" | "patch" | "delete");

export type ApiResponseData<P extends ApiPaths, M extends ApiMethod<P>> =
    paths[P][M] extends { responses: { 200: { content: { "application/json": infer T } } } } ? T :
    paths[P][M] extends { responses: { 201: { content: { "application/json": infer T } } } } ? T :
    any;

export async function apiRequest<P extends ApiPaths, M extends ApiMethod<P>>(
    path: P,
    method: M,
    options?: {
        params?: any;
        body?: any;
    }
): Promise<ApiResponseData<P, M>> {
    console.log(`[DEBUG] apiRequest: ${method.toUpperCase()} ${path}`, options);
    const response = await apiClient.request({
        url: path.replace(/{(\w+)}/g, (_, key) => options?.params?.[key] || ""),
        method,
        data: options?.body,
        params: options?.params,
    });
    return response.data;
}
