import Editor from "@monaco-editor/react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Save, RefreshCw } from "lucide-react";

interface ContractEditorProps {
    initialValue?: string;
    onSave?: (value: string) => void;
    readOnly?: boolean;
}

export function ContractEditor({ initialValue = "{}", onSave, readOnly = false }: ContractEditorProps) {
    const [value, setValue] = useState(initialValue);

    const handleEditorChange = (newValue: string | undefined) => {
        setValue(newValue || "{}");
    };

    return (
        <div className="flex flex-col h-full gap-4">
            <div className="flex items-center justify-between">
                <span className="text-xs font-mono text-muted-foreground">JSON Schema (Draft 2020-12)</span>
                <div className="flex gap-2">
                    <Button variant="outline" size="sm">
                        <RefreshCw className="mr-2 h-3.5 w-3.5" /> Validate
                    </Button>
                    {!readOnly && (
                        <Button size="sm" onClick={() => onSave?.(value)}>
                            <Save className="mr-2 h-3.5 w-3.5" /> Save Changes
                        </Button>
                    )}
                </div>
            </div>
            <div className="flex-1 border rounded-md overflow-hidden">
                <Editor
                    height="100%"
                    defaultLanguage="json"
                    theme="vs-light"
                    value={value}
                    onChange={handleEditorChange}
                    options={{
                        readOnly,
                        minimap: { enabled: false },
                        fontSize: 13,
                        scrollBeyondLastLine: false,
                        automaticLayout: true,
                    }}
                />
            </div>
        </div>
    );
}
