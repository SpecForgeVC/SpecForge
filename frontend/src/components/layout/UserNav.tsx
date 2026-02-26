import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { LogOut, Settings } from "lucide-react";
import { Link } from "react-router-dom";

export function UserNav() {
    const { user, logout } = useAuth();

    return (
        <div className="flex items-center gap-3">
            <div className="flex flex-col items-end text-xs">
                <span className="font-semibold">{user?.name}</span>
                <span className="text-muted-foreground">{user?.role}</span>
            </div>
            <div className="h-8 w-8 rounded-full bg-slate-200 flex items-center justify-center cursor-pointer hover:ring-2 ring-primary/20 transition-all overflow-hidden border">
                <Avatar className="h-8 w-8">
                    <AvatarImage src={`https://avatar.vercel.sh/${user?.id}.png`} />
                    <AvatarFallback>{user?.name?.charAt(0)}</AvatarFallback>
                </Avatar>
            </div>
            <Button asChild variant="ghost" size="icon" className="h-8 w-8" title="Settings">
                <Link to="/settings">
                    <Settings className="h-4 w-4" />
                </Link>
            </Button>
            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={logout} title="Log out">
                <LogOut className="h-4 w-4" />
            </Button>
        </div>
    );
}
