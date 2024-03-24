import React from "react";
import { Button, Result } from "antd";
import { SignUpResultProps } from "./data-types";

const SignUpResult: React.FC<SignUpResultProps> = (
  props: SignUpResultProps
) => {
  return (
    <Result
      status="success"
      title="Sign up Successfully!"
      subTitle={props.username}
      extra={[
        <Button type="primary" key="console">
          Go to Login
        </Button>,
      ]}
    />
  );
};

export default SignUpResult;
