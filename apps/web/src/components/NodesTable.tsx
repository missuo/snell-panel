import { useState } from "react";
import { Button, Chip, Spinner, Table } from "@heroui/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { NodeDTO } from "@snell-panel/shared";
import { api } from "../api/client";
import { useNodes } from "../api/hooks";
import { countryFlag } from "../lib/format";
import { CommandModal } from "./CommandModal";
import { RelayModal } from "./RelayModal";
import { RenameModal } from "./RenameModal";

type Action = {
  type: "install" | "upgrade" | "relay" | "rename";
  node: NodeDTO;
};

export function NodesTable() {
  const qc = useQueryClient();
  const nodes = useNodes();
  const [action, setAction] = useState<Action | null>(null);

  const del = useMutation({
    mutationFn: (id: string) => api.deleteNode(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["nodes"] }),
  });

  const toggle = useMutation({
    mutationFn: (v: { id: string; enabled: boolean }) =>
      api.patchNode(v.id, { enabled: v.enabled }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["nodes"] }),
  });

  if (nodes.isLoading)
    return (
      <div className="flex justify-center p-10">
        <Spinner />
      </div>
    );
  if (nodes.isError)
    return <p className="text-danger">{(nodes.error as Error).message}</p>;

  const list = nodes.data ?? [];
  const close = (open: boolean) => !open && setAction(null);

  return (
    <>
      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label="Nodes" className="min-w-[820px]">
            <Table.Header>
              <Table.Column isRowHeader>Name</Table.Column>
              <Table.Column>Status</Table.Column>
              <Table.Column>Ver</Table.Column>
              <Table.Column>IP</Table.Column>
              <Table.Column>Port</Table.Column>
              <Table.Column>ISP / ASN</Table.Column>
              <Table.Column>Actions</Table.Column>
            </Table.Header>
            <Table.Body>
              {list.map((n) => (
                <Table.Row key={n.node_id} id={n.node_id}>
                  <Table.Cell>
                    {countryFlag(n.country_code)} {n.node_name}
                  </Table.Cell>
                  <Table.Cell>
                    <div className="flex flex-wrap gap-1.5">
                      <Chip color={n.status === "active" ? "success" : "warning"}>
                        {n.status}
                      </Chip>
                      {!n.enabled && <Chip>hidden</Chip>}
                    </div>
                  </Table.Cell>
                  <Table.Cell>V{n.version}</Table.Cell>
                  <Table.Cell>{n.ip ?? "—"}</Table.Cell>
                  <Table.Cell>{n.port ?? "—"}</Table.Cell>
                  <Table.Cell>
                    <div className="flex flex-col leading-tight">
                      <span>{n.isp ?? "—"}</span>
                      {n.asn != null && (
                        <span className="text-xs text-muted">AS{n.asn}</span>
                      )}
                    </div>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="flex flex-wrap gap-1.5">
                      {n.status === "pending" ? (
                        <Button
                          size="sm"
                          variant="primary"
                          onPress={() => setAction({ type: "install", node: n })}
                        >
                          Install
                        </Button>
                      ) : (
                        <>
                          <Button
                            size="sm"
                            variant="secondary"
                            onPress={() => setAction({ type: "relay", node: n })}
                          >
                            Relay
                          </Button>
                          {n.version === "5" && (
                            <Button
                              size="sm"
                              variant="secondary"
                              onPress={() =>
                                setAction({ type: "upgrade", node: n })
                              }
                            >
                              Upgrade
                            </Button>
                          )}
                        </>
                      )}
                      <Button
                        size="sm"
                        variant="outline"
                        isDisabled={toggle.isPending}
                        onPress={() =>
                          toggle.mutate({ id: n.node_id, enabled: !n.enabled })
                        }
                      >
                        {n.enabled ? "Disable" : "Enable"}
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onPress={() => setAction({ type: "rename", node: n })}
                      >
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        variant="danger"
                        onPress={() => {
                          if (confirm(`Delete "${n.node_name}"?`)) {
                            del.mutate(n.node_id);
                          }
                        }}
                      >
                        Delete
                      </Button>
                    </div>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>

      {list.length === 0 && (
        <p className="p-6 text-center text-sm text-muted">
          No nodes yet. Click “Add Node” to create one.
        </p>
      )}

      {action?.type === "install" && (
        <CommandModal node={action.node} purpose="install" isOpen onOpenChange={close} />
      )}
      {action?.type === "upgrade" && (
        <CommandModal node={action.node} purpose="upgrade" isOpen onOpenChange={close} />
      )}
      {action?.type === "relay" && (
        <RelayModal origin={action.node} isOpen onOpenChange={close} />
      )}
      {action?.type === "rename" && (
        <RenameModal node={action.node} isOpen onOpenChange={close} />
      )}
    </>
  );
}
