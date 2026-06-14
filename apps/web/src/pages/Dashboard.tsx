import { useState } from "react";
import { Button } from "@heroui/react";
import { clearToken } from "../lib/auth";
import { NodesTable } from "../components/NodesTable";
import { AddNodeModal } from "../components/AddNodeModal";
import { SubscriptionModal } from "../components/SubscriptionModal";
import { ThemeToggle } from "../components/ThemeToggle";

export function Dashboard({ onLogout }: { onLogout: () => void }) {
  const [addOpen, setAddOpen] = useState(false);
  const [subOpen, setSubOpen] = useState(false);

  return (
    <div className="mx-auto flex max-w-6xl flex-col gap-6 p-4 sm:p-6">
      <header className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Snell Panel</h1>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <Button
            variant="tertiary"
            onPress={() => {
              clearToken();
              onLogout();
            }}
          >
            Logout
          </Button>
        </div>
      </header>

      <div className="flex items-center justify-between">
        <h2 className="text-lg font-medium">Nodes</h2>
        <div className="flex gap-2">
          <Button variant="secondary" onPress={() => setSubOpen(true)}>
            Subscription
          </Button>
          <Button variant="primary" onPress={() => setAddOpen(true)}>
            Add Node
          </Button>
        </div>
      </div>

      <NodesTable />

      <AddNodeModal isOpen={addOpen} onOpenChange={setAddOpen} />
      <SubscriptionModal isOpen={subOpen} onOpenChange={setSubOpen} />
    </div>
  );
}
