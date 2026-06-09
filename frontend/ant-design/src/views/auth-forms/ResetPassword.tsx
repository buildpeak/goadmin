import React from "react";
import { LockOutlined } from "@ant-design/icons";
import { Button, Card, Col, Form, Input, message, Row, Flex, Typography } from "antd";
import { useNavigate, useSearchParams } from "react-router-dom";
import Logo from "../../components/Logo";
import { resetPassword } from "../../services/backend-api";

const { Title, Link } = Typography;

const ResetPassword: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get("token") || "";

  const onFinish = async (values: { password: string; confirm: string }) => {
    if (values.password !== values.confirm) {
      messageApi.open({
        type: "error",
        content: "Passwords do not match",
      });
      return;
    }

    try {
      await resetPassword(token, values.password);
      messageApi.open({
        type: "success",
        content: "Password has been reset successfully.",
      });
      setTimeout(() => navigate("/login", { replace: true }), 1500);
    } catch (error) {
      messageApi.open({
        type: "error",
        content: error instanceof Error ? error.message : "Failed to reset password",
      });
    }
  };

  if (!token) {
    return (
      <Row justify="center" style={{ marginTop: 80 }}>
        <Col xs={22} sm={12} md={10} lg={8}>
          <Flex justify="space-evenly" vertical>
            <Flex justify="center">
              <Logo width={64} />
            </Flex>
            <Typography>
              <Title level={3}>Invalid Link</Title>
            </Typography>
            <Card>
              <p>This password reset link is invalid or has expired.</p>
              <Link onClick={() => navigate("/forgot-password")}>Request a new reset link</Link>
            </Card>
          </Flex>
        </Col>
      </Row>
    );
  }

  return (
    <Row justify="center" style={{ marginTop: 80 }}>
      <Col xs={22} sm={12} md={10} lg={8}>
        <Flex justify="space-evenly" vertical>
          <Flex justify="center">
            <Logo width={64} />
          </Flex>
          <Typography>
            <Title level={3}>Reset Password</Title>
          </Typography>
          <Card>
            {contextHolder}
            <Form name="reset_password" onFinish={onFinish} layout="vertical">
              <Form.Item
                name="password"
                rules={[
                  { required: true, message: "Please input your new password!" },
                  { min: 6, message: "Password must be at least 6 characters" },
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="New password"
                  size="large"
                />
              </Form.Item>
              <Form.Item
                name="confirm"
                dependencies={["password"]}
                rules={[
                  { required: true, message: "Please confirm your new password!" },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue("password") === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error("Passwords do not match"));
                    },
                  }),
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="Confirm new password"
                  size="large"
                />
              </Form.Item>
              <Form.Item style={{ marginBottom: 0 }}>
                <Button type="primary" htmlType="submit" block size="large">
                  Reset Password
                </Button>
              </Form.Item>
            </Form>
            <div style={{ marginTop: 16, textAlign: "center" }}>
              <Link onClick={() => navigate("/login")}>Back to Sign In</Link>
            </div>
          </Card>
        </Flex>
      </Col>
    </Row>
  );
};

export default ResetPassword;
