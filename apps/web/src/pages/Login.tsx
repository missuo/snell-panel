import { useState } from "react";
import { Button, Card, Input, Label, TextField } from "@heroui/react";
import { api } from "../api/client";
import { clearToken, setToken } from "../lib/auth";
import { Logo } from "../components/Logo";

export function Login({ onAuthed }: { onAuthed: () => void }) {
  const [value, setValue] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit() {
    if (!value.trim()) return;
    setLoading(true);
    setError("");
    setToken(value.trim());
    try {
      await api.settings();
      onAuthed();
    } catch {
      clearToken();
      setError("Invalid access token.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="relative flex min-h-full items-center justify-center p-6">
      <div className="pointer-events-none absolute inset-x-0 top-0 h-72 bg-gradient-to-b from-accent/10 to-transparent" />
      <Card className="relative w-full max-w-sm">
        <Card.Header>
          <Logo className="mb-3 h-11 w-11" />
          <Card.Title>Snell Panel</Card.Title>
          <Card.Description>Enter your access token to continue.</Card.Description>
        </Card.Header>
        <div className="flex flex-col gap-3 px-6 pb-2">
          <TextField value={value} onChange={setValue} type="password">
            <Label>Access Token</Label>
            <Input
              placeholder="access token"
              autoFocus
              onKeyDown={(e) => {
                if (e.key === "Enter") submit();
              }}
            />
          </TextField>
          {error && <p className="text-sm text-danger">{error}</p>}
        </div>
        <Card.Footer>
          <Button
            variant="primary"
            className="w-full"
            isDisabled={loading}
            onPress={submit}
          >
            {loading ? "Signing in…" : "Sign in"}
          </Button>
        </Card.Footer>
      </Card>
    </div>
  );
}
