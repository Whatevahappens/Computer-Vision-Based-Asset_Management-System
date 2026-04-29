import React, { useEffect, useState } from "react";
import { Typography, Table, Button, Modal, Select, Radio, message, Statistic, Row, Col, Card } from "antd";
import { listAssets, calculateDepreciation, revalueAsset } from "../api";
import type { Asset } from "../types";

const { Title, Text } = Typography;

const DepreciationPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(false);
  const [calcModal, setCalcModal] = useState(false);
  const [selectedAsset, setSelectedAsset] = useState<string>("");
  const [method, setMethod] = useState<string>("STRAIGHT_LINE");

  const load = () => {
    setLoading(true);
    listAssets({ limit: 200, status: "ACTIVE" }).then((r) => setAssets(r.data.data || [])).finally(() => setLoading(false));
  };
  useEffect(() => { load(); }, []);

  const totalOriginal = assets.reduce((a, b) => a + b.acquisitionPrice, 0);
  const totalCurrent = assets.reduce((a, b) => a + b.currentValue, 0);
  const totalDepreciated = totalOriginal - totalCurrent;

  const handleCalc = async () => {
    if (!selectedAsset) { message.warning("Эд хөрөнгө сонгоно уу"); return; }
    try {
      await calculateDepreciation(selectedAsset, method);
      message.success("Элэгдэл тооцоологдлоо");
      setCalcModal(false); load();
    } catch (e: any) { message.error(e.response?.data?.error || "Алдаа"); }
  };

  const columns = [
    { title: "Эд хөрөнгө", render: (_: unknown, r: Asset) => <div><Text strong>{r.assetName}</Text><br/><Text type="secondary" style={{ fontSize: 11 }}>{r.barcode}</Text></div> },
    { title: "Анхны өртөг", dataIndex: "acquisitionPrice", render: (v: number) => `${v?.toLocaleString()} ₮` },
    { title: "Ашиглалтын жил", dataIndex: "usefulLifeMonths", render: (v: number) => `${Math.round(v / 12)} жил` },
    { title: "Хуримтлагдсан", render: (_: unknown, r: Asset) => `${(r.acquisitionPrice - r.currentValue).toLocaleString()} ₮` },
    { title: "Үлдэгдэл", dataIndex: "currentValue", render: (v: number) => `${v?.toLocaleString()} ₮` },
  ];

  return (
    <>
      <Title level={4}>Элэгдэл тооцоолол</Title>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}><Card><Statistic title="Элэгдэлтэй хөрөнгө" value={assets.length} /></Card></Col>
        <Col span={8}><Card><Statistic title="Нийт хуримтлагдсан" value={totalDepreciated} suffix="₮" /></Card></Col>
        <Col span={8}><Card><Statistic title="Нийт үлдэгдэл үнэ" value={totalCurrent} suffix="₮" /></Card></Col>
      </Row>
      <Button type="primary" style={{ marginBottom: 12 }} onClick={() => setCalcModal(true)}>Элэгдэл тооцоолох</Button>
      <Table rowKey="id" columns={columns} dataSource={assets} loading={loading} size="middle" />
      <Modal title="Элэгдэл тооцоолох" open={calcModal} onCancel={() => setCalcModal(false)} onOk={handleCalc} okText="Тооцоолох">
        <Select style={{ width: "100%", marginBottom: 12 }} placeholder="Эд хөрөнгө сонгох" value={selectedAsset || undefined} onChange={setSelectedAsset}
          options={assets.map((a) => ({ value: a.id, label: `${a.assetName} (${a.barcode})` }))} />
        <Radio.Group value={method} onChange={(e) => setMethod(e.target.value)}>
          <Radio value="STRAIGHT_LINE">Шулуун шугамын арга</Radio>
          <Radio value="DECLINING_BALANCE">Буурах үлдэгдлийн арга</Radio>
        </Radio.Group>
      </Modal>
    </>
  );
};
export default DepreciationPage;
