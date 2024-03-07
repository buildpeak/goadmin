import React, { useEffect } from "react";
import { LockOutlined, UserOutlined } from "@ant-design/icons";
import {
  Button,
  Card,
  Checkbox,
  Col,
  Divider,
  Form,
  Row,
  Input,
  Flex,
  Typography,
} from "antd";

import Logo from "../../components/Logo";
import "./Login.css";

const { Title } = Typography;

declare global {
  interface Window {
    google: any;
  }
}

type GoogleSignInResponse = {
  credential: string;
};

const LoginForm: React.FC = () => {
  const onFinish = (values: any) => {
    console.log("Received values of form: ", values);
  };

  const googleSignInCallback = async (response: GoogleSignInResponse) => {
    console.log("Encoded JWT ID token: " + response.credential);
  };

  useEffect(() => {
    const params = {
      clientId:
        "33550892324-ct318rvjim5q61i846nsg726tb5vo4jm.apps.googleusercontent.com",
      scope: "email",
    };

    window.google.accounts.id.initialize({
      client_id: params.clientId,
      callback: googleSignInCallback,
    });

    window.google.accounts.id.renderButton(
      document.getElementById("googleSignInDiv"),
      { theme: "outline", size: "large" } // customization attributes
    );
    window.google.accounts.id.prompt(); // also display the One Tap dialog
  }, []);

  return (
    <Row justify="center" className="login-page">
      <Col xs={22} sm={20} md={12} lg={8} xl={6} xxl={4}>
        <Flex justify="space-evenly" vertical>
          <Flex justify="center">
            <Logo width={64} />
          </Flex>
          <Typography>
            <Title level={2}>Sign In</Title>
          </Typography>
          <Card>
            <Form
              name="normal_login"
              className="login-form"
              initialValues={{ remember: true }}
              onFinish={onFinish}
            >
              <Form.Item
                name="username"
                rules={[
                  { required: true, message: "Please input your Username!" },
                ]}
              >
                <Input
                  prefix={<UserOutlined className="site-form-item-icon" />}
                  placeholder="Username"
                />
              </Form.Item>
              <Form.Item
                name="password"
                rules={[
                  { required: true, message: "Please input your Password!" },
                ]}
              >
                <Input
                  prefix={<LockOutlined className="site-form-item-icon" />}
                  type="password"
                  placeholder="Password"
                />
              </Form.Item>
              <Form.Item>
                <Form.Item name="remember" valuePropName="checked" noStyle>
                  <Checkbox className="login-form-remember">
                    Remember me
                  </Checkbox>
                </Form.Item>

                <a className="login-form-forgot" href="./#">
                  Forgot password
                </a>
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  className="login-form-button"
                >
                  Log in
                </Button>
                Or <a href="./#">register now!</a>
              </Form.Item>
            </Form>
          </Card>
          <Divider>Or</Divider>
          <Card>
            <Flex justify="center">
              <div id="googleSignInDiv"></div>
            </Flex>
          </Card>
        </Flex>
      </Col>
    </Row>
  );
};

export default LoginForm;
