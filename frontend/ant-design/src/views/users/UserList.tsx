import { useEffect, useState } from "react";
import { Button, Table, Tag, Typography, message, Modal, Input, Form, Space } from "antd";
import { PlusOutlined, EditOutlined } from "@ant-design/icons";
import { getUsers, updateUser } from "../../services/backend-api";
import type { User } from "../../services/types";
import type { ColumnsType } from "antd/es/table";

const { Title } = Typography;

function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchUsers = () => {
    setLoading(true);
    getUsers()
      .then(setUsers)
      .catch(() => message.error("Failed to load users"))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleEdit = (user: User) => {
    setEditUser(user);
    form.setFieldsValue(user);
  };

  const handleSave = async () => {
    if (!editUser) return;
    setSubmitting(true);
    try {
      const values = await form.validateFields();
      await updateUser(editUser.id, values);
      message.success("User updated");
      setEditUser(null);
      fetchUsers();
    } catch {
      // validation failed or API error
    } finally {
      setSubmitting(false);
    }
  };

  const columns: ColumnsType<User> = [
    { title: "Username", dataIndex: "username", key: "username" },
    { title: "Email", dataIndex: "email", key: "email" },
    {
      title: "Name",
      key: "name",
      render: (_: unknown, r: User) => `${r.first_name} ${r.last_name}`,
    },
    {
      title: "Active",
      dataIndex: "active",
      key: "active",
      render: (v: boolean) => v ? <Tag color="green">Active</Tag> : <Tag color="red">Inactive</Tag>,
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: unknown, r: User) => (
        <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(r)}>
          Edit
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Space style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={3} style={{ margin: 0 }}>Users</Title>
        <Button type="primary" icon={<PlusOutlined />} disabled>
          Add User
        </Button>
      </Space>
      <Table dataSource={users} columns={columns} rowKey="id" loading={loading} />
      <Modal
        title="Edit User"
        open={!!editUser}
        onOk={handleSave}
        onCancel={() => setEditUser(null)}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="first_name" label="First Name">
            <Input />
          </Form.Item>
          <Form.Item name="last_name" label="Last Name">
            <Input />
          </Form.Item>
          <Form.Item name="email" label="Email" rules={[{ type: "email" }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

export default UserList;
