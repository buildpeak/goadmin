import React from "react";
import type { NavigateFunction } from "react-router-dom";
import type { MessageInstance } from "antd/lib/message/interface";
import { LogoutOutlined, ProfileOutlined } from "@ant-design/icons";

import { Dropdown } from "antd";
import { logout } from "../services/backend-api";

const AvatarDropdown = (
  props: {
    messageApi: MessageInstance;
    navigate: NavigateFunction;
  },
  dom: React.ReactNode
) => {
  const handleLogout = async () => {
    try {
      const accessToken = localStorage.getItem("accessToken");

      if (accessToken) {
        await logout(accessToken);

        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
        localStorage.removeItem("googleIdToken");
      }

      // go to login
      props.navigate("/login", { replace: true });
    } catch (error) {
      console.error(error);

      // show error message
      props.messageApi.open({
        type: "error",
        content: error instanceof Error ? error.message : "An error occurred",
      });
    }
  };

  const handleMenuClick = async (key: string) => {
    if (key === "profile") {
      console.log("profile");
    }

    if (key === "logout") {
      await handleLogout();
    }
  };

  return (
    <Dropdown
      menu={{
        items: [
          {
            key: "profile",
            icon: <ProfileOutlined />,
            label: "Profile",
          },
          {
            key: "logout",
            icon: <LogoutOutlined />,
            label: "Logout",
          },
        ],
        onClick: ({ key }) => handleMenuClick(key),
      }}
    >
      {dom}
    </Dropdown>
  );
};

export default AvatarDropdown;
