import React from "react";

import { ProLayout, ProConfigProvider } from "@ant-design/pro-components";
import { ConfigProvider } from "antd";
import { Outlet } from "react-router-dom";
import Logo from "../../components/Logo";
import { UserOutlined } from "@ant-design/icons";
import AvatarDropdown from "../../components/AvatarDropdown";

const MainLayout: React.FC = () => {
  return (
    <div id="main-layout">
      <ProConfigProvider hashed={false}>
        <ConfigProvider
          getTargetContainer={() => {
            return document.getElementById("main-layout") || document.body;
          }}
        >
          <ProLayout
            title="GoAdmin"
            logo={<Logo width={32} />}
            avatarProps={{
              size: "small",
              title: "User",
              icon: <UserOutlined />,
              render: AvatarDropdown,
            }}
            layout="mix"
          >
            <Outlet />
          </ProLayout>
        </ConfigProvider>
      </ProConfigProvider>
    </div>
  );
};

export default MainLayout;
