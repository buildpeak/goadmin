import { Route, Routes } from "react-router-dom";
import LoginForm from "../views/auth-forms/Login";
import Dashboard from "../views/dashboard/Dashboard";
import MainLayout from "../views/layouts/MainLayout";
import SimpleLayout from "../views/layouts/SimpleLayout";
import SignUpForm from "../views/auth-forms/SignUp";

function MainRouter() {
  return (
    <Routes>
      <Route element={<SimpleLayout />}>
        <Route path="/login" element={<LoginForm />} />
        <Route path="/signup" element={<SignUpForm />} />
      </Route>
      <Route element={<MainLayout />}>
        <Route path="/dashboard" element={<Dashboard />} />
      </Route>
    </Routes>
  );
}

export default MainRouter;
