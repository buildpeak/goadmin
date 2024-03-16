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
import { verifyGoogleIdToken } from "../../services/google-auth";

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
    const idToken = response.credential;

    const backendAccessToken = await verifyGoogleIdToken(idToken);
    console.log(backendAccessToken);
  };

  useEffect(() => {
    const params = {
      clientId: process.env.REACT_APP_GOOGLE_CLIENT_ID,
      scope: "email",
    };

    let timer: number;

    const googleAccountInit = () => {
      if (window.google) {
        window.google.accounts.id.initialize({
          client_id: params.clientId,
          callback: googleSignInCallback,
        });

        window.google.accounts.id.renderButton(
          document.getElementById("googleSignInDiv"),
          { theme: "outline", size: "large" } // customization attributes
        );
        window.google.accounts.id.prompt(); // also display the One Tap dialog

        clearTimeout(timer);
      } else {
        timer = window.setTimeout(googleAccountInit, 100);
      }
    };

    googleAccountInit();

    return () => clearTimeout(timer);
  }, []);

  return (
    <Row justify="center" className="login-page">
      <Col xs={22} sm={22} md={22} lg={22} xl={22} xxl={22}>
        <Flex justify="space-evenly" vertical>
          <Flex justify="center">
            <Logo width={64} />
          </Flex>
          <Typography>
            <Title level={3}>Sign In</Title>
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

              <Form.Item style={{ marginBottom: 0 }}>
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
            <Divider>Or</Divider>
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
