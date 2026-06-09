import { Navigate, Route, Routes } from "react-router-dom";
import ForgotPasswordForm from "../views/auth-forms/ForgotPassword";
import LoginForm from "../views/auth-forms/Login";
import ResetPasswordForm from "../views/auth-forms/ResetPassword";
import AuditLogList from "../views/audit/AuditLogList";
import Dashboard from "../views/dashboard/Dashboard";
import MainLayout from "../views/layouts/MainLayout";
import RoleList from "../views/roles/RoleList";
import PermissionList from "../views/permissions/PermissionList";
import SimpleLayout from "../views/layouts/SimpleLayout";
import SignUpForm from "../views/auth-forms/SignUp";
import UserList from "../views/users/UserList";
import Profile from "../views/profile/Profile";

function MainRouter() {
  return (
    <Routes>
      <Route element={<SimpleLayout />}>
        <Route path="/login" element={<LoginForm />} />
        <Route path="/signup" element={<SignUpForm />} />
        <Route path="/forgot-password" element={<ForgotPasswordForm />} />
        <Route path="/reset-password" element={<ResetPasswordForm />} />
      </Route>
      <Route element={<MainLayout />}>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/users" element={<UserList />} />
        <Route path="/roles" element={<RoleList />} />
        <Route path="/permissions" element={<PermissionList />} />
        <Route path="/audit-logs" element={<AuditLogList />} />
        <Route path="/profile" element={<Profile />} />
      </Route>
      <Route path="*" element={<Navigate to="/dashboard" />} />
    </Routes>
  );
}

export default MainRouter;
