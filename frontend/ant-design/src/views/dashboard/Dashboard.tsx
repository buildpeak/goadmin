import { useEffect, useState } from "react";
import { Card, Col, Row, Statistic, Table, Typography } from "antd";
import { UserOutlined, TeamOutlined, FileProtectOutlined } from "@ant-design/icons";
import { getUsers } from "../../services/backend-api";
import type { User } from "../../services/types";

const { Title } = Typography;

function Dashboard() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getUsers()
      .then(setUsers)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const activeUsers = users.filter((u) => u.active).length;

  const columns = [
    { title: "Username", dataIndex: "username", key: "username" },
    { title: "Email", dataIndex: "email", key: "email" },
    { title: "Name", key: "name", render: (_: unknown, r: User) => `${r.first_name} ${r.last_name}` },
    { title: "Active", dataIndex: "active", key: "active", render: (v: boolean) => (v ? "Yes" : "No") },
  ];

  return (
    <div>
      <Title level={3}>Dashboard</Title>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={8}>
          <Card>
            <Statistic title="Total Users" value={users.length} prefix={<TeamOutlined />} loading={loading} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="Active Users" value={activeUsers} prefix={<UserOutlined />} loading={loading} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="Roles" value={3} prefix={<FileProtectOutlined />} />
          </Card>
        </Col>
      </Row>
      <Card title="Recent Users">
        <Table dataSource={users} columns={columns} rowKey="id" loading={loading} size="small" />
      </Card>
    </div>
  );
}

export default Dashboard;
