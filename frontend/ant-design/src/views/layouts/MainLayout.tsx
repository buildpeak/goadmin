import React from "react";

import { ProLayout, ProConfigProvider } from "@ant-design/pro-components";
import { ConfigProvider, message } from "antd";
import { Outlet, useNavigate } from "react-router-dom";
import Logo from "../../components/Logo";
import { UserOutlined } from "@ant-design/icons";
import AvatarDropdown from "../../components/AvatarDropdown";

const MainLayout: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

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
              render: (props, dom, _siderProps) => {
                return AvatarDropdown(
                  {
                    messageApi,
                    navigate,
                    ...props,
                  },
                  dom
                );
              },
            }}
            layout="mix"
          >
            {contextHolder}
            <Outlet />
          </ProLayout>
        </ConfigProvider>
      </ProConfigProvider>
    </div>
  );
};

export default MainLayout;
