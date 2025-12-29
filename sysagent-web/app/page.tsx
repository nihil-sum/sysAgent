'use client';

import { useState } from 'react';

export default function SysAgentTerminal() {
    const [input, setInput] = useState('');
    const [log, setLog] = useState<string[]>([]);
    const [loading, setLoading] = useState(false);

    const sendCommand = async () => {
        if (!input) return;

        setLoading(true);
        // Add user command to log
        setLog(prev => [...prev, `> ${input}`]);

        try {
            const res = await fetch('http://localhost:8080/api/chat', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ message: input }),
            });

            const data = await res.json();

            // Add agent response to log
            setLog(prev => [...prev, `SysAgent: ${data.response}`]);
        } catch (err) {
            setLog(prev => [...prev, `Error: Connection refused. Is Go server running?`]);
        }

        setLoading(false);
        setInput('');
    };

    return (
        <div className="min-h-screen bg-black text-green-400 p-8 font-mono">
            <h1 className="text-2xl mb-4 border-b border-green-800 pb-2">SysAgent // Terminal Interface</h1>

            <div className="mb-4 h-96 overflow-y-auto border border-green-900 p-4 rounded bg-gray-900">
                {log.map((line, i) => (
                    <div key={i} className="mb-2 whitespace-pre-wrap">{line}</div>
                ))}
                {loading && <div className="animate-pulse">_ Processing...</div>}
            </div>

            <div className="flex gap-2">
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && sendCommand()}
                    className="flex-1 bg-gray-800 border border-green-700 p-2 text-white focus:outline-none focus:border-green-400"
                    placeholder="Enter system command..."
                />
                <button
                    onClick={sendCommand}
                    disabled={loading}
                    className="bg-green-700 text-black px-6 py-2 hover:bg-green-600 disabled:opacity-50"
                >
                    EXECUTE
                </button>
            </div>
        </div>
    );
}