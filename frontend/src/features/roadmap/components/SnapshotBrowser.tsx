import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { History, Eye, ArrowLeftRight } from "lucide-react";

interface Snapshot {
    id: string;
    version: string;
    created_at: string;
    created_by: string;
    comment: string;
}

const mockSnapshots: Snapshot[] = [
    { id: "1", version: "v1.2.0", created_at: "2026-02-10T14:30:00Z", created_by: "Admin", comment: "Approved production schema" },
    { id: "2", version: "v1.1.0", created_at: "2026-02-05T09:15:00Z", created_by: "AI Agent", comment: "Auto-snapshot before OAuth change" },
    { id: "3", version: "v1.0.0", created_at: "2026-01-20T11:00:00Z", created_by: "System", comment: "Initial baseline" },
];

export function SnapshotBrowser() {
    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h3 className="font-semibold text-lg">Snapshot History</h3>
                <Button variant="outline" size="sm">
                    <ArrowLeftRight className="mr-2 h-4 w-4" /> Compare Any Two
                </Button>
            </div>

            <div className="space-y-3">
                {mockSnapshots.map((snapshot) => (
                    <Card key={snapshot.id} className="hover:bg-slate-50 transition-colors cursor-pointer">
                        <CardContent className="p-4 flex items-center justify-between">
                            <div className="flex items-center gap-4">
                                <div className="h-10 w-10 bg-slate-100 rounded-full flex items-center justify-center">
                                    <History className="h-5 w-5 text-muted-foreground" />
                                </div>
                                <div className="flex flex-col">
                                    <div className="flex items-center gap-2">
                                        <span className="font-bold">{snapshot.version}</span>
                                        <span className="text-[10px] bg-slate-200 px-1.5 py-0.5 rounded text-muted-foreground">ID: {snapshot.id}</span>
                                    </div>
                                    <span className="text-xs text-muted-foreground">{snapshot.comment}</span>
                                </div>
                            </div>
                            <div className="flex flex-col items-end gap-1">
                                <span className="text-xs font-medium">{new Date(snapshot.created_at).toLocaleDateString()}</span>
                                <div className="flex gap-2">
                                    <Button variant="ghost" size="sm" className="h-7 px-2">
                                        <Eye className="h-3.5 w-3.5 mr-1" /> View
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
