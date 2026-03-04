import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LLMConfigSection } from "./components/LLMConfigSection";
import { Settings, Zap, Shield, Bell, User } from "lucide-react";
import { useAuth } from "@/hooks/use-auth";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

export default function SettingsPage() {
    const { user } = useAuth();

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-3xl font-bold tracking-tight">Settings</h2>
                <p className="text-muted-foreground">Manage your application preferences and integrations.</p>
            </div>

            <Tabs defaultValue="llm" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="llm" className="flex gap-2">
                        <Zap className="h-4 w-4" /> Intelligence (LLM)
                    </TabsTrigger>
                    <TabsTrigger value="general" className="flex gap-2">
                        <User className="h-4 w-4" /> General
                    </TabsTrigger>
                    <TabsTrigger value="security" className="flex gap-2">
                        <Shield className="h-4 w-4" /> Security
                    </TabsTrigger>
                    <TabsTrigger value="notifications" className="flex gap-2" disabled>
                        <Bell className="h-4 w-4" /> Notifications
                        <Badge variant="outline" className="text-[10px] ml-1">Soon</Badge>
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="llm" className="space-y-4 pt-4">
                    <LLMConfigSection />
                </TabsContent>

                <TabsContent value="general" className="space-y-4 pt-4">
                    <GeneralTab user={user} />
                </TabsContent>

                <TabsContent value="security" className="space-y-4 pt-4">
                    <SecurityTab user={user} />
                </TabsContent>

                <TabsContent value="notifications">
                    {/* Intentionally empty — notifications tab is disabled */}
                </TabsContent>
            </Tabs>
        </div>
    );
}

function GeneralTab({ user }: { user: { id: string; workspace_id: string; role: string; name?: string } | null }) {
    return (
        <div className="space-y-4 max-w-2xl">
            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Account</CardTitle>
                    <CardDescription>Your identity within this SpecForge workspace.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">User ID</p>
                            <p className="font-mono text-xs break-all">{user?.id ?? "—"}</p>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Role</p>
                            <Badge variant="secondary">{user?.role ?? "—"}</Badge>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Workspace ID</p>
                            <p className="font-mono text-xs break-all">{user?.workspace_id ?? "—"}</p>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Display Name</p>
                            <p className="text-sm">{user?.name ?? "SpecForge User"}</p>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Preferences</CardTitle>
                    <CardDescription>UI and workflow preferences.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                    <div className="flex items-center justify-between py-2">
                        <div>
                            <p className="text-sm font-medium">Default Roadmap View</p>
                            <p className="text-xs text-muted-foreground">How roadmap items are displayed by default.</p>
                        </div>
                        <Badge variant="outline">List</Badge>
                    </div>
                    <Separator />
                    <div className="flex items-center justify-between py-2">
                        <div>
                            <p className="text-sm font-medium">AI Assistance</p>
                            <p className="text-xs text-muted-foreground">Inline AI features during editing.</p>
                        </div>
                        <Badge variant="outline">Enabled</Badge>
                    </div>
                    <p className="text-xs text-muted-foreground pt-2">
                        <Settings className="h-3 w-3 inline mr-1" />
                        Additional preferences will be available in a future update.
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}

function SecurityTab({ user }: { user: { id: string; workspace_id: string; role: string; name?: string } | null }) {
    const sessionIssuedAt = new Date().toLocaleString(); // Approximation — real value would come from the JWT claims

    return (
        <div className="space-y-4 max-w-2xl">
            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Current Session</CardTitle>
                    <CardDescription>Details about your active authentication session.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">User ID</p>
                            <p className="font-mono text-xs break-all">{user?.id ?? "—"}</p>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Role</p>
                            <Badge variant="secondary">{user?.role ?? "—"}</Badge>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Session Started (approx.)</p>
                            <p className="text-xs">{sessionIssuedAt}</p>
                        </div>
                        <div>
                            <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Token Storage</p>
                            <p className="text-xs">LocalStorage (refresh token)</p>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Password & Access</CardTitle>
                    <CardDescription>Manage your credentials and access controls.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                    <div className="flex items-center justify-between py-1">
                        <div>
                            <p className="text-sm font-medium">Change Password</p>
                            <p className="text-xs text-muted-foreground">Update your login credentials.</p>
                        </div>
                        <Badge variant="outline" className="text-muted-foreground">Coming Soon</Badge>
                    </div>
                    <Separator />
                    <div className="flex items-center justify-between py-1">
                        <div>
                            <p className="text-sm font-medium">Two-Factor Authentication</p>
                            <p className="text-xs text-muted-foreground">Add an extra layer of security.</p>
                        </div>
                        <Badge variant="outline" className="text-muted-foreground">Coming Soon</Badge>
                    </div>
                    <Separator />
                    <div className="flex items-center justify-between py-1">
                        <div>
                            <p className="text-sm font-medium">Active Sessions</p>
                            <p className="text-xs text-muted-foreground">View and revoke other active login sessions.</p>
                        </div>
                        <Badge variant="outline" className="text-muted-foreground">Coming Soon</Badge>
                    </div>
                    <p className="text-xs text-muted-foreground pt-2">
                        <Shield className="h-3 w-3 inline mr-1" />
                        Password management and session controls will be available in a future update.
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}
