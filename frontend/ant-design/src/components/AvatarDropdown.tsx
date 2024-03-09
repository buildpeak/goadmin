import { LogoutOutlined, ProfileOutlined } from "@ant-design/icons";
import React from "react";

import { Dropdown } from "antd";

const AvatarDropdown: React.FC = (_props, dom) => (
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
    }}
  >
    {dom}
  </Dropdown>
);

export default AvatarDropdown;
