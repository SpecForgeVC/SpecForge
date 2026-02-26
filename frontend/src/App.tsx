import { lazy, Suspense } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate, useParams } from "react-router-dom";
import MainLayout from "./components/layout/MainLayout";
import { AuthProvider } from "./hooks/use-auth";
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { NavigationProvider } from "./hooks/use-navigation";
import { WorkspaceMiddleware } from "./components/layout/WorkspaceMiddleware";
import { Loader2 } from "lucide-react";

// Lazy load pages to improve bundle size and initial load time
const DashboardPage = lazy(() => import("./features/projects/DashboardPage"));
const RoadmapItemPage = lazy(() => import("./features/roadmap/RoadmapItemPage"));
const WorkspaceListPage = lazy(() => import("./features/workspaces/WorkspaceListPage"));
const ProjectListPage = lazy(() => import("./features/projects/ProjectListPage"));
const RoadmapListPage = lazy(() => import("./features/roadmap/RoadmapListPage").then(m => ({ default: m.RoadmapListPage })));
const ContractListPage = lazy(() => import("./features/projects/ContractListPage").then(m => ({ default: m.ContractListPage })));
const VariableListPage = lazy(() => import("./features/projects/VariableListPage").then(m => ({ default: m.VariableListPage })));
const SnapshotListPage = lazy(() => import("./features/projects/SnapshotListPage").then(m => ({ default: m.SnapshotListPage })));
const RequirementsListPage = lazy(() => import("./features/roadmap/RequirementsListPage").then(m => ({ default: m.RequirementsListPage })));
const ValidationRulesListPage = lazy(() => import("./features/projects/ValidationRulesListPage").then(m => ({ default: m.ValidationRulesListPage })));
const WebhooksListPage = lazy(() => import("./features/projects/WebhooksListPage").then(m => ({ default: m.WebhooksListPage })));
const ProposalQueuePage = lazy(() => import("./features/proposals/ProposalQueuePage").then(m => ({ default: m.ProposalQueuePage })));
const LoginPage = lazy(() => import("./features/auth/LoginPage"));
const SettingsPage = lazy(() => import("./features/settings/SettingsPage"));
const IntelligenceDashboard = lazy(() => import("./features/intelligence/IntelligenceDashboard").then(m => ({ default: m.IntelligenceDashboard })));
const VariableLineagePage = lazy(() => import("./features/intelligence/VariableLineagePage").then(m => ({ default: m.VariableLineagePage })));
const ContractMutationReviewPage = lazy(() => import("./features/proposals/ContractMutationReviewPage").then(m => ({ default: m.ContractMutationReviewPage })));
const DriftHistoryPage = lazy(() => import("./features/intelligence/DriftHistoryPage").then(m => ({ default: m.DriftHistoryPage })));
const ProjectBootstrapDashboard = lazy(() => import("./features/projects/ProjectBootstrapDashboard").then(m => ({ default: m.ProjectBootstrapDashboard })));
const AlignmentDashboardPage = lazy(() => import("./features/projects/AlignmentDashboardPage").then(m => ({ default: m.AlignmentDashboardPage })));
const UIRoadmapListPage = lazy(() => import("./features/ui_roadmap/UIRoadmapListPage").then(m => ({ default: m.UIRoadmapListPage })));
const UIRoadmapWizardPage = lazy(() => import("./features/ui_roadmap/UIRoadmapWizardPage").then(m => ({ default: m.UIRoadmapWizardPage })));
const UIRoadmapItemPage = lazy(() => import("./features/ui_roadmap/UIRoadmapItemPage"));
const ImportWizardPage = lazy(() => import("./features/projects/ImportWizard").then(m => ({ default: m.default })));

const PageLoader = () => (
  <div className="flex h-screen w-full items-center justify-center">
    <Loader2 className="h-10 w-10 animate-spin text-indigo-600" />
  </div>
);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
});

function ProjectImportPageWrapper() {
  const { projectId } = useParams<{ projectId: string }>();
  return <ImportWizardPage projectId={projectId || ""} />;
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <NavigationProvider>
            <Suspense fallback={<PageLoader />}>
              <Routes>
                <Route path="/login" element={<LoginPage />} />

                <Route element={<ProtectedRoute />}>
                  <Route element={<WorkspaceMiddleware />}>
                    <Route element={<MainLayout />}>
                      <Route path="/" element={<Navigate to="/workspaces" replace />} />
                      <Route path="/workspaces" element={<WorkspaceListPage />} />
                      <Route path="/workspaces/:workspaceId/projects" element={<ProjectListPage />} />
                      <Route path="/projects/:projectId" element={<DashboardPage />} />
                      <Route path="/projects/:projectId/roadmap" element={<RoadmapListPage />} />
                      <Route path="/projects/:projectId/requirements" element={<RequirementsListPage />} />
                      <Route path="/projects/:projectId/contracts" element={<ContractListPage />} />
                      <Route path="/projects/:projectId/variables" element={<VariableListPage />} />
                      <Route path="/projects/:projectId/validation-rules" element={<ValidationRulesListPage />} />
                      <Route path="/projects/:projectId/webhooks" element={<WebhooksListPage />} />
                      <Route path="/projects/:projectId/proposals" element={<ProposalQueuePage />} />
                      <Route path="/projects/:projectId/snapshots" element={<SnapshotListPage />} />
                      <Route path="/roadmap/:roadmapItemId" element={<RoadmapItemPage />} />
                      <Route path="/roadmap/:roadmapItemId/intelligence" element={<IntelligenceDashboard />} />
                      <Route path="/settings" element={<SettingsPage />} />

                      <Route path="/variables/:variableId/lineage" element={<VariableLineagePage />} />
                      <Route path="/proposals/:proposalId/review" element={<ContractMutationReviewPage />} />
                      <Route path="/projects/:projectId/drift" element={<DriftHistoryPage />} />
                      <Route path="/projects/:projectId/bootstrap" element={<ProjectBootstrapDashboard />} />
                      <Route path="/projects/:projectId/alignment" element={<AlignmentDashboardPage />} />
                      <Route path="/projects/:projectId/import" element={<ProjectImportPageWrapper />} />
                      <Route path="/projects/:projectId/ui-roadmap" element={<UIRoadmapListPage />} />
                      <Route path="/projects/:projectId/ui-roadmap/new" element={<UIRoadmapWizardPage />} />
                      <Route path="/projects/:projectId/ui-roadmap/:id" element={<UIRoadmapItemPage />} />
                      <Route path="/projects/:projectId/ui-roadmap/:id/edit" element={<UIRoadmapWizardPage />} />
                    </Route>
                  </Route>
                </Route>

                <Route path="*" element={<Navigate to="/workspaces" replace />} />
              </Routes>
            </Suspense>
          </NavigationProvider>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  );
}
