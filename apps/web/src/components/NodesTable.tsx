import { useState } from "react";
import { Button, Chip, Dropdown, Label, Spinner, Table } from "@heroui/react";
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

function MoreIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" className="size-4" aria-hidden="true">
      <circle cx="5" cy="12" r="1.6" />
      <circle cx="12" cy="12" r="1.6" />
      <circle cx="19" cy="12" r="1.6" />
    </svg>
  );
}

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

  function handleMenu(key: string, n: NodeDTO) {
    switch (key) {
      case "relay":
        setAction({ type: "relay", node: n });
        break;
      case "upgrade":
        setAction({ type: "upgrade", node: n });
        break;
      case "toggle":
        toggle.mutate({ id: n.node_id, enabled: !n.enabled });
        break;
      case "edit":
        setAction({ type: "rename", node: n });
        break;
      case "delete":
        if (confirm(`Delete "${n.node_name}"?`)) del.mutate(n.node_id);
        break;
    }
  }

  return (
    <>
      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label="Nodes" className="min-w-[760px] text-sm">
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
                  <Table.Cell>
                    <span className="font-mono text-xs">{n.ip ?? "—"}</span>
                  </Table.Cell>
                  <Table.Cell>
                    <span className="font-mono text-xs">{n.port ?? "—"}</span>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="flex flex-col leading-tight">
                      <span>{n.isp ?? "—"}</span>
                      {n.asn != null && (
                        <span className="font-mono text-xs text-muted">
                          AS{n.asn}
                        </span>
                      )}
                    </div>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="flex items-center justify-end gap-1.5">
                      {n.status === "pending" && (
                        <Button
                          size="sm"
                          variant="primary"
                          onPress={() => setAction({ type: "install", node: n })}
                        >
                          Install
                        </Button>
                      )}
                      <Dropdown>
                        <Button
                          size="sm"
                          variant="outline"
                          aria-label="More actions"
                          className="px-2"
                        >
                          <MoreIcon />
                        </Button>
                        <Dropdown.Popover>
                          <Dropdown.Menu
                            onAction={(key) => handleMenu(String(key), n)}
                          >
                            <Dropdown.Item
                              id="relay"
                              textValue="Add relay"
                              isDisabled={n.status !== "active"}
                            >
                              <Label>Add relay</Label>
                            </Dropdown.Item>
                            {n.status === "active" && n.version === "5" && (
                              <Dropdown.Item id="upgrade" textValue="Upgrade to V6">
                                <Label>Upgrade to V6</Label>
                              </Dropdown.Item>
                            )}
                            <Dropdown.Item
                              id="toggle"
                              textValue={n.enabled ? "Disable" : "Enable"}
                            >
                              <Label>{n.enabled ? "Disable" : "Enable"}</Label>
                            </Dropdown.Item>
                            <Dropdown.Item id="edit" textValue="Edit">
                              <Label>Edit</Label>
                            </Dropdown.Item>
                            <Dropdown.Item
                              id="delete"
                              textValue="Delete"
                              variant="danger"
                            >
                              <Label>Delete</Label>
                            </Dropdown.Item>
                          </Dropdown.Menu>
                        </Dropdown.Popover>
                      </Dropdown>
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
