import React, { useCallback, useEffect, useState } from "react";
import {
  Button,
  Table,
  Space,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Tag,
  Transfer,
} from "antd";
import { PlusOutlined, SettingOutlined } from "@ant-design/icons";
import type { ColumnsType } from "antd/es/table";
import type { TransferDirection } from "antd/es/transfer";
import {
  listRoles,
  createRole,
  updateRole,
  deleteRole,
  listPermissions,
  listPermissionsByRole,
  addPermissionToRole,
  removePermissionFromRole,
} from "../../services/backend-api";
import type { Role, Permission } from "../../services/types";

const RoleList: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [permModalOpen, setPermModalOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [allPerms, setAllPerms] = useState<Permission[]>([]);
  const [assignedPerms, setAssignedPerms] = useState<string[]>([]);
  const [permLoading, setPermLoading] = useState(false);
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchRoles = useCallback(async () => {
    setLoading(true);
    try {
      const data = await listRoles();
      setRoles(data);
    } catch {
      messageApi.error("Failed to load roles");
    } finally {
      setLoading(false);
    }
  }, [messageApi]);

  useEffect(() => {
    fetchRoles();
  }, [fetchRoles]);

  const openCreate = () => {
    setEditingRole(null);
    form.resetFields();
    setModalOpen(true);
  };

  const openEdit = (role: Role) => {
    setEditingRole(role);
    form.setFieldsValue({ name: role.name });
    setModalOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    try {
      if (editingRole) {
        await updateRole(editingRole.id, values);
        messageApi.success("Role updated");
      } else {
        await createRole(values);
        messageApi.success("Role created");
      }
      setModalOpen(false);
      fetchRoles();
    } catch {
      messageApi.error("Failed to save role");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteRole(id);
      messageApi.success("Role deleted");
      fetchRoles();
    } catch {
      messageApi.error("Failed to delete role");
    }
  };

  const openPermManager = async (role: Role) => {
    setSelectedRole(role);
    setPermModalOpen(true);
    setPermLoading(true);
    try {
      const [all, assigned] = await Promise.all([
        listPermissions(),
        listPermissionsByRole(role.id),
      ]);
      setAllPerms(all);
      setAssignedPerms(assigned.map((p) => p.id));
    } catch {
      messageApi.error("Failed to load permissions");
    } finally {
      setPermLoading(false);
    }
  };

  const handlePermChange = async (
    targetKeys: string[],
    direction: TransferDirection,
    moveKeys: string[]
  ) => {
    if (!selectedRole) return;
    try {
      if (direction === "right") {
        await Promise.all(
          moveKeys.map((permId) =>
            addPermissionToRole(selectedRole.id, permId)
          )
        );
      } else {
        await Promise.all(
          moveKeys.map((permId) =>
            removePermissionFromRole(selectedRole.id, permId)
          )
        );
      }
      setAssignedPerms(targetKeys);
      messageApi.success("Permissions updated");
    } catch {
      messageApi.error("Failed to update permissions");
    }
  };

  const columns: ColumnsType<Role> = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    { title: "Name", dataIndex: "name", key: "name" },
    {
      title: "Permissions",
      key: "perms",
      render: (_, record) => (
        <Button
          size="small"
          icon={<SettingOutlined />}
          onClick={() => openPermManager(record)}
        >
          Manage
        </Button>
      ),
    },
    { title: "Created", dataIndex: "created_at", key: "created_at", width: 180 },
    {
      title: "Actions",
      key: "actions",
      width: 160,
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => openEdit(record)}>
            Edit
          </Button>
          <Popconfirm
            title="Delete this role?"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button size="small" danger>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      {contextHolder}
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          Create Role
        </Button>
      </Space>
      <Table
        columns={columns}
        dataSource={roles}
        rowKey="id"
        loading={loading}
        size="middle"
      />

      <Modal
        title={editingRole ? "Edit Role" : "Create Role"}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="Role Name"
            rules={[{ required: true, message: "Please enter a role name" }]}
          >
            <Input placeholder="e.g. admin, editor" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`Permissions — ${selectedRole?.name || ""}`}
        open={permModalOpen}
        onCancel={() => setPermModalOpen(false)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <Transfer
          dataSource={allPerms.map((p) => ({
            key: p.id,
            title: p.name,
            description: p.rule_type || "",
          }))}
          targetKeys={assignedPerms}
          onChange={handlePermChange}
          render={(item) => item.title}
          listStyle={{ width: 250, height: 400 }}
          titles={["Available", "Assigned"]}
          showSearch
          disabled={permLoading}
        />
      </Modal>
    </>
  );
};

export default RoleList;
