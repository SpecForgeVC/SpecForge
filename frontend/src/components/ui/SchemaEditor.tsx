import React, { useState, useEffect } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

interface SchemaEditorProps {
    label: string;
    description?: string;
    initialValue: any;
    onChange: (value: any) => void;
}

export function SchemaEditor({ label, description, initialValue, onChange }: SchemaEditorProps) {
    const [jsonString, setJsonString] = useState(JSON.stringify(initialValue || {}, null, 2));
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setJsonString(JSON.stringify(initialValue || {}, null, 2));
    }, [initialValue]);

    const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const value = e.target.value;
        setJsonString(value);

        try {
            const parsed = JSON.parse(value);
            setError(null);
            onChange(parsed);
        } catch (err: any) {
            setError("Invalid JSON format");
        }
    };

    return (
        <div className="space-y-4 pt-4 border-t border-border mt-4">
            <div>
                <Label className="text-sm font-semibold">{label}</Label>
                {description && <p className="text-xs text-muted-foreground">{description}</p>}
            </div>
            <Textarea
                value={jsonString}
                onChange={handleTextChange}
                placeholder='{ "type": "object", ... }'
                className="font-mono text-xs min-h-[200px]"
            />
            {error && (
                <Alert variant="destructive" className="py-2">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription className="text-xs">{error}</AlertDescription>
                </Alert>
            )}
        </div>
    );
}
