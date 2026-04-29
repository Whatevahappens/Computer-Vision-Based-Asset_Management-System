import React, { useEffect, useState } from "react";
import { Typography, Table, Button, Modal, Form, Input, InputNumber, Select, message } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { listAssetModels, createAssetModel } from "../api";
import type { AssetModel } from "../types";

const { Title } = Typography;

const AssetModelsPage: React.FC = () => {
  const [models, setModels] = useState<AssetModel[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  const load = () => { setLoading(true); listAssetModels().then((r) => setModels(r.data)).finally(() => setLoading(false)); };
  useEffect(() => { load(); }, []);

  const handleCreate = async (vals: any) => {
    try { await createAssetModel(vals); message.success("Загвар үүсгэлээ"); setModalOpen(false); form.resetFields(); load(); }
    catch { message.error("Алдаа"); }
  };

  const columns = [
    { title: "Брэнд", dataIndex: "brand" },
    { title: "Загвар", dataIndex: "modelName" },
    { title: "Ангилал", dataIndex: "category" },
    { title: "Төрөл", dataIndex: "assetType" },
    { title: "Үнэ", dataIndex: "defaultUnitPrice", render: (v: number) => v ? `${v.toLocaleString()} ₮` : "—" },
    { title: "Элэгдлийн арга", dataIndex: "depreciationMethod" },
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Эд хөрөнгийн загвар</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>Шинэ загвар</Button>
      </div>
      <Table rowKey="id" columns={columns} dataSource={models} loading={loading} size="middle" />
      <Modal title="Шинэ загвар" open={modalOpen} onCancel={() => setModalOpen(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="brand" label="Брэнд" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="modelName" label="Загварын нэр" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="assetType" label="Төрөл" rules={[{ required: true }]}>
            <Select options={["EQUIPMENT","FURNITURE","VEHICLE","ELECTRONIC","OTHER"].map((v) => ({ value: v, label: v }))} />
          </Form.Item>
          <Form.Item name="category" label="Ангилал" rules={[{ required: true }]}>
            <Select options={["IT_EQUIPMENT","OFFICE_EQUIPMENT","FURNITURE","VEHICLE","OTHER"].map((v) => ({ value: v, label: v }))} />
          </Form.Item>
          <Form.Item name="defaultUnitPrice" label="Анхны үнэ"><InputNumber min={0} style={{ width: "100%" }} /></Form.Item>
          <Form.Item name="defaultUsefulLifeMonths" label="Ашиглалт (сар)"><InputNumber min={1} style={{ width: "100%" }} /></Form.Item>
          <Form.Item name="depreciationMethod" label="Элэгдлийн арга">
            <Select options={[{ value: "STRAIGHT_LINE", label: "Шулуун шугам" }, { value: "DECLINING_BALANCE", label: "Буурах үлдэгдэл" }]} />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>Үүсгэх</Button>
        </Form>
      </Modal>
    </>
  );
};
export default AssetModelsPage;
