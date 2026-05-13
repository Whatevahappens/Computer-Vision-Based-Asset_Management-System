import React, { useEffect, useState } from "react";
import {
  Table, Button, Space, Tag, Input, Typography, Modal, Form, Select,
  InputNumber, DatePicker, message, Popconfirm, Drawer, Timeline,
} from "antd";
import { PlusOutlined, SearchOutlined, HistoryOutlined } from "@ant-design/icons";
import { listAssets, createAsset, deleteAsset, getAssetHistory, listLocations, listDepartments, listAssetModels } from "../api";
import type { Asset, Location, Department, AssetModel, AssetHistory } from "../types";
import dayjs from "dayjs";

const { Title } = Typography;

const STATUS_COLORS: Record<string, string> = {
  ACTIVE: "green", INACTIVE: "default", DISPOSED: "red",
  LOST: "orange", UNDER_MAINTENANCE: "blue",
};
const STATUS_MN: Record<string, string> = {
  ACTIVE: "Идэвхтэй", INACTIVE: "Идэвхгүй", DISPOSED: "Актлагдсан",
  LOST: "Алдагдсан", UNDER_MAINTENANCE: "Засварт",
};

const AssetsPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [modalOpen, setModalOpen] = useState(false);
  const [historyDrawer, setHistoryDrawer] = useState<{ open: boolean; assetId: string; name: string }>({ open: false, assetId: "", name: "" });
  const [history, setHistory] = useState<AssetHistory[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [models, setModels] = useState<AssetModel[]>([]);
  const [form] = Form.useForm();

  const load = () => {
    setLoading(true);
    listAssets({ page, limit: 15, search })
      .then((r) => { setAssets(r.data.data || []); setTotal(r.data.total); })
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, [page, search]);
  useEffect(() => {
    listLocations().then((r) => setLocations(r.data));
    listDepartments().then((r) => setDepartments(r.data));
    listAssetModels().then((r) => setModels(r.data));
  }, []);

  const openHistory = (assetId: string, name: string) => {
    setHistoryDrawer({ open: true, assetId, name });
    getAssetHistory(assetId).then((r) => setHistory(r.data));
  };

  const handleCreate = async (vals: any) => {
    try {
      await createAsset({
        ...vals,
        acquisitionDate: vals.acquisitionDate.format("YYYY-MM-DD"),
      });
      message.success("Эд хөрөнгө бүртгэгдлээ");
      setModalOpen(false);
      form.resetFields();
      load();
    } catch {
      message.error("Алдаа гарлаа");
    }
  };

  const columns = [
    { title: "Код", dataIndex: "barcode", width: 130 },
    { title: "Нэр", dataIndex: "assetName", ellipsis: true },
    {
      title: "Ангилал", dataIndex: ["assetModel", "category"], width: 140,
      render: (v: string) => v || "—",
    },
    {
      title: "Байршил", dataIndex: ["location", "name"], width: 100,
      render: (v: string) => v || "—",
    },
    {
      title: "Үнэ", dataIndex: "acquisitionPrice", width: 120,
      render: (v: number) => `${v?.toLocaleString()} ₮`,
    },
    {
      title: "Төлөв", dataIndex: "status", width: 110,
      render: (s: string) => <Tag color={STATUS_COLORS[s]}>{STATUS_MN[s] || s}</Tag>,
    },
    {
      title: "Үйлдэл", width: 120,
      render: (_: unknown, rec: Asset) => (
        <Space size="small">
          <Button size="small" icon={<HistoryOutlined />} onClick={() => openHistory(rec.id, rec.assetName)} />
          <Popconfirm title="Актлах уу?" onConfirm={() => deleteAsset(rec.id).then(load)}>
            <Button size="small" danger>Устгах</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Эд хөрөнгө удирдах</Title>
        <Space>
          <Input
            placeholder="Хайх..." prefix={<SearchOutlined />} allowClear
            style={{ width: 220 }}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            Шинэ эд хөрөнгө
          </Button>
        </Space>
      </div>

      <Table
        rowKey="id" columns={columns} dataSource={assets} loading={loading}
        pagination={{ current: page, total, pageSize: 15, onChange: setPage, showSizeChanger: false }}
        size="middle" scroll={{ x: 800 }}
      />

      {/* Create Modal */}
      <Modal
        title="Шинэ эд хөрөнгө бүртгэх" open={modalOpen}
        onCancel={() => setModalOpen(false)} footer={null} width={520}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="assetName" label="Нэр" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="serialNumber" label="Сериал дугаар">
            <Input />
          </Form.Item>
          <Form.Item name="assetModelId" label="Загвар">
            <Select allowClear placeholder="Сонгох" options={
              models.map((m) => ({ value: m.id, label: `${m.brand} ${m.modelName}` }))
            } />
          </Form.Item>
          <Space style={{ width: "100%" }} size="middle">
            <Form.Item name="acquisitionPrice" label="Үнэ" rules={[{ required: true }]} style={{ flex: 1 }}>
              <InputNumber style={{ width: "100%" }} min={0} />
            </Form.Item>
            <Form.Item name="usefulLifeMonths" label="Ашиглалт (сар)" rules={[{ required: true }]} style={{ flex: 1 }}>
              <InputNumber style={{ width: "100%" }} min={1} />
            </Form.Item>
          </Space>
          <Form.Item name="acquisitionDate" label="Худалдан авсан огноо" rules={[{ required: true }]}>
            <DatePicker style={{ width: "100%" }} />
          </Form.Item>
          <Space style={{ width: "100%" }} size="middle">
            <Form.Item name="locationId" label="Байршил" style={{ flex: 1 }}>
              <Select allowClear placeholder="Сонгох" options={locations.map((l) => ({ value: l.id, label: l.name }))} />
            </Form.Item>
            <Form.Item name="departmentId" label="Хэлтэс" style={{ flex: 1 }}>
              <Select allowClear placeholder="Сонгох" options={departments.map((d) => ({ value: d.id, label: d.name }))} />
            </Form.Item>
          </Space>
          <Form.Item name="description" label="Тайлбар">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>Бүртгэх</Button>
        </Form>
      </Modal>

      {/* History Drawer */}
      <Drawer
        title={`${historyDrawer.name} — Түүх`} open={historyDrawer.open}
        onClose={() => setHistoryDrawer({ open: false, assetId: "", name: "" })}
        width={400}
      >
        <Timeline items={history.map((h) => ({
          children: (
            <div>
              <Tag>{h.changeType}</Tag>
              <span style={{ fontSize: 12, color: "#888" }}>{dayjs(h.changedAt).format("YYYY-MM-DD HH:mm")}</span>
              <div style={{ marginTop: 4 }}>{h.description}</div>
            </div>
          ),
        }))} />
      </Drawer>
    </>
  );
};

export default AssetsPage;
