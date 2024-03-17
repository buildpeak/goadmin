import React, { useEffect, useState } from "react";

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

const { Title } = Typography;

const SignUpForm: React.FC = () => {
  const [form] = Form.useForm();

  const [messageApi, contextHolder] = message.useMessage();

  const formItemLayout = {
    labelCol: {
      xs: { span: 24 },
      sm: { span: 8 },
    },
    wrapperCol: {
      xs: { span: 24 },
      sm: { span: 16 },
    },
  };
  const tailFormItemLayout = {};

  const onFinish = async (values: any) => {
    console.log(values);
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
            <Card>
              <Form
                {...formItemLayout}
                form={form}
                name="register"
                className="signup-form"
                onFinish={onFinish}
                scrollToFirstError
              >
                {contextHolder}
                <Form.Item
                  name="email"
                  label="E-mail"
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
                  <Input />
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
                  <Input.Password />
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
                  <Input.Password />
                </Form.Item>
                <Form.Item
                  name="full_name"
                  label="Full Name"
                  tooltip="Your full name."
                  rules={[
                    {
                      required: true,
                      message: "Please input your full name!",
                      whitespace: true,
                    },
                  ]}
                >
                  <Input />
                </Form.Item>

                <Form.Item
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
                  {...tailFormItemLayout}
                >
                  <Checkbox>
                    I have read the <a href="/terms.txt">agreement</a>
                  </Checkbox>
                </Form.Item>

                <Form.Item {...tailFormItemLayout}>
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
