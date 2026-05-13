import React, { useEffect, useState } from "react";
import { Table, Button, Tag, Modal, Form, Input, Select, message, Space, Typography, Popconfirm } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { listUsers, createUser, deactivateUser, listDepartments } from "../api";
import type { User, Department } from "../types";

const { Title } = Typography;
const ROLE_MN: Record<string, string> = { ADMIN: "Админ", ASSET_CUSTODIAN: "Нярав", ACCOUNTANT: "Нягтлан бодогч", EMPLOYEE: "Ажилтан" };
const STATUS_COLORS: Record<string, string> = { ACTIVE: "green", INACTIVE: "default", SUSPENDED: "red" };

const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [form] = Form.useForm();

  const load = () => {
    setLoading(true);
    listUsers({ limit: 100 }).then((r) => setUsers(r.data.data || [])).finally(() => setLoading(false));
  };
  useEffect(() => { load(); listDepartments().then((r) => setDepartments(r.data)); }, []);

  const handleCreate = async (vals: any) => {
    try {
      await createUser(vals);
      message.success("Хэрэглэгч үүсгэгдлээ");
      setModalOpen(false); form.resetFields(); load();
    } catch { message.error("Алдаа гарлаа"); }
  };

  const columns = [
    { title: "Нэр", render: (_: unknown, r: User) => `${r.lastName?.charAt(0)}.${r.firstName}` },
    { title: "И-мэйл", dataIndex: "email" },
    { title: "Үүрэг", dataIndex: "role", render: (v: string) => <Tag color="blue">{ROLE_MN[v] || v}</Tag> },
    { title: "Хэлтэс", dataIndex: ["department", "name"], render: (v: string) => v || "—" },
    { title: "Төлөв", dataIndex: "status", render: (v: string) => <Tag color={STATUS_COLORS[v]}>{v}</Tag> },
    { title: "Үйлдэл", render: (_: unknown, r: User) => (
      <Popconfirm title="Идэвхгүй болгох уу?" onConfirm={() => deactivateUser(r.id).then(load)}>
        <Button size="small" danger disabled={r.status !== "ACTIVE"}>Идэвхгүй</Button>
      </Popconfirm>
    )},
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Хэрэглэгч удирдах</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>Шинэ хэрэглэгч</Button>
      </div>
      <Table rowKey="id" columns={columns} dataSource={users} loading={loading} size="middle" />
      <Modal title="Шинэ хэрэглэгч" open={modalOpen} onCancel={() => setModalOpen(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Space style={{ width: "100%" }}><Form.Item name="firstName" label="Овог" rules={[{ required: true }]} style={{ flex: 1 }}><Input /></Form.Item>
          <Form.Item name="lastName" label="Нэр" rules={[{ required: true }]} style={{ flex: 1 }}><Input /></Form.Item></Space>
          <Form.Item name="email" label="И-мэйл" rules={[{ required: true, type: "email" }]}><Input /></Form.Item>
          <Form.Item name="username" label="Нэвтрэх нэр" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="password" label="Нууц үг" rules={[{ required: true, min: 6 }]}><Input.Password /></Form.Item>
          <Form.Item name="phone" label="Утас"><Input /></Form.Item>
          <Form.Item name="role" label="Үүрэг" rules={[{ required: true }]}>
            <Select options={Object.entries(ROLE_MN).map(([k, v]) => ({ value: k, label: v }))} />
          </Form.Item>
          <Form.Item name="departmentId" label="Хэлтэс">
            <Select allowClear options={departments.map((d) => ({ value: d.id, label: d.name }))} />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>Үүсгэх</Button>
        </Form>
      </Modal>
    </>
  );
};
export default UsersPage;
