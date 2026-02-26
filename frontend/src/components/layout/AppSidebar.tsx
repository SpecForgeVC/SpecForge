import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarFooter,
} from "@/components/ui/sidebar";
import { LayoutDashboard, ListTree, Settings, ShieldCheck, Database, History, FileText, CheckSquare, Webhook, Sparkles, Terminal } from "lucide-react";
import { Link, useLocation } from "react-router-dom";

import { useNavigation } from "@/hooks/use-navigation";

const menuItems = [
    { title: "Dashboard", icon: LayoutDashboard, url: "/projects/:id" },
    { title: "Requirements", icon: FileText, url: "/projects/:id/requirements" },
    { title: "API Roadmap", icon: ListTree, url: "/projects/:id/roadmap" },
    { title: "UI Roadmap", icon: Sparkles, url: "/projects/:id/ui-roadmap" },
    { title: "Import IDE", icon: Terminal, url: "/projects/:id/import" },
    { title: "Contracts", icon: ShieldCheck, url: "/projects/:id/contracts" },
    { title: "Variables", icon: Database, url: "/projects/:id/variables" },
    { title: "Validation Rules", icon: CheckSquare, url: "/projects/:id/validation-rules" },
    { title: "Webhooks", icon: Webhook, url: "/projects/:id/webhooks" },
    { title: "AI Proposals", icon: Sparkles, url: "/projects/:id/proposals" },
    { title: "Drift Reports", icon: History, url: "/projects/:id/drift" },
    { title: "Snapshots", icon: Database, url: "/projects/:id/snapshots" },
];

export function AppSidebar() {
    const location = useLocation();
    const { getProjectLink } = useNavigation();

    // Use activeProjectId from context, which is synced with URL and localStorage

    const getUrl = (url: string) => {
        return getProjectLink(url);
    };

    return (
        <Sidebar>
            <SidebarHeader className="p-4 border-b">
                <Link to="/workspaces" className="flex items-center gap-2">
                    <img src="/logo_dark.png" alt="SpecForge Logo" className="h-6 w-6 object-contain" />
                    <h1 className="text-xl font-bold tracking-tight">SpecForge</h1>
                </Link>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>Project Workspace</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {menuItems.map((item) => {
                                const url = getUrl(item.url);
                                return (
                                    <SidebarMenuItem key={item.title}>
                                        <SidebarMenuButton asChild isActive={location.pathname === url} tooltip={item.title}>
                                            <Link to={url}>
                                                <item.icon className="h-4 w-4" />
                                                <span>{item.title}</span>
                                            </Link>
                                        </SidebarMenuButton>
                                    </SidebarMenuItem>
                                );
                            })}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
            <SidebarFooter className="p-4 border-t">
                <div className="flex items-center gap-2 text-muted-foreground hover:text-foreground cursor-pointer transition-colors">
                    <Settings className="h-4 w-4" />
                    <span className="text-sm font-medium">Settings</span>
                </div>
            </SidebarFooter>
        </Sidebar>
    );
}
