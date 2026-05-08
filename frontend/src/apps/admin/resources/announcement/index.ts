import EventIcon from "@mui/icons-material/Event";
import ScheduleEntriesList from "./AnnouncementList";
import AnnouncementCreate from "./AnnouncementCreate";
import AnnouncementEdit from "./AnnouncementEdit";

export default {
  name: "announcements",
  options: { label: "Announcements" },
  list: ScheduleEntriesList,
  create: AnnouncementCreate,
  edit: AnnouncementEdit,
  icon: EventIcon,
};
