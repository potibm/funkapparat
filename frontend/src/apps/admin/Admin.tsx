import { Admin, Resource } from "react-admin";
import { MyTheme, MyDarkTheme } from "./theme/MyTheme";
import { MyLayout } from "./theme/MyLayout";
import { dataProvider } from "./providers/dataProvider";
import announcements from "./resources/announcement";

export const AdminApp = () => (
  <Admin
    dataProvider={dataProvider}
    theme={MyTheme}
    darkTheme={MyDarkTheme}
    layout={MyLayout}
    title="Funkapparat Admin"
  >
    <Resource {...announcements} />
  </Admin>
);
