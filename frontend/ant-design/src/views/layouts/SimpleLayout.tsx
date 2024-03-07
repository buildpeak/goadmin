import { Footer, Content, Header } from "antd/es/layout/layout";
import { Outlet } from "react-router-dom";

function SimpleLayout() {
  return (
    <div>
      <Header />
      <Content>
        <Outlet />
      </Content>
      <Footer>@2024 GoAdmin</Footer>
    </div>
  );
}

export default SimpleLayout;
