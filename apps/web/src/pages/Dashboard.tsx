import { useState } from "react";
import { Button } from "@heroui/react";
import { clearToken } from "../lib/auth";
import { SubscriptionCard } from "../components/SubscriptionCard";
import { NodesTable } from "../components/NodesTable";
import { AddNodeModal } from "../components/AddNodeModal";
import { ThemeToggle } from "../components/ThemeToggle";

export function Dashboard({ onLogout }: { onLogout: () => void }) {
  const [addOpen, setAddOpen] = useState(false);

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

      <SubscriptionCard />

      <div className="flex items-center justify-between">
        <h2 className="text-lg font-medium">Nodes</h2>
        <Button variant="primary" onPress={() => setAddOpen(true)}>
          Add Node
        </Button>
      </div>

      <NodesTable />

      <AddNodeModal isOpen={addOpen} onOpenChange={setAddOpen} />
    </div>
  );
}
