import { Flex, Layout } from "antd";
import { Outlet } from "react-router-dom";

import { default as GFooter } from "../../components/Footer";

const { Footer, Content } = Layout;

function SimpleLayout() {
  return (
    <Flex vertical justify="space-between" style={{ minHeight: "100vh" }}>
      <Content>
        <Outlet />
      </Content>
      <Footer>
        <GFooter />
      </Footer>
    </Flex>
  );
}

export default SimpleLayout;

// I have read the <a href="/terms.txt">agreement</a>
