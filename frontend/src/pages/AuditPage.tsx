import React, { useEffect, useState } from "react";
import {
  Typography, Button, Select, Upload, Card, Table, Tag, Statistic,
  Row, Col, message, Space, Spin, Modal, List,
} from "antd";
import { CameraOutlined, UploadOutlined, CheckCircleOutlined, WarningOutlined, CloseCircleOutlined } from "@ant-design/icons";
import { listLocations, startAudit, runCVAudit, listAudits, getAudit } from "../api";
import type { Location, AuditSession, AuditSummary } from "../types";
import dayjs from "dayjs";

const { Title, Text } = Typography;

const FINDING_COLORS: Record<string, string> = {
  MATCHED: "green", MISSING: "red", UNREGISTERED: "orange",
  DAMAGED: "volcano", MISPLACED: "blue",
};
const FINDING_MN: Record<string, string> = {
  MATCHED: "Тохирсон", MISSING: "Олдоогүй", UNREGISTERED: "Бүртгэлгүй",
  DAMAGED: "Гэмтсэн", MISPLACED: "Буруу байрласан",
};
const STATUS_COLORS: Record<string, string> = {
  COMPLETED: "green", IN_PROGRESS: "blue", PLANNED: "default", CANCELLED: "red",
};

const AuditPage: React.FC = () => {
  const [locations, setLocations] = useState<Location[]>([]);
  const [audits, setAudits] = useState<AuditSession[]>([]);
  const [loading, setLoading] = useState(false);
  const [cvModalOpen, setCvModalOpen] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState<string>("");
  const [cvLoading, setCvLoading] = useState(false);
  const [auditResult, setAuditResult] = useState<AuditSession | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [detailModal, setDetailModal] = useState<AuditSession | null>(null);

  useEffect(() => {
    listLocations().then((r) => setLocations(r.data));
    loadAudits();
  }, []);

  const loadAudits = () => {
    setLoading(true);
    listAudits({ limit: 50 }).then((r) => setAudits(r.data.data || [])).finally(() => setLoading(false));
  };

  const handleRunAudit = async () => {
    if (!selectedLocation || !file) {
      message.warning("Байршил болон зураг сонгоно уу");
      return;
    }
    setCvLoading(true);
    try {
      const { data: session } = await startAudit(selectedLocation);
      const { data: result } = await runCVAudit(session.id, file);
      setAuditResult(result);
      message.success("CV аудит амжилттай!");
      loadAudits();
    } catch (err: any) {
      message.error(err.response?.data?.error || "CV аудит амжилтгүй");
    } finally {
      setCvLoading(false);
    }
  };

  const viewDetail = async (id: string) => {
    const { data } = await getAudit(id);
    setDetailModal(data);
  };

  const totalRegistered = (s: AuditSummary[]) => s.reduce((a, b) => a + b.registeredCount, 0);
  const totalDetected = (s: AuditSummary[]) => s.reduce((a, b) => a + b.detectedCount, 0);
  const totalDiff = (s: AuditSummary[]) => s.reduce((a, b) => a + b.difference, 0);

  const columns = [
    { title: "Огноо", dataIndex: "startedAt", render: (v: string) => dayjs(v).format("YYYY-MM-DD HH:mm") },
    { title: "Байршил", dataIndex: ["location", "name"] },
    { title: "Төлөв", dataIndex: "status", render: (s: string) => <Tag color={STATUS_COLORS[s]}>{s}</Tag> },
    {
      title: "Үйлдэл", render: (_: unknown, rec: AuditSession) => (
        <Button size="small" onClick={() => viewDetail(rec.id)}>Дэлгэрэнгүй</Button>
      ),
    },
  ];

  return (
    <>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Компьютерийн харааны аудит</Title>
        <Button type="primary" icon={<CameraOutlined />} onClick={() => { setCvModalOpen(true); setAuditResult(null); setFile(null); }}>
          CV аудит хийх
        </Button>
      </div>

      {/* Stats */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}><Card><Statistic title="Нийт аудит" value={audits.length} /></Card></Col>
        <Col span={8}><Card><Statistic title="Тохирсон" value={audits.filter((a) => a.status === "COMPLETED").length} prefix={<CheckCircleOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
        <Col span={8}><Card><Statistic title="Зөрүүтэй" value={audits.filter((a) => a.summaries?.some((s) => s.difference !== 0)).length} prefix={<WarningOutlined />} valueStyle={{ color: "#fa8c16" }} /></Card></Col>
      </Row>

      <Table rowKey="id" columns={columns} dataSource={audits} loading={loading} size="middle" />

      {/* CV Audit Modal */}
      <Modal
        title="CV аудит хийх" open={cvModalOpen} width={600}
        onCancel={() => setCvModalOpen(false)} footer={null}
      >
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <div>
            <Text strong>Аудит хийх байршил</Text>
            <Select
              style={{ width: "100%", marginTop: 4 }} placeholder="Байршил сонгох"
              value={selectedLocation || undefined}
              onChange={setSelectedLocation}
              options={locations.map((l) => ({ value: l.id, label: `${l.name} — ${l.building || ""} ${l.room || ""}` }))}
            />
          </div>
          <div>
            <Text strong>Зураг боловсруулах хэсэг</Text>
            <Upload.Dragger
              accept="image/*" maxCount={1} beforeUpload={(f) => { setFile(f); return false; }}
              onRemove={() => setFile(null)}
              style={{ marginTop: 4 }}
            >
              <p className="ant-upload-drag-icon"><UploadOutlined style={{ fontSize: 32, color: "#1677ff" }} /></p>
              <p>Зураг оруулах</p>
            </Upload.Dragger>
          </div>

          <Button type="primary" block size="large" loading={cvLoading} onClick={handleRunAudit} disabled={!selectedLocation || !file}>
            Аудит эхлүүлэх
          </Button>

          {cvLoading && <Spin tip="YOLO26 загвараар илрүүлж байна..." style={{ display: "block", textAlign: "center" }} />}

          {auditResult && auditResult.summaries && (
            <Card title="Боловсруулалтын төлөв" style={{ marginTop: 8 }}>
              <Row gutter={16}>
                <Col span={8}><Statistic title="Бүртгэлтэй" value={totalRegistered(auditResult.summaries)} /></Col>
                <Col span={8}><Statistic title="Илэрсэн" value={totalDetected(auditResult.summaries)} /></Col>
                <Col span={8}>
                  <Statistic
                    title="Зөрүү" value={totalDiff(auditResult.summaries)}
                    valueStyle={{ color: totalDiff(auditResult.summaries) === 0 ? "#52c41a" : "#ff4d4f" }}
                  />
                </Col>
              </Row>
              {auditResult.findings && auditResult.findings.length > 0 && (
                <div style={{ marginTop: 12 }}>
                  <Text strong>Олдворууд:</Text>
                  <List
                    size="small" dataSource={auditResult.findings}
                    renderItem={(f) => (
                      <List.Item>
                        <Tag color={FINDING_COLORS[f.type]}>{FINDING_MN[f.type] || f.type}</Tag>
                        <Text style={{ fontSize: 12 }}>{f.notes}</Text>
                      </List.Item>
                    )}
                  />
                </div>
              )}
            </Card>
          )}
        </Space>
      </Modal>

      {/* Detail Modal */}
      <Modal
        title="Аудитын дэлгэрэнгүй" open={!!detailModal} width={600}
        onCancel={() => setDetailModal(null)} footer={null}
      >
        {detailModal && (
          <Space direction="vertical" size="middle" style={{ width: "100%" }}>
            <Row gutter={16}>
              <Col span={8}><Statistic title="Бүртгэлтэй" value={totalRegistered(detailModal.summaries || [])} /></Col>
              <Col span={8}><Statistic title="Илэрсэн" value={totalDetected(detailModal.summaries || [])} /></Col>
              <Col span={8}><Statistic title="Зөрүү" value={totalDiff(detailModal.summaries || [])} /></Col>
            </Row>
            {detailModal.findings && (
              <List
                header={<Text strong>Олдворууд</Text>} size="small" bordered
                dataSource={detailModal.findings}
                renderItem={(f) => (
                  <List.Item>
                    <Tag color={FINDING_COLORS[f.type]}>{FINDING_MN[f.type] || f.type}</Tag>
                    <Text>{f.notes}</Text>
                    {f.confidence > 0 && <Tag>{(f.confidence * 100).toFixed(1)}%</Tag>}
                  </List.Item>
                )}
              />
            )}
          </Space>
        )}
      </Modal>
    </>
  );
};

export default AuditPage;
