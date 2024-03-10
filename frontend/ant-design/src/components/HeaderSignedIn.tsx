import React from "react";
import { UserOutlined } from "@ant-design/icons";
import { Avatar, Flex } from "antd";

const HeaderSignedIn: React.FC = () => {
  return (
    <>
      <Flex justify="flex-end">
        <Avatar icon={<UserOutlined />} />
      </Flex>
    </>
  );
};

export default HeaderSignedIn;
