import React, { useEffect, useState } from "react";
import { Typography, Table, Button, Modal, Form, Input, message, Popconfirm } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { listDepartments, createDepartment, deleteDepartment } from "../api";
import type { Department } from "../types";

const { Title } = Typography;

const DepartmentsPage: React.FC = () => {
  const [deps, setDeps] = useState<Department[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  const load = () => { setLoading(true); listDepartments().then((r) => setDeps(r.data)).finally(() => setLoading(false)); };
  useEffect(() => { load(); }, []);

  const handleCreate = async (vals: any) => {
    try { await createDepartment(vals); message.success("Хэлтэс үүсгэлээ"); setModalOpen(false); form.resetFields(); load(); }
    catch { message.error("Алдаа"); }
  };

  const columns = [
    { title: "Нэр", dataIndex: "name" },
    { title: "Тайлбар", dataIndex: "description" },
    { title: "", render: (_: unknown, r: Department) => <Popconfirm title="Устгах?" onConfirm={() => deleteDepartment(r.id).then(load)}><Button size="small" danger>Устгах</Button></Popconfirm> },
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Хэлтэс</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>Шинэ хэлтэс</Button>
      </div>
      <Table rowKey="id" columns={columns} dataSource={deps} loading={loading} size="middle" />
      <Modal title="Шинэ хэлтэс" open={modalOpen} onCancel={() => setModalOpen(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="Нэр" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="description" label="Тайлбар"><Input.TextArea rows={2} /></Form.Item>
          <Button type="primary" htmlType="submit" block>Үүсгэх</Button>
        </Form>
      </Modal>
    </>
  );
};
export default DepartmentsPage;
