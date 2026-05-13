import React, { useState } from "react";
import { Typography, Card, Row, Col, Button, message } from "antd";
import { FileTextOutlined, CalculatorOutlined, SwapOutlined, AuditOutlined } from "@ant-design/icons";
import { generateReport } from "../api";

const { Title } = Typography;

const reports = [
  { key: "asset_detail", title: "Эд хөрөнгийн дэлгэрэнгүй тайлан", desc: "Бүх бүртгэлтэй эд хөрөнгийн жагсаалт", icon: <FileTextOutlined style={{ fontSize: 28 }} /> },
  { key: "depreciation", title: "Элэгдлийн тайлан", desc: "Хуримтлагдсан элэгдлийн тайлан", icon: <CalculatorOutlined style={{ fontSize: 28 }} /> },
];

const ReportsPage: React.FC = () => {
  const [loading, setLoading] = useState<string | null>(null);

  const handleGenerate = async (type: string) => {
    setLoading(type);
    try {
      const { data } = await generateReport(type, "csv");
      message.success("Тайлан үүслээ");
      if (data.download) window.open(`/api/v1${data.download}`, "_blank");
    } catch { message.error("Алдаа"); } finally { setLoading(null); }
  };

  return (
    <>
      <Title level={4}>Тайлан</Title>
      <Row gutter={[16, 16]}>
        {reports.map((r) => (
          <Col xs={24} sm={12} lg={8} key={r.key}>
            <Card hoverable>
              <div style={{ textAlign: "center", marginBottom: 12 }}>{r.icon}</div>
              <Card.Meta title={r.title} description={r.desc} />
              <Button type="primary" block style={{ marginTop: 16 }} loading={loading === r.key} onClick={() => handleGenerate(r.key)}>
                Тайлан үүсгэх
              </Button>
            </Card>
          </Col>
        ))}
      </Row>
    </>
  );
};
export default ReportsPage;
