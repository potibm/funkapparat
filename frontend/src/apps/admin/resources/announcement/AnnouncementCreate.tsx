import { Create, SimpleForm } from "react-admin";
import { AnnouncementInputs } from "./AnnouncementInputs";

export const AnnouncementCreate = () => {
  return (
    <Create title="Add Annoucement">
      <SimpleForm>
        <AnnouncementInputs />
      </SimpleForm>
    </Create>
  );
};

export default AnnouncementCreate;
