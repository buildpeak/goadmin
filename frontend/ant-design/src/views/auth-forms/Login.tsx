import React, { useCallback, useEffect } from "react";
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
  message,
  Flex,
  Typography,
} from "antd";

import { useNavigate } from "react-router-dom";

import Logo from "../../components/Logo";
import "./Login.css";
import { login, verifyGoogleIdToken } from "../../services/backend-api";

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
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    try {
      const token = await login(values.username, values.password);

      localStorage.setItem("accessToken", token.access_token);
      localStorage.setItem("refreshToken", token.refresh_token);

      // redirect to dashboard
      navigate("/dashboard", { replace: true });
    } catch (error) {
      console.error(error);

      // show error message
      messageApi.open({
        type: "error",
        content: error instanceof Error ? error.message : "An error occurred",
      });
    }
  };

  const googleSignInCallback = useCallback(
    async (response: GoogleSignInResponse) => {
      const idToken = response.credential;

      localStorage.setItem("googleIdToken", idToken);

      try {
        const backendJwtToken = await verifyGoogleIdToken(idToken);

        localStorage.setItem("accessToken", backendJwtToken.access_token);
        localStorage.setItem("refreshToken", backendJwtToken.refresh_token);

        // redirect to dashboard
        navigate("/dashboard", { replace: true });
      } catch (error) {
        console.error(error);

        // show error message
        messageApi.open({
          type: "error",
          content: error instanceof Error ? error.message : "An error occurred",
        });
      }
    },
    [messageApi, navigate]
  );

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
  }, [googleSignInCallback]);

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
              {contextHolder}
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
                Or <a href="./signup">Sign up now!</a>
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
