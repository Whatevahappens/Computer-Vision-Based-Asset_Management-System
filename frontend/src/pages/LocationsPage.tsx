import React, { useEffect, useState } from "react";
import { Typography, Table, Button, Modal, Form, Input, InputNumber, message, Popconfirm } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { listLocations, createLocation, deleteLocation } from "../api";
import type { Location } from "../types";

const { Title } = Typography;

const LocationsPage: React.FC = () => {
  const [locs, setLocs] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  const load = () => { setLoading(true); listLocations().then((r) => setLocs(r.data)).finally(() => setLoading(false)); };
  useEffect(() => { load(); }, []);

  const handleCreate = async (vals: any) => {
    try { await createLocation(vals); message.success("Байршил үүсгэлээ"); setModalOpen(false); form.resetFields(); load(); }
    catch { message.error("Алдаа"); }
  };

  const columns = [
    { title: "Нэр", dataIndex: "name" },
    { title: "Барилга", dataIndex: "building" },
    { title: "Давхар", dataIndex: "floor" },
    { title: "Өрөө", dataIndex: "room" },
    { title: "Багтаамж", dataIndex: "capacity" },
    { title: "", render: (_: unknown, r: Location) => <Popconfirm title="Устгах?" onConfirm={() => deleteLocation(r.id).then(load)}><Button size="small" danger>Устгах</Button></Popconfirm> },
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Байршил</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>Шинэ байршил</Button>
      </div>
      <Table rowKey="id" columns={columns} dataSource={locs} loading={loading} size="middle" />
      <Modal title="Шинэ байршил" open={modalOpen} onCancel={() => setModalOpen(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="Нэр" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="building" label="Барилга"><Input /></Form.Item>
          <Form.Item name="floor" label="Давхар"><Input /></Form.Item>
          <Form.Item name="room" label="Өрөө"><Input /></Form.Item>
          <Form.Item name="capacity" label="Багтаамж"><InputNumber min={0} style={{ width: "100%" }} /></Form.Item>
          <Button type="primary" htmlType="submit" block>Үүсгэх</Button>
        </Form>
      </Modal>
    </>
  );
};
export default LocationsPage;
