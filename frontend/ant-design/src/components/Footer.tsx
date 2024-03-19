import React from "react";
import { Flex } from "antd";

const Footer: React.FC = () => {
  return (
    <Flex justify="center">
      Goadmin ©{new Date().getFullYear()} Powered by Gadmin
    </Flex>
  );
};

export default Footer;
