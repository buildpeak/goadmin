import React, { useEffect } from "react";

import {
  Button,
  Card,
  Checkbox,
  Col,
  Flex,
  Form,
  Input,
  message,
  Row,
  Typography,
} from "antd";
import Logo from "../../components/Logo";
import "./SignUp.css";
import { jwtDecode } from "jwt-decode";
import { GoogleJwtPayload } from "./data-types";

const { Title } = Typography;

const SignUpForm: React.FC = () => {
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const formItemLayout = {
    labelCol: { span: 8 },
    wrapperCol: { span: 16 },
  };

  useEffect(() => {
    const googleIdToken = localStorage.getItem("googleIdToken");
    if (!googleIdToken) {
      return;
    }

    const userInfo = jwtDecode(googleIdToken) as GoogleJwtPayload;

    form.setFieldsValue({
      username: userInfo.email,
      email: userInfo.email,
      full_name: userInfo.name,
    });

    form.getFieldInstance("email").input.disabled = true;
  }, [form]);

  const onFinish = async (values: any) => {
    console.log(JSON.stringify(values));

    messageApi.loading("Registering...");
    // try {
    //   const data = await register(values);
    //
    //   console.log(data);
    //
    //   messageApi.open({
    //     type: "success",
    //     content: "Registration successful",
    //   });
    // } catch (error) {
    //   console.error(error);
    //   if (error.response?.status === 400) {
    //     messageApi.open({
    //       type: "error",
    //       content: "Invalid input",
    //     });
    //     return;
    //   }
    //   messageApi.open({
    //     type: "error",
    //     content: error.message,
    //   });
    // }
  };

  return (
    <Row justify="center" className="signup-page">
      <Col span={22}>
        <Flex justify="space-evenly" vertical>
          <Flex justify="center">
            <Logo width={64} />
          </Flex>
          <Typography>
            <Title level={3}>Sign Up</Title>
          </Typography>
          <Flex justify="center">
            <Card className="signup-form">
              <Form
                {...formItemLayout}
                form={form}
                layout="horizontal"
                name="register"
                onFinish={onFinish}
                scrollToFirstError
              >
                {contextHolder}
                <Form.Item
                  name="username"
                  label="Username"
                  rules={[
                    {
                      required: true,
                      message: "Please input your username!",
                    },
                  ]}
                >
                  <Input placeholder="Username" />
                </Form.Item>

                <Form.Item
                  name="email"
                  label="Email"
                  rules={[
                    {
                      type: "email",
                      message: "The input is not valid E-mail!",
                    },
                    {
                      required: true,
                      message: "Please input your E-mail!",
                    },
                  ]}
                >
                  <Input placeholder="Email" />
                </Form.Item>

                <Form.Item
                  name="password"
                  label="Password"
                  rules={[
                    {
                      required: true,
                      message: "Please input your password!",
                    },
                  ]}
                  hasFeedback
                >
                  <Input.Password placeholder="Password" />
                </Form.Item>

                <Form.Item
                  name="confirm"
                  label="Confirm Password"
                  dependencies={["password"]}
                  hasFeedback
                  rules={[
                    {
                      required: true,
                      message: "Please confirm your password!",
                    },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue("password") === value) {
                          return Promise.resolve();
                        }
                        return Promise.reject(
                          new Error(
                            "The new password that you entered do not match!"
                          )
                        );
                      },
                    }),
                  ]}
                >
                  <Input.Password placeholder="Confirm Password" />
                </Form.Item>

                <Form.Item
                  name="full_name"
                  label="Full Name"
                  tooltip="Your full name: First Name, Last Name."
                  rules={[
                    {
                      required: true,
                      message: "Please input your full name!",
                      whitespace: true,
                    },
                  ]}
                >
                  <Input placeholder="Full Name" />
                </Form.Item>

                <Form.Item
                  wrapperCol={{ span: 24 }}
                  name="agreement"
                  valuePropName="checked"
                  rules={[
                    {
                      validator: (_, value) =>
                        value
                          ? Promise.resolve()
                          : Promise.reject(
                              new Error("Should accept agreement")
                            ),
                    },
                  ]}
                >
                  <Checkbox>
                    I have read the <a href="/terms.txt">agreement</a>
                  </Checkbox>
                </Form.Item>

                <Form.Item wrapperCol={{ span: 24 }}>
                  <Button
                    type="primary"
                    htmlType="submit"
                    className="signup-form-button"
                  >
                    Sign Up
                  </Button>
                  <div>
                    Or <a href="/login">Go back to login!</a>
                  </div>
                </Form.Item>
              </Form>
            </Card>
          </Flex>
        </Flex>
      </Col>
    </Row>
  );
};

export default SignUpForm;
