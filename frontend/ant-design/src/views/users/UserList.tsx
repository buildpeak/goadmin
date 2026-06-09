import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Table,
  Tag,
  message,
  Modal,
  Input,
  Form,
  Space,
  Transfer,
} from "antd";
import { PlusOutlined, EditOutlined, SafetyOutlined } from "@ant-design/icons";
import type { ColumnsType } from "antd/es/table";
import {
  getUsers,
  updateUser,
  signUp,
  listRoles,
  listRolesByUser,
  assignRoleToUser,
  unassignRoleFromUser,
} from "../../services/backend-api";
import type { User, Role } from "../../services/types";

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [rolesOpen, setRolesOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [allRoles, setAllRoles] = useState<Role[]>([]);
  const [assignedRoles, setAssignedRoles] = useState<string[]>([]);
  const [rolesLoading, setRolesLoading] = useState(false);
  const [editForm] = Form.useForm();
  const [createForm] = Form.useForm();

  const fetchUsers = useCallback(() => {
    setLoading(true);
    getUsers()
      .then(setUsers)
      .catch(() => message.error("Failed to load users"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleEdit = (user: User) => {
    setEditUser(user);
    editForm.setFieldsValue(user);
  };

  const handleSaveEdit = async () => {
    if (!editUser) return;
    setSubmitting(true);
    try {
      const values = await editForm.validateFields();
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

  const handleCreate = async () => {
    setSubmitting(true);
    try {
      const values = await createForm.validateFields();
      await signUp({ ...values, agreement: true });
      message.success("User created");
      setCreateOpen(false);
      createForm.resetFields();
      fetchUsers();
    } catch {
      // validation or API error
    } finally {
      setSubmitting(false);
    }
  };

  const openRoleManager = async (user: User) => {
    setSelectedUser(user);
    setRolesOpen(true);
    setRolesLoading(true);
    try {
      const [all, assigned] = await Promise.all([
        listRoles(),
        listRolesByUser(user.id),
      ]);
      setAllRoles(all);
      setAssignedRoles(assigned.map((r) => r.id));
    } catch {
      message.error("Failed to load roles");
    } finally {
      setRolesLoading(false);
    }
  };

  const handleRoleChange = async (targetKeys: string[]) => {
    if (!selectedUser) return;
    const toAdd = targetKeys.filter((k) => !assignedRoles.includes(k));
    const toRemove = assignedRoles.filter((k) => !targetKeys.includes(k));
    try {
      await Promise.all([
        ...toAdd.map((roleId) => assignRoleToUser(selectedUser.id, roleId)),
        ...toRemove.map((roleId) => unassignRoleFromUser(selectedUser.id, roleId)),
      ]);
      setAssignedRoles(targetKeys);
      message.success("Roles updated");
    } catch {
      message.error("Failed to update roles");
    }
  };

  const columns: ColumnsType<User> = [
    { title: "ID", dataIndex: "id", key: "id", width: 60 },
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
      width: 100,
      render: (v: boolean) =>
        v ? <Tag color="green">Active</Tag> : <Tag color="red">Inactive</Tag>,
    },
    {
      title: "Actions",
      key: "actions",
      width: 220,
      render: (_: unknown, r: User) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(r)}>
            Edit
          </Button>
          <Button
            icon={<SafetyOutlined />}
            size="small"
            onClick={() => openRoleManager(r)}
          >
            Roles
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Space style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <span style={{ fontSize: 18, fontWeight: 600 }}>Users</span>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
          Add User
        </Button>
      </Space>
      <Table dataSource={users} columns={columns} rowKey="id" loading={loading} />

      {/* Edit User Modal */}
      <Modal
        title="Edit User"
        open={!!editUser}
        onOk={handleSaveEdit}
        onCancel={() => setEditUser(null)}
        confirmLoading={submitting}
        destroyOnClose
      >
        <Form form={editForm} layout="vertical" style={{ marginTop: 16 }}>
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

      {/* Create User Modal */}
      <Modal
        title="Create User"
        open={createOpen}
        onOk={handleCreate}
        onCancel={() => { setCreateOpen(false); createForm.resetFields(); }}
        confirmLoading={submitting}
        destroyOnClose
      >
        <Form form={createForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="username"
            label="Username"
            rules={[{ required: true, message: "Required" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="email"
            label="Email"
            rules={[{ required: true, type: "email", message: "Valid email required" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="password"
            label="Password"
            rules={[{ required: true, min: 6, message: "Min 6 characters" }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item name="first_name" label="First Name">
            <Input />
          </Form.Item>
          <Form.Item name="last_name" label="Last Name">
            <Input />
          </Form.Item>
        </Form>
      </Modal>

      {/* Role Manager Modal */}
      <Modal
        title={`Roles — ${selectedUser?.username || ""}`}
        open={rolesOpen}
        onCancel={() => setRolesOpen(false)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <Transfer
          dataSource={allRoles.map((r) => ({ key: r.id, title: r.name }))}
          targetKeys={assignedRoles}
          onChange={handleRoleChange}
          render={(item) => item.title}
          listStyle={{ width: 250, height: 400 }}
          titles={["Available", "Assigned"]}
          showSearch
          disabled={rolesLoading}
        />
      </Modal>
    </div>
  );
};

export default UserList;
