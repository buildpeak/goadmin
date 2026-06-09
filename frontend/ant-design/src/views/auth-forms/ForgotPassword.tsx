import React from "react";
import { MailOutlined } from "@ant-design/icons";
import { Button, Card, Col, Form, Input, message, Row, Flex, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import Logo from "../../components/Logo";
import { forgotPassword } from "../../services/backend-api";

const { Title, Text, Link } = Typography;

const ForgotPassword: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const [submitted, setSubmitted] = React.useState(false);
  const navigate = useNavigate();

  const onFinish = async (values: { email: string }) => {
    try {
      await forgotPassword(values.email);
      setSubmitted(true);
    } catch (error) {
      messageApi.open({
        type: "error",
        content: error instanceof Error ? error.message : "An error occurred",
      });
    }
  };

  if (submitted) {
    return (
      <Row justify="center" style={{ marginTop: 80 }}>
        <Col xs={22} sm={12} md={10} lg={8}>
          <Flex justify="space-evenly" vertical>
            <Flex justify="center">
              <Logo width={64} />
            </Flex>
            <Typography>
              <Title level={3}>Check your email</Title>
            </Typography>
            <Card>
              <Text>
                If an account with that email exists, a password reset link has been sent.
                Please check your inbox.
              </Text>
              <div style={{ marginTop: 16 }}>
                <Link onClick={() => navigate("/login")}>Back to Sign In</Link>
              </div>
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
            <Title level={3}>Forgot Password</Title>
          </Typography>
          <Card>
            {contextHolder}
            <Form name="forgot_password" onFinish={onFinish} layout="vertical">
              <Form.Item
                name="email"
                rules={[
                  { required: true, message: "Please input your Email!" },
                  { type: "email", message: "Please enter a valid email" },
                ]}
              >
                <Input
                  prefix={<MailOutlined />}
                  placeholder="Email"
                  size="large"
                />
              </Form.Item>
              <Form.Item style={{ marginBottom: 0 }}>
                <Button type="primary" htmlType="submit" block size="large">
                  Send Reset Link
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

export default ForgotPassword;
