import { SideMenuItemProps } from "./SideMenu"
import {
  BsGearFill,
  BsPaletteFill,
  BsCameraFill,
  BsWindow,
  BsPersonCircle,
  BsJoystick,
  BsMedium,
  BsFingerprint,
  BsCloudUploadFill,
  BsBucket,
  BsHddNetwork,
  BsArrowLeftRight,
} from "solid-icons/bs"
import { CgDatabase, CgShare } from "solid-icons/cg"
import { OcWorkflow2 } from "solid-icons/oc"
import { IoCopy, IoMove, IoHome, IoMagnetOutline } from "solid-icons/io"
import { Component, lazy } from "solid-js"
import { Group, UserRole } from "~/types"
import { TbArchive } from "solid-icons/tb"

export type SideMenuItem = SideMenuItemProps & {
  component?: Component
  children?: SideMenuItem[]
}

const CommonSettings = lazy(() => import("./settings/Common"))

export const side_menu_items: SideMenuItem[] = [
  {
    title: "manage.sidemenu.profile",
    icon: BsFingerprint,
    to: "/@manage",
    role: UserRole.GUEST,
    component: lazy(() => import("./users/Profile")),
  },
  {
    title: "manage.sidemenu.settings",
    icon: BsGearFill,
    to: "/@manage/settings",
    children: [
      {
        title: "manage.sidemenu.site",
        icon: BsWindow,
        to: "/@manage/settings/site",
        component: () => <CommonSettings group={Group.SITE} />,
      },
      {
        title: "manage.sidemenu.style",
        icon: BsPaletteFill,
        to: "/@manage/settings/style",
        component: () => <CommonSettings group={Group.STYLE} />,
      },
      {
        title: "manage.sidemenu.preview",
        icon: BsCameraFill,
        to: "/@manage/settings/preview",
        component: () => <CommonSettings group={Group.PREVIEW} />,
      },
      {
        title: "manage.sidemenu.global",
        icon: BsJoystick,
        to: "/@manage/settings/global",
        component: () => <CommonSettings group={Group.GLOBAL} />,
      },
      {
        title: "manage.sidemenu.s3",
        icon: BsBucket,
        to: "/@manage/settings/s3",
        component: lazy(() => import("./settings/S3")),
      },
      {
        title: "manage.sidemenu.ftp",
        icon: BsHddNetwork,
        to: "/@manage/settings/ftp",
        component: () => <CommonSettings group={Group.FTP} />,
      },
      {
        title: "manage.sidemenu.traffic",
        icon: BsArrowLeftRight,
        to: "/@manage/settings/traffic",
        component: () => <CommonSettings group={Group.TRAFFIC} />,
      },
      {
        title: "manage.sidemenu.other",
        icon: BsMedium,
        to: "/@manage/settings/other",
        component: lazy(() => import("./settings/Other")),
      },
    ],
  },
  {
    title: "manage.sidemenu.tasks",
    icon: OcWorkflow2,
    to: "/@manage/tasks",
    role: UserRole.GENERAL,
    children: [
      {
        title: "manage.sidemenu.upload",
        icon: BsCloudUploadFill,
        to: "/@manage/tasks/upload",
        role: UserRole.GENERAL,
        component: lazy(() => import("./tasks/Upload")),
      },
      {
        title: "manage.sidemenu.copy",
        icon: IoCopy,
        to: "/@manage/tasks/copy",
        role: UserRole.GENERAL,
        component: lazy(() => import("./tasks/Copy")),
      },
      {
        title: "manage.sidemenu.move",
        icon: IoMove,
        to: "/@manage/tasks/move",
        role: UserRole.GENERAL,
        component: lazy(() => import("./tasks/Move")),
      },
      {
        title: "manage.sidemenu.decompress",
        icon: TbArchive,
        to: "/@manage/tasks/decompress",
        role: UserRole.GENERAL,
        component: lazy(() => import("./tasks/Decompress")),
      },
    ],
  },
  {
    title: "manage.sidemenu.users",
    icon: BsPersonCircle,
    to: "/@manage/users",
    component: lazy(() => import("./users/Users")),
  },
  {
    title: "manage.sidemenu.storages",
    icon: CgDatabase,
    to: "/@manage/storages",
    component: lazy(() => import("./storages/Storages")),
  },
  {
    title: "manage.sidemenu.shares",
    icon: CgShare,
    to: "/@manage/shares",
    role: UserRole.GENERAL,
    component: lazy(() => import("./shares/Shares")),
  },
  {
    title: "manage.sidemenu.home",
    icon: IoHome,
    to: "/",
    role: UserRole.GUEST,
    refresh: true,
  },
]
