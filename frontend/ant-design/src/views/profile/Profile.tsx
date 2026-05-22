import { useEffect, useState } from "react";
import { Card, Descriptions, Typography, Spin, Tag } from "antd";
import { getProfile } from "../../services/backend-api";
import type { User } from "../../services/types";

const { Title } = Typography;

function Profile() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getProfile()
      .then(setUser)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: "center", padding: 48 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!user) return null;

  return (
    <div>
      <Title level={3}>Profile</Title>
      <Card>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="Username">{user.username}</Descriptions.Item>
          <Descriptions.Item label="Email">{user.email}</Descriptions.Item>
          <Descriptions.Item label="First Name">{user.first_name}</Descriptions.Item>
          <Descriptions.Item label="Last Name">{user.last_name}</Descriptions.Item>
          <Descriptions.Item label="Status">
            {user.active ? <Tag color="green">Active</Tag> : <Tag color="red">Inactive</Tag>}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
}

export default Profile;
