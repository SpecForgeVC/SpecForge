import { Link, useLocation } from "react-router-dom";
import { ChevronRight, Home } from "lucide-react";

export function Breadcrumbs() {
    const location = useLocation();
    const pathnames = location.pathname.split("/").filter((x) => x);

    return (
        <nav className="flex items-center space-x-1 text-sm font-medium text-muted-foreground">
            <Link to="/" className="flex items-center hover:text-foreground transition-colors">
                <Home className="h-4 w-4" />
            </Link>
            {pathnames.length > 0 && <ChevronRight className="h-4 w-4" />}
            {pathnames.map((value, index) => {
                const last = index === pathnames.length - 1;
                const to = `/${pathnames.slice(0, index + 1).join("/")}`;
                const label = value.charAt(0).toUpperCase() + value.slice(1).replace(/-/g, " ");

                return (
                    <div key={to} className="flex items-center space-x-1">
                        {last ? (
                            <span className="text-foreground font-semibold">{label}</span>
                        ) : (
                            <Link to={to} className="hover:text-foreground transition-colors">
                                {label}
                            </Link>
                        )}
                        {!last && <ChevronRight className="h-4 w-4" />}
                    </div>
                );
            })}
        </nav>
    );
}
