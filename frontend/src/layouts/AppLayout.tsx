import React, { useState, useEffect } from "react";
import { Layout, Menu, Avatar, Badge, Dropdown, Typography, theme } from "antd";
import {
  DashboardOutlined, AppstoreOutlined, AuditOutlined, SwapOutlined,
  FileTextOutlined, EnvironmentOutlined, TeamOutlined, BankOutlined,
  BellOutlined, UserOutlined, LogoutOutlined, CalculatorOutlined,
  ToolOutlined, InboxOutlined,
} from "@ant-design/icons";
import { Outlet, useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { listNotifications } from "../api";
import type { Notification } from "../types";

const { Sider, Content, Header } = Layout;
const { Text } = Typography;

const ROLE_LABELS: Record<string, string> = {
  ADMIN: "Админ", ASSET_CUSTODIAN: "Нярав",
  ACCOUNTANT: "Нягтлан бодогч", EMPLOYEE: "Ажилтан",
};

const AppLayout: React.FC = () => {
  const { user, logout, hasRole } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const { token: themeToken } = theme.useToken();

  useEffect(() => {
    listNotifications().then((r) => setNotifications(r.data)).catch(() => {});
  }, [location.pathname]);

  const unread = notifications.filter((n) => !n.isRead).length;

  const menuItems = [
    { key: "/", icon: <DashboardOutlined />, label: "Хянах самбар" },
    hasRole("ASSET_CUSTODIAN", "ADMIN", "ACCOUNTANT") && {
      key: "/assets", icon: <AppstoreOutlined />, label: "Эд хөрөнгө",
    },
    hasRole("ASSET_CUSTODIAN", "ADMIN") && {
      key: "/audits", icon: <AuditOutlined />, label: "CV аудит",
    },
    hasRole("ASSET_CUSTODIAN", "ADMIN") && {
      key: "/assignments", icon: <SwapOutlined />, label: "Хуваарилалт",
    },
    hasRole("ACCOUNTANT", "ADMIN") && {
      key: "/depreciation", icon: <CalculatorOutlined />, label: "Элэгдэл",
    },
    hasRole("ASSET_CUSTODIAN", "ACCOUNTANT", "ADMIN") && {
      key: "/reports", icon: <FileTextOutlined />, label: "Тайлан",
    },
    hasRole("EMPLOYEE") && {
      key: "/my-assets", icon: <InboxOutlined />, label: "Миний эд хөрөнгө",
    },
    hasRole("ADMIN", "ASSET_CUSTODIAN") && {
      key: "/locations", icon: <EnvironmentOutlined />, label: "Байршил",
    },
    hasRole("ADMIN") && {
      key: "/departments", icon: <BankOutlined />, label: "Хэлтэс",
    },
    hasRole("ADMIN") && {
      key: "/asset-models", icon: <ToolOutlined />, label: "Загвар",
    },
    hasRole("ADMIN") && {
      key: "/users", icon: <TeamOutlined />, label: "Хэрэглэгч",
    },
  ].filter(Boolean) as { key: string; icon: React.ReactNode; label: string }[];

  const userMenu = {
    items: [
      { key: "role", label: ROLE_LABELS[user?.role || ""], disabled: true },
      { type: "divider" as const },
      { key: "logout", icon: <LogoutOutlined />, label: "Гарах", danger: true },
    ],
    onClick: ({ key }: { key: string }) => {
      if (key === "logout") { logout(); navigate("/login"); }
    },
  };

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        collapsible collapsed={collapsed} onCollapse={setCollapsed}
        style={{ background: "#fff", borderRight: "1px solid #f0f0f0" }}
        width={220}
      >
        <div style={{ padding: "20px 16px 12px", textAlign: collapsed ? "center" : "left" }}>
          <Text strong style={{ fontSize: collapsed ? 14 : 16, color: themeToken.colorPrimary }}>
            {collapsed ? "AS" : "Asset System"}
          </Text>
          {!collapsed && (
            <Text type="secondary" style={{ display: "block", fontSize: 11 }}>Prototype</Text>
          )}
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ border: "none" }}
        />
      </Sider>
      <Layout>
        <Header style={{
          background: "#fff", padding: "0 24px",
          display: "flex", justifyContent: "flex-end", alignItems: "center",
          gap: 16, borderBottom: "1px solid #f0f0f0", height: 56,
        }}>
          <Badge count={unread} size="small">
            <BellOutlined style={{ fontSize: 18, cursor: "pointer" }} onClick={() => navigate("/notifications")} />
          </Badge>
          <Dropdown menu={userMenu} placement="bottomRight">
            <div style={{ cursor: "pointer", display: "flex", alignItems: "center", gap: 8 }}>
              <Avatar size="small" icon={<UserOutlined />} style={{ background: themeToken.colorPrimary }} />
              <Text>{user?.lastName?.charAt(0)}.{user?.firstName}</Text>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: 16, padding: 20, background: "#fff", borderRadius: 8, minHeight: 400 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
