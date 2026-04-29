import React, { useEffect, useState } from "react";
import { Typography, Card, Tag, Statistic, Row, Col, List, Empty, Spin } from "antd";
import { getMyAssets } from "../api";
import type { Asset } from "../types";

const { Title, Text } = Typography;

const MyAssetsPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => { getMyAssets().then((r) => setAssets(r.data)).finally(() => setLoading(false)); }, []);

  if (loading) return <Spin style={{ display: "block", margin: "80px auto" }} />;

  const totalValue = assets.reduce((a, b) => a + b.currentValue, 0);

  return (
    <>
      <Title level={4}>Миний эд хөрөнгө</Title>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}><Card><Statistic title="Миний хариуцсан хөрөнгө" value={assets.length} /></Card></Col>
        <Col span={12}><Card><Statistic title="Нийт үнэ цэнэ" value={totalValue} suffix="₮" /></Card></Col>
      </Row>
      {assets.length === 0 ? <Empty description="Хуваарилагдсан эд хөрөнгө байхгүй" /> : (
        <List dataSource={assets} renderItem={(a) => (
          <List.Item>
            <List.Item.Meta
              title={<><Text strong>{a.assetName}</Text> <Tag>{a.barcode}</Tag></>}
              description={<>Байршил: {a.location?.name || "—"} · Үнэ: {a.currentValue?.toLocaleString()} ₮ · Огноо: {a.acquisitionDate?.slice(0, 10)}</>}
            />
            <Tag color={a.status === "ACTIVE" ? "green" : "red"}>{a.status}</Tag>
          </List.Item>
        )} />
      )}
    </>
  );
};
export default MyAssetsPage;
