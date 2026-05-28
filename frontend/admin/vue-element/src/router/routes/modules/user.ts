import type { RouteRecordRaw } from "vue-router";
import { Layout } from "@/layouts";

const userRoutes: RouteRecordRaw[] = [
  {
    name: "Profile",
    path: "/profile",
    component: Layout,
    redirect: "/profile/user",
    meta: {
      title: "routes.profile.settings",
      hideInMenu: true,
    },

    children: [
      {
        path: "user",
        name: "ProfilePage",
        component: () => import("@/pages/core/profile/index.vue"),
        meta: {
          title: "routes.profile.settings",
          icon: "lucide:user-pen",
          hideInMenu: true,
        },
      },
    ],
  },
];

export default userRoutes;
