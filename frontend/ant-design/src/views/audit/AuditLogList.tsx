import React, { useCallback, useEffect, useState } from "react";
import { Table, Tag, Select, Space, Input } from "antd";
import type { ColumnsType } from "antd/es/table";
import { listAuditLogs } from "../../services/backend-api";
import type { AuditLog } from "../../services/types";

const actionColors: Record<string, string> = {
  login: "blue",
  logout: "default",
  signup: "green",
  create: "cyan",
  update: "orange",
  delete: "red",
};

const AuditLogList: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [actionFilter, setActionFilter] = useState<string>("");
  const [actorFilter, setActorFilter] = useState<string>("");

  const fetchLogs = useCallback(async () => {
    setLoading(true);
    try {
      const data = await listAuditLogs({
        limit: 100,
        action: actionFilter || undefined,
        actor: actorFilter || undefined,
      });
      setLogs(data);
    } catch {
      // API error handled by interceptor
    } finally {
      setLoading(false);
    }
  }, [actionFilter, actorFilter]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  const columns: ColumnsType<AuditLog> = [
    { title: "ID", dataIndex: "id", key: "id", width: 60 },
    {
      title: "Action",
      dataIndex: "action",
      key: "action",
      width: 100,
      render: (v: string) => (
        <Tag color={actionColors[v] || "default"}>{v}</Tag>
      ),
    },
    { title: "Actor", dataIndex: "actor_username", key: "actor", width: 120 },
    { title: "Resource", dataIndex: "resource_type", key: "resource", width: 100 },
    { title: "Resource ID", dataIndex: "resource_id", key: "resource_id", width: 80 },
    { title: "Detail", dataIndex: "detail", key: "detail", ellipsis: true },
    { title: "IP", dataIndex: "ip_address", key: "ip", width: 130 },
    { title: "Timestamp", dataIndex: "created_at", key: "created_at", width: 180 },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: "flex", justifyContent: "space-between" }}>
        <span style={{ fontSize: 18, fontWeight: 600 }}>Audit Logs</span>
        <Space>
          <Input
            placeholder="Filter by actor"
            value={actorFilter}
            onChange={(e) => setActorFilter(e.target.value)}
            style={{ width: 160 }}
            allowClear
          />
          <Select
            placeholder="Filter by action"
            value={actionFilter || undefined}
            onChange={(v) => setActionFilter(v || "")}
            style={{ width: 140 }}
            allowClear
            options={[
              { value: "login", label: "Login" },
              { value: "logout", label: "Logout" },
              { value: "signup", label: "Signup" },
              { value: "create", label: "Create" },
              { value: "update", label: "Update" },
              { value: "delete", label: "Delete" },
            ]}
          />
        </Space>
      </Space>
      <Table
        columns={columns}
        dataSource={logs}
        rowKey="id"
        loading={loading}
        size="middle"
        scroll={{ x: 900 }}
      />
    </div>
  );
};

export default AuditLogList;
