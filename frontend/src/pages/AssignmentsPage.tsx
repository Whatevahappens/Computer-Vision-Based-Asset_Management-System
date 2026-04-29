import React, { useEffect, useState } from "react";
import { Typography, Table, Button, Modal, Form, Select, Input, message, Tag } from "antd";
import { listAssets, assignAsset, transferAsset, disposeAsset, listUsers, listLocations } from "../api";
import type { Asset, User, Location } from "../types";

const { Title } = Typography;

const AssignmentsPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState<{ type: "assign" | "transfer" | "dispose"; assetId: string } | null>(null);
  const [form] = Form.useForm();

  const load = () => { setLoading(true); listAssets({ limit: 200 }).then((r) => setAssets(r.data.data || [])).finally(() => setLoading(false)); };
  useEffect(() => {
    load();
    listUsers({ limit: 200 }).then((r) => setUsers(r.data.data || []));
    listLocations().then((r) => setLocations(r.data));
  }, []);

  const handleSubmit = async (vals: any) => {
    if (!modal) return;
    try {
      if (modal.type === "assign") await assignAsset(modal.assetId, { userId: vals.userId, locationId: vals.locationId, notes: vals.notes });
      else if (modal.type === "transfer") await transferAsset(modal.assetId, { toUserId: vals.userId, locationId: vals.locationId, notes: vals.notes });
      else await disposeAsset(modal.assetId, { reason: vals.reason, residualValue: vals.residualValue || 0, notes: vals.notes });
      message.success("Амжилттай"); setModal(null); form.resetFields(); load();
    } catch { message.error("Алдаа"); }
  };

  const columns = [
    { title: "Код", dataIndex: "barcode", width: 120 },
    { title: "Нэр", dataIndex: "assetName" },
    { title: "Хариуцагч", render: (_: unknown, r: Asset) => r.assignedUser ? `${r.assignedUser.lastName?.charAt(0)}.${r.assignedUser.firstName}` : "—" },
    { title: "Байршил", dataIndex: ["location", "name"], render: (v: string) => v || "—" },
    { title: "Төлөв", dataIndex: "status", render: (s: string) => <Tag color={s === "ACTIVE" ? "green" : "red"}>{s}</Tag> },
    { title: "Үйлдэл", render: (_: unknown, r: Asset) => r.status === "ACTIVE" && (
      <>
        <Button size="small" onClick={() => { setModal({ type: "assign", assetId: r.id }); form.resetFields(); }}>Хуваарилах</Button>{" "}
        {r.assignedUserId && <Button size="small" onClick={() => { setModal({ type: "transfer", assetId: r.id }); form.resetFields(); }}>Шилжүүлэх</Button>}{" "}
        <Button size="small" danger onClick={() => { setModal({ type: "dispose", assetId: r.id }); form.resetFields(); }}>Актлах</Button>
      </>
    )},
  ];

  const isDispose = modal?.type === "dispose";
  const titleMap = { assign: "Эд хөрөнгө хуваарилах", transfer: "Эд хөрөнгө шилжүүлэх", dispose: "Эд хөрөнгө актлах" };

  return (
    <>
      <Title level={4}>Хуваарилалт удирдах</Title>
      <Table rowKey="id" columns={columns} dataSource={assets} loading={loading} size="middle" scroll={{ x: 800 }} />
      <Modal title={modal ? titleMap[modal.type] : ""} open={!!modal} onCancel={() => setModal(null)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          {!isDispose && <Form.Item name="userId" label="Хэрэглэгч" rules={[{ required: true }]}>
            <Select options={users.map((u) => ({ value: u.id, label: `${u.lastName?.charAt(0)}.${u.firstName} (${u.role})` }))} />
          </Form.Item>}
          {!isDispose && <Form.Item name="locationId" label="Байршил">
            <Select allowClear options={locations.map((l) => ({ value: l.id, label: l.name }))} />
          </Form.Item>}
          {isDispose && <Form.Item name="reason" label="Шалтгаан" rules={[{ required: true }]}><Input /></Form.Item>}
          <Form.Item name="notes" label="Тайлбар"><Input.TextArea rows={2} /></Form.Item>
          <Button type="primary" htmlType="submit" block>{isDispose ? "Актлах" : "Хадгалах"}</Button>
        </Form>
      </Modal>
    </>
  );
};
export default AssignmentsPage;
