export default [
  {
    path: "/attendees",
    name: "attendees",
    component: () => import("@/pages/AttendeeTable.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("@/pages/LoginPage.vue"),
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/pages/RegisterForm.vue"),
  },
  {
    path: "/forgot-password",
    name: "forgot-password",
    component: () => import("@/pages/ForgotPassword.vue"),
  },
  {
    path: "/forgot-password/sent",
    name: "forgot-password-sent",
    component: () => import("@/pages/ForgotPasswordSent.vue"),
  },
  {
    path: "/",
    name: "home",
    component: () => import("@/pages/HomePage.vue"),
  },
]
