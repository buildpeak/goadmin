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
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import type { ColumnsType } from "antd/es/table";
import {
  listPermissions,
  createPermission,
  updatePermission,
  deletePermission,
} from "../../services/backend-api";
import type { Permission } from "../../services/types";

const PermissionList: React.FC = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingPerm, setEditingPerm] = useState<Permission | null>(null);
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchPerms = useCallback(async () => {
    setLoading(true);
    try {
      const data = await listPermissions();
      setPermissions(data);
    } catch {
      messageApi.error("Failed to load permissions");
    } finally {
      setLoading(false);
    }
  }, [messageApi]);

  useEffect(() => {
    fetchPerms();
  }, [fetchPerms]);

  const openCreate = () => {
    setEditingPerm(null);
    form.resetFields();
    setModalOpen(true);
  };

  const openEdit = (perm: Permission) => {
    setEditingPerm(perm);
    form.setFieldsValue(perm);
    setModalOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    try {
      if (editingPerm) {
        await updatePermission(editingPerm.id, values);
        messageApi.success("Permission updated");
      } else {
        await createPermission(values);
        messageApi.success("Permission created");
      }
      setModalOpen(false);
      fetchPerms();
    } catch {
      messageApi.error("Failed to save permission");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deletePermission(id);
      messageApi.success("Permission deleted");
      fetchPerms();
    } catch {
      messageApi.error("Failed to delete permission");
    }
  };

  const columns: ColumnsType<Permission> = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    { title: "Name", dataIndex: "name", key: "name" },
    {
      title: "Rule Type",
      dataIndex: "rule_type",
      key: "rule_type",
      width: 120,
      render: (val: string) => val ? <Tag>{val}</Tag> : "-",
    },
    { title: "Rule", dataIndex: "rule", key: "rule", ellipsis: true },
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
            title="Delete this permission?"
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
          Create Permission
        </Button>
      </Space>
      <Table
        columns={columns}
        dataSource={permissions}
        rowKey="id"
        loading={loading}
        size="middle"
      />

      <Modal
        title={editingPerm ? "Edit Permission" : "Create Permission"}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: "Please enter a name" }]}
          >
            <Input placeholder="e.g. users.read" />
          </Form.Item>
          <Form.Item name="rule_type" label="Rule Type">
            <Input placeholder="e.g. path, policy" />
          </Form.Item>
          <Form.Item name="rule" label="Rule">
            <Input placeholder="e.g. /v1/users/**" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default PermissionList;
