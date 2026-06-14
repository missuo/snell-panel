import { useMemo, useState } from "react";
import {
  Button,
  Card,
  Input,
  Label,
  ListBox,
  Select,
  Switch,
  TextField,
} from "@heroui/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { SettingsResponse, SubscriptionFormat } from "@snell-panel/shared";
import { api } from "../api/client";
import { CopyButton } from "./CopyButton";

export function SubscriptionCard() {
  const qc = useQueryClient();
  const settings = useQuery({ queryKey: ["settings"], queryFn: api.settings });

  const [format, setFormat] = useState<SubscriptionFormat>("surge");
  const [flag, setFlag] = useState(true);
  const [filter, setFilter] = useState("");
  const [via, setVia] = useState("");

  const reset = useMutation({
    mutationFn: api.resetSubscribeToken,
    onSuccess: (d: SettingsResponse) => qc.setQueryData(["settings"], d),
  });

  const url = useMemo(() => {
    if (!settings.data) return "";
    const p = new URLSearchParams();
    p.set("token", settings.data.subscribe_token);
    p.set("format", format);
    if (!flag) p.set("flag", "false");
    if (filter.trim()) p.set("filter", filter.trim());
    if (via.trim()) p.set("via", via.trim());
    return `${window.location.origin}/api/subscribe?${p.toString()}`;
  }, [settings.data, format, flag, filter, via]);

  return (
    <Card>
      <Card.Header>
        <Card.Title>Subscription</Card.Title>
        <Card.Description>
          Build a subscription URL. Its token is independent of your access token
          and can be reset anytime.
        </Card.Description>
      </Card.Header>

      <div className="flex flex-col gap-4 px-6 pb-2">
        <div className="flex flex-wrap items-end gap-3">
          <Select
            className="w-44"
            selectedKey={format}
            onSelectionChange={(k) => setFormat(String(k) as SubscriptionFormat)}
          >
            <Label>Format</Label>
            <Select.Trigger>
              <Select.Value />
              <Select.Indicator />
            </Select.Trigger>
            <Select.Popover>
              <ListBox>
                <ListBox.Item id="surge" textValue="Surge">
                  Surge
                  <ListBox.ItemIndicator />
                </ListBox.Item>
                <ListBox.Item id="shadowrocket" textValue="Shadowrocket">
                  Shadowrocket
                  <ListBox.ItemIndicator />
                </ListBox.Item>
                <ListBox.Item id="mihomo" textValue="Mihomo">
                  Mihomo
                  <ListBox.ItemIndicator />
                </ListBox.Item>
              </ListBox>
            </Select.Popover>
          </Select>

          <TextField value={filter} onChange={setFilter} className="w-40">
            <Label>Filter (name)</Label>
            <Input placeholder="optional" />
          </TextField>

          <TextField value={via} onChange={setVia} className="w-40">
            <Label>Via (relay)</Label>
            <Input placeholder="optional" />
          </TextField>

          <Switch isSelected={flag} onChange={setFlag}>
            <Switch.Control>
              <Switch.Thumb />
            </Switch.Control>
            <Switch.Content>
              <Label className="text-sm">Country flag</Label>
            </Switch.Content>
          </Switch>
        </div>

        <TextField value={url} className="w-full">
          <Label className="sr-only">Subscription URL</Label>
          <Input readOnly className="font-mono text-xs" />
        </TextField>
      </div>

      <Card.Footer className="gap-2">
        <CopyButton text={url} label="Copy URL" />
        <Button
          variant="danger"
          isDisabled={reset.isPending}
          onPress={() => {
            if (
              confirm(
                "Reset the subscribe token? Existing subscription URLs will stop working.",
              )
            ) {
              reset.mutate();
            }
          }}
        >
          {reset.isPending ? "Resetting…" : "Reset token"}
        </Button>
      </Card.Footer>
    </Card>
  );
}
