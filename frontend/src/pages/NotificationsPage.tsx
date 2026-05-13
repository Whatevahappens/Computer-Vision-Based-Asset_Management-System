import React, { useEffect, useState } from "react";
import { Typography, List, Badge, Button, Empty, message } from "antd";
import { BellOutlined } from "@ant-design/icons";
import { listNotifications, markNotificationRead, markAllRead } from "../api";
import type { Notification } from "../types";
import dayjs from "dayjs";

const { Title, Text } = Typography;

const NotificationsPage: React.FC = () => {
  const [notifs, setNotifs] = useState<Notification[]>([]);
  useEffect(() => { listNotifications().then((r) => setNotifs(r.data)); }, []);

  const unread = notifs.filter((n) => !n.isRead).length;
  const handleMarkAll = async () => { await markAllRead(); setNotifs(notifs.map((n) => ({ ...n, isRead: true }))); message.success("Бүгд уншсан"); };
  const handleRead = async (id: string) => { await markNotificationRead(id); setNotifs(notifs.map((n) => n.id === id ? { ...n, isRead: true } : n)); };

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Мэдэгдэл ({unread} уншаагүй)</Title>
        {unread > 0 && <Button onClick={handleMarkAll}>Бүгдийг уншсан</Button>}
      </div>
      {notifs.length === 0 ? <Empty description="Мэдэгдэл байхгүй" /> : (
        <List dataSource={notifs} renderItem={(n) => (
          <List.Item actions={!n.isRead ? [<Button size="small" onClick={() => handleRead(n.id)}>Уншсан</Button>] : []}>
            <List.Item.Meta
              avatar={<Badge dot={!n.isRead}><BellOutlined style={{ fontSize: 20 }} /></Badge>}
              title={n.title}
              description={<><Text>{n.message}</Text><br/><Text type="secondary" style={{ fontSize: 11 }}>{dayjs(n.createdAt).format("YYYY-MM-DD HH:mm")}</Text></>}
            />
          </List.Item>
        )} />
      )}
    </>
  );
};
export default NotificationsPage;
