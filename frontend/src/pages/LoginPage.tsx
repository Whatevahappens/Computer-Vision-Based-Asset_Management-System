import React from "react";
import { Form, Input, Button, Card, Typography, message, Space } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import { login } from "../api";
import { useAuth } from "../context/AuthContext";

const { Title, Text } = Typography;

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const { setAuth } = useAuth();
  const [loading, setLoading] = React.useState(false);

  const onFinish = async (vals: { username: string; password: string }) => {
    setLoading(true);
    try {
      const { data } = await login(vals.username, vals.password);
      setAuth(data.token, data.user);
      message.success("Амжилттай нэвтэрлээ");
      navigate("/");
    } catch {
      message.error("Нэвтрэх нэр эсвэл нууц үг буруу байна");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center",
      background: "linear-gradient(135deg, #e8f4fd 0%, #f0f5ff 100%)",
    }}>
      <Card style={{ width: 400, borderRadius: 12, boxShadow: "0 4px 24px rgba(0,0,0,0.08)" }}>
        <Space direction="vertical" size="middle" style={{ width: "100%", textAlign: "center" }}>
          <div>
            <Title level={4} style={{ margin: 0, color: "#1677ff" }}>Asset Management</Title>
            <Text type="secondary">Ухаалаг эд хөрөнгийн бүртгэлийн систем</Text>
          </div>
          <Form layout="vertical" onFinish={onFinish} autoComplete="off" style={{ textAlign: "left" }}>
            <Form.Item label="Нэвтрэх нэр" name="username" rules={[{ required: true, message: "Нэвтрэх нэр оруулна уу" }]}>
              <Input prefix={<UserOutlined />} placeholder="admin" size="large" />
            </Form.Item>
            <Form.Item label="Нууц үг" name="password" rules={[{ required: true, message: "Нууц үг оруулна уу" }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="••••••" size="large" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" block size="large" loading={loading}>
                Нэвтрэх
              </Button>
            </Form.Item>
          </Form>
        </Space>
      </Card>
    </div>
  );
};

export default LoginPage;
