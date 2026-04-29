import React, { useEffect, useState } from "react";
import { Card, Col, Row, Statistic, Typography, Spin } from "antd";
import {
  AppstoreOutlined, CheckCircleOutlined, DollarOutlined, AuditOutlined,
} from "@ant-design/icons";
import { getDashboard } from "../api";
import type { DashboardStats } from "../types";

const { Title } = Typography;

const DashboardPage: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getDashboard()
      .then((r) => setStats(r.data))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin size="large" style={{ display: "block", margin: "100px auto" }} />;

  const cards = [
    { title: "Нийт эд хөрөнгө", value: stats?.totalAssets ?? 0, icon: <AppstoreOutlined />, color: "#1677ff" },
    { title: "Идэвхтэй", value: stats?.activeAssets ?? 0, icon: <CheckCircleOutlined />, color: "#52c41a" },
    { title: "Нийт үнэ цэнэ", value: stats?.totalValue ?? 0, icon: <DollarOutlined />, color: "#fa8c16", suffix: "₮", },
    { title: "Нийт аудит", value: stats?.totalAudits ?? 0, icon: <AuditOutlined />, color: "#722ed1" },
  ];

  return (
    <>
      <Title level={4}>Хянах самбар</Title>
      <Row gutter={[16, 16]}>
        {cards.map((c) => (
          <Col xs={24} sm={12} lg={6} key={c.title}>
            <Card hoverable style={{ borderTop: `3px solid ${c.color}` }}>
              <Statistic
                title={c.title}
                value={c.value}
                prefix={c.icon}
                suffix={c.suffix}
                valueStyle={{ color: c.color }}
              />
            </Card>
          </Col>
        ))}
      </Row>
    </>
  );
};

export default DashboardPage;
