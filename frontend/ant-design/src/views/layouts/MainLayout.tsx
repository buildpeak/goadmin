import React, { useEffect } from "react";

import { ProLayout, ProConfigProvider } from "@ant-design/pro-components";
import { ConfigProvider, message } from "antd";
import { Outlet, useNavigate, useLocation } from "react-router-dom";
import Logo from "../../components/Logo";
import { UserOutlined, DashboardOutlined, TeamOutlined } from "@ant-design/icons";
import AvatarDropdown from "../../components/AvatarDropdown";

const MainLayout: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();
  const location = useLocation();

  // check login status
  useEffect(() => {
    const accessToken = localStorage.getItem("accessToken");
    if (!accessToken) {
      navigate("/login", { replace: true });
    }
  }, [navigate]);

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
            location={location}
            route={{
              routes: [
                { path: "/dashboard", name: "Dashboard", icon: <DashboardOutlined /> },
                { path: "/users", name: "Users", icon: <TeamOutlined /> },
              ],
            }}
            menuItemRender={(item, dom) => (
              <a onClick={() => item.path && navigate(item.path)}>
                {dom}
              </a>
            )}
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
