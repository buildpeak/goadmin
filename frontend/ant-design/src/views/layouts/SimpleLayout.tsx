import { Flex, Layout } from "antd";
import { Outlet } from "react-router-dom";

const { Footer, Content } = Layout;

function SimpleLayout() {
  return (
    <div>
      <Content>
        <Outlet />
      </Content>
      <Footer style={{ bottom: 0, position: "fixed", width: "100%" }}>
        <Flex justify="center">@2024 GoAdmin</Flex>
      </Footer>
    </div>
  );
}

export default SimpleLayout;
