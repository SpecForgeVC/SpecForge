import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LLMConfigSection } from "./components/LLMConfigSection";
import { Settings, Zap, Shield, Bell } from "lucide-react";

export default function SettingsPage() {
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
                    <TabsTrigger value="general" className="flex gap-2" disabled>
                        <Settings className="h-4 w-4" /> General
                    </TabsTrigger>
                    <TabsTrigger value="security" className="flex gap-2" disabled>
                        <Shield className="h-4 w-4" /> Security
                    </TabsTrigger>
                    <TabsTrigger value="notifications" className="flex gap-2" disabled>
                        <Bell className="h-4 w-4" /> Notifications
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="llm" className="space-y-4 pt-4">
                    <LLMConfigSection />
                </TabsContent>

                <TabsContent value="general">
                    {/* Placeholder */}
                </TabsContent>
            </Tabs>
        </div>
    );
}
